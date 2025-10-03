package middleware

import (
	"sync"
	"time"

	kperrors "github.com/cuiyuanxin/kunpeng/pkg/errors"
	"github.com/cuiyuanxin/kunpeng/pkg/logger"
	"github.com/cuiyuanxin/kunpeng/pkg/response"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"golang.org/x/time/rate"
)

// RateLimiterConfig 限流配置
type RateLimiterConfig struct {
	Rate       int           // 每秒允许的请求数
	Burst      int           // 允许的突发请求数
	ExpiresIn  time.Duration // 限流器过期时间
	CleanupInt time.Duration // 清理间隔
}

// DefaultRateLimiterConfig 默认限流配置
var DefaultRateLimiterConfig = RateLimiterConfig{
	Rate:       1000,        // 每秒1000个请求（更宽松）
	Burst:      2000,        // 突发2000个请求（更宽松）
	ExpiresIn:  time.Hour,   // 1小时过期
	CleanupInt: time.Minute, // 每分钟清理一次
}

// IPRateLimiter IP限流器
type IPRateLimiter struct {
	ips      map[string]*ipLimiter
	mu       *sync.RWMutex
	config   RateLimiterConfig
	stopChan chan struct{}
}

// ipLimiter IP限流器
type ipLimiter struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

// NewIPRateLimiter 创建IP限流器
func NewIPRateLimiter(config RateLimiterConfig) *IPRateLimiter {
	ipl := &IPRateLimiter{
		ips:      make(map[string]*ipLimiter),
		mu:       &sync.RWMutex{},
		config:   config,
		stopChan: make(chan struct{}),
	}

	// 启动清理过期限流器的协程
	go ipl.cleanupLoop()

	return ipl
}

// addIPUnsafe 添加IP限流器（不加锁版本，内部使用）
func (ipl *IPRateLimiter) addIPUnsafe(ip string) *rate.Limiter {
	limiter := rate.NewLimiter(rate.Limit(ipl.config.Rate), ipl.config.Burst)
	ipl.ips[ip] = &ipLimiter{
		limiter:  limiter,
		lastSeen: time.Now(),
	}
	return limiter
}

// AddIP 添加IP限流器
func (ipl *IPRateLimiter) AddIP(ip string) *rate.Limiter {
	ipl.mu.Lock()
	defer ipl.mu.Unlock()
	return ipl.addIPUnsafe(ip)
}

// GetLimiter 获取IP限流器
func (ipl *IPRateLimiter) GetLimiter(ip string) *rate.Limiter {
	ipl.mu.Lock()
	defer ipl.mu.Unlock()

	limiter, exists := ipl.ips[ip]
	if !exists {
		return ipl.addIPUnsafe(ip)
	}

	// 更新最后访问时间
	limiter.lastSeen = time.Now()
	return limiter.limiter
}

// cleanupLoop 清理过期限流器
func (ipl *IPRateLimiter) cleanupLoop() {
	ticker := time.NewTicker(ipl.config.CleanupInt)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			ipl.mu.Lock()
			for ip, limiter := range ipl.ips {
				if time.Since(limiter.lastSeen) > ipl.config.ExpiresIn {
					delete(ipl.ips, ip)
				}
			}
			ipl.mu.Unlock()
		case <-ipl.stopChan:
			return
		}
	}
}

// Stop 停止清理协程
func (ipl *IPRateLimiter) Stop() {
	close(ipl.stopChan)
}

// 全局IP限流器
var globalIPLimiter *IPRateLimiter

// InitRateLimiter 初始化限流器
func InitRateLimiter(config RateLimiterConfig) {
	globalIPLimiter = NewIPRateLimiter(config)
}

// RateLimiter 中间件，用于限制请求频率
func RateLimiter() gin.HandlerFunc {
	// 如果全局限流器未初始化，则使用默认配置初始化
	if globalIPLimiter == nil {
		InitRateLimiter(DefaultRateLimiterConfig)
		logger.GetLogger().Info("限流中间件已初始化",
			zap.Int("rate", DefaultRateLimiterConfig.Rate),
			zap.Int("burst", DefaultRateLimiterConfig.Burst))
	}

	return func(c *gin.Context) {
		// 获取客户端IP
		ip := c.ClientIP()

		// 获取限流器
		limiter := globalIPLimiter.GetLimiter(ip)

		// 检查是否允许请求
		if !limiter.Allow() {
			logger.GetLogger().Warn("请求被限流",
				zap.String("ip", ip),
				zap.String("path", c.Request.URL.Path))
			response.FailWithCode(c, kperrors.ErrTooManyRequests)
			c.Abort()
			return
		}

		// 记录允许的请求（仅在调试模式下）
		if gin.Mode() == gin.DebugMode {
			logger.GetLogger().Debug("请求通过限流检查",
				zap.String("ip", ip),
				zap.String("path", c.Request.URL.Path))
		}

		c.Next()
	}
}
