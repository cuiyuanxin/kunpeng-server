package jwt

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/cuiyuanxin/kunpeng/pkg/config"
	kperrors "github.com/cuiyuanxin/kunpeng/pkg/errors"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// TokenType token类型
type TokenType string

const (
	AccessTokenType  TokenType = "access"
	RefreshTokenType TokenType = "refresh"
)

// CustomClaims 自定义JWT声明
type CustomClaims struct {
	UserID     uint      `json:"user_id"`
	Username   string    `json:"username"`
	RoleID     uint      `json:"role_id"`
	AppKey     string    `json:"app_key"`
	RememberMe bool      `json:"remember_me"`
	TokenType  TokenType `json:"token_type"`
	jwt.RegisteredClaims
}

// TokenPair token对
type TokenPair struct {
	AccessToken      string `json:"access_token"`
	RefreshToken     string `json:"refresh_token"`
	ExpiresIn        int64  `json:"expires_in"`         // access token过期时间（秒）
	RefreshExpiresIn int64  `json:"refresh_expires_in"` // refresh token过期时间（秒）
}

// GenerateTokenPair 生成JWT令牌对（access token和refresh token）
func GenerateTokenPair(userID uint, username string, roleID uint, appKey, appSecret string, rememberMe bool) (*TokenPair, error) {
	// 生成access token
	accessToken, accessExpiresIn, err := generateSingleToken(userID, username, roleID, appKey, appSecret, rememberMe, AccessTokenType)
	if err != nil {
		return nil, err
	}

	// 生成refresh token
	refreshToken, refreshExpiresIn, err := generateSingleToken(userID, username, roleID, appKey, appSecret, rememberMe, RefreshTokenType)
	if err != nil {
		return nil, err
	}

	return &TokenPair{
		AccessToken:      accessToken,
		RefreshToken:     refreshToken,
		ExpiresIn:        accessExpiresIn,
		RefreshExpiresIn: refreshExpiresIn,
	}, nil
}

// generateSingleToken 生成单个JWT令牌
func generateSingleToken(userID uint, username string, roleID uint, appKey, appSecret string, rememberMe bool, tokenType TokenType) (string, int64, error) {
	jwtConfig := config.GetJWTConfig()

	// 二次加密
	encryptedKey := encryptWithAppSecret(appKey, appSecret)

	// 根据token类型和记住我选项确定过期时间
	var expireTime time.Duration
	if tokenType == RefreshTokenType {
		// refresh token有更长的过期时间
		if rememberMe {
			expireTime = jwtConfig.RememberMeExpireTime * 2 // refresh token过期时间是记住我的2倍
		} else {
			expireTime = jwtConfig.ExpireTime * 7 // refresh token过期时间是普通token的7倍
		}
	} else {
		// access token使用正常过期时间
		if rememberMe {
			expireTime = jwtConfig.RememberMeExpireTime
		} else {
			expireTime = jwtConfig.ExpireTime
		}
	}

	// 创建声明
	claims := CustomClaims{
		UserID:     userID,
		Username:   username,
		RoleID:     roleID,
		RememberMe: rememberMe,
		AppKey:     encryptedKey,
		TokenType:  tokenType,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expireTime)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    jwtConfig.Issuer,
			Subject:   username,
		},
	}

	// 创建令牌
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// 签名令牌
	tokenString, err := token.SignedString([]byte(jwtConfig.Secret))
	if err != nil {
		return "", 0, err
	}

	// 返回token字符串和过期时间（秒）
	return tokenString, int64(expireTime.Seconds()), nil
}

// ParseToken 解析JWT令牌
func ParseToken(tokenString string) (*CustomClaims, error) {
	jwtConfig := config.GetJWTConfig()

	// 解析令牌
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		// 验证签名方法
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(jwtConfig.Secret), nil
	})

	if err != nil {
		// 检查是否是过期错误
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, kperrors.New(kperrors.ErrTokenExpired, err)
		}
		return nil, kperrors.New(kperrors.ErrInvalidToken, err)
	}

	// 获取声明
	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, kperrors.New(kperrors.ErrInvalidToken, nil)
}

// RefreshTokenPair 通过refresh token刷新token对
func RefreshTokenPair(refreshTokenString string, appSecret string) (*TokenPair, error) {
	// 解析refresh token
	claims, err := ParseToken(refreshTokenString)
	if err != nil {
		return nil, err
	}

	// 验证token类型必须是refresh token
	if claims.TokenType != RefreshTokenType {
		return nil, kperrors.New(kperrors.ErrInvalidToken, nil)
	}

	// 验证二次加密
	if !validateEncryption(claims.AppKey, claims.Username, appSecret) {
		return nil, kperrors.New(kperrors.ErrInvalidToken, nil)
	}

	// 生成新的token对
	return GenerateTokenPair(claims.UserID, claims.Username, claims.RoleID, claims.Username, appSecret, claims.RememberMe)
}

// encryptWithAppSecret 使用AppSecret进行二次加密
func encryptWithAppSecret(appKey, appSecret string) string {
	// 使用HMAC-SHA256进行安全加密
	h := hmac.New(sha256.New, []byte(appSecret))
	h.Write([]byte(appKey))
	return hex.EncodeToString(h.Sum(nil))
}

// validateEncryption 验证二次加密
func validateEncryption(encryptedKey, appKey, appSecret string) bool {
	return encryptedKey == encryptWithAppSecret(appKey, appSecret)
}

// GenerateAppKeyAndSecret 生成AppKey和AppSecret
func GenerateAppKeyAndSecret(username string) (string, string) {
	// 生成随机的AppKey (20位字符)
	appKeyBytes := make([]byte, 15) // 15字节 = 20位十六进制字符
	if _, err := rand.Read(appKeyBytes); err != nil {
		// 如果随机数生成失败，使用时间戳和用户名作为后备方案
		timestamp := time.Now().UnixNano()
		h := sha256.New()
		h.Write([]byte(fmt.Sprintf("%s%d", username, timestamp)))
		hash := h.Sum(nil)
		appKey := hex.EncodeToString(hash[:10])    // 取前10字节，生成20位十六进制
		appSecret := hex.EncodeToString(hash[10:]) // 取后面部分作为AppSecret
		return appKey, appSecret
	}
	appKey := hex.EncodeToString(appKeyBytes)

	// 生成加密安全的随机AppSecret
	randomBytes := make([]byte, 32) // 256位随机数
	if _, err := rand.Read(randomBytes); err != nil {
		// 如果随机数生成失败，使用时间戳作为后备方案
		timestamp := time.Now().UnixNano()
		h := sha256.New()
		h.Write([]byte(fmt.Sprintf("%s%d", username, timestamp)))
		appSecret := hex.EncodeToString(h.Sum(nil))
		return appKey, appSecret
	}

	// 使用SHA256哈希随机字节
	h := sha256.New()
	h.Write(randomBytes)
	h.Write([]byte(username)) // 添加用户名增加唯一性
	appSecret := hex.EncodeToString(h.Sum(nil))

	return appKey, appSecret
}

// GetUserID 从上下文中获取用户ID
func GetUserID(c *gin.Context) uint {
	if userID, exists := c.Get("user_id"); exists {
		if id, ok := userID.(uint); ok {
			return id
		}
	}
	return 0
}

// GetUsername 从上下文中获取用户名
func GetUsername(c *gin.Context) string {
	if username, exists := c.Get("username"); exists {
		if name, ok := username.(string); ok {
			return name
		}
	}
	return ""
}

// GetRoleID 从上下文中获取角色ID
func GetRoleID(c *gin.Context) uint {
	if roleID, exists := c.Get("role_id"); exists {
		if id, ok := roleID.(uint); ok {
			return id
		}
	}
	return 0
}

// GetAppKey 从上下文中获取AppKey
func GetAppKey(c *gin.Context) string {
	if appKey, exists := c.Get("app_key"); exists {
		if key, ok := appKey.(string); ok {
			return key
		}
	}
	return ""
}

// ExtractTokenFromHeader 从HTTP请求头中提取JWT token
// 支持 "Authorization: Bearer <token>" 格式
// 返回提取的token字符串，如果格式不正确或不存在则返回空字符串
func ExtractTokenFromHeader(authHeader string) string {
	if authHeader == "" {
		return ""
	}

	// 检查是否以 "Bearer " 开头
	if !strings.HasPrefix(authHeader, "Bearer ") {
		return ""
	}

	// 分割字符串并验证格式
	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || parts[0] != "Bearer" {
		return ""
	}

	return parts[1]
}
