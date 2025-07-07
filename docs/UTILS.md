# 工具函数库文档

## 概述

本项目提供了一套完整的工具函数库，涵盖字符串处理、时间操作、文件处理、加密解密、数据验证、网络请求等常用功能。工具函数库遵循单一职责原则，提供高性能、易用的 API 接口。

## 系统架构

### 工具库组织结构

```
pkg/utils/
├── string.go              # 字符串处理工具
├── time.go                # 时间处理工具
├── file.go                # 文件操作工具
├── crypto.go              # 加密解密工具
├── validator.go           # 数据验证工具
├── http.go                # HTTP 请求工具
├── json.go                # JSON 处理工具
├── slice.go               # 切片操作工具
├── map.go                 # 映射操作工具
├── convert.go             # 类型转换工具
├── random.go              # 随机数生成工具
├── hash.go                # 哈希计算工具
├── compress.go            # 压缩解压工具
├── image.go               # 图片处理工具
└── pagination.go          # 分页工具
```

## 字符串处理工具

### 基础字符串操作

```go
// string.go
package utils

import (
    "regexp"
    "strings"
    "unicode"
    "unicode/utf8"
)

// IsEmpty 检查字符串是否为空
func IsEmpty(s string) bool {
    return len(strings.TrimSpace(s)) == 0
}

// IsNotEmpty 检查字符串是否不为空
func IsNotEmpty(s string) bool {
    return !IsEmpty(s)
}

// Capitalize 首字母大写
func Capitalize(s string) string {
    if s == "" {
        return s
    }
    
    r, size := utf8.DecodeRuneInString(s)
    if r == utf8.RuneError {
        return s
    }
    
    return string(unicode.ToUpper(r)) + s[size:]
}

// Reverse 反转字符串
func Reverse(s string) string {
    runes := []rune(s)
    for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
        runes[i], runes[j] = runes[j], runes[i]
    }
    return string(runes)
}

// ContainsIgnoreCase 忽略大小写检查子字符串
func ContainsIgnoreCase(s, substr string) bool {
    return strings.Contains(strings.ToLower(s), strings.ToLower(substr))
}

// TruncateString 截断字符串
func TruncateString(s string, maxLen int, suffix string) string {
    if len(s) <= maxLen {
        return s
    }
    
    if len(suffix) >= maxLen {
        return suffix[:maxLen]
    }
    
    return s[:maxLen-len(suffix)] + suffix
}

// PadLeft 左填充字符串
func PadLeft(s string, length int, pad string) string {
    if len(s) >= length {
        return s
    }
    
    padLen := length - len(s)
    padStr := strings.Repeat(pad, padLen/len(pad)+1)
    return padStr[:padLen] + s
}

// PadRight 右填充字符串
func PadRight(s string, length int, pad string) string {
    if len(s) >= length {
        return s
    }
    
    padLen := length - len(s)
    padStr := strings.Repeat(pad, padLen/len(pad)+1)
    return s + padStr[:padLen]
}

// CamelCase 转换为驼峰命名
func CamelCase(s string) string {
    words := strings.FieldsFunc(s, func(r rune) bool {
        return !unicode.IsLetter(r) && !unicode.IsNumber(r)
    })
    
    if len(words) == 0 {
        return ""
    }
    
    result := strings.ToLower(words[0])
    for i := 1; i < len(words); i++ {
        result += Capitalize(strings.ToLower(words[i]))
    }
    
    return result
}

// SnakeCase 转换为蛇形命名
func SnakeCase(s string) string {
    re := regexp.MustCompile(`([a-z0-9])([A-Z])`)
    snake := re.ReplaceAllString(s, `${1}_${2}`)
    return strings.ToLower(snake)
}

// KebabCase 转换为短横线命名
func KebabCase(s string) string {
    re := regexp.MustCompile(`([a-z0-9])([A-Z])`)
    kebab := re.ReplaceAllString(s, `${1}-${2}`)
    return strings.ToLower(kebab)
}

// RemoveSpecialChars 移除特殊字符
func RemoveSpecialChars(s string) string {
    re := regexp.MustCompile(`[^a-zA-Z0-9\s]`)
    return re.ReplaceAllString(s, "")
}

// ExtractNumbers 提取字符串中的数字
func ExtractNumbers(s string) []string {
    re := regexp.MustCompile(`\d+`)
    return re.FindAllString(s, -1)
}

// MaskString 掩码字符串（用于敏感信息）
func MaskString(s string, start, end int, mask string) string {
    if start < 0 || end > len(s) || start >= end {
        return s
    }
    
    maskLen := end - start
    maskStr := strings.Repeat(mask, maskLen)
    
    return s[:start] + maskStr + s[end:]
}

// MaskEmail 掩码邮箱地址
func MaskEmail(email string) string {
    parts := strings.Split(email, "@")
    if len(parts) != 2 {
        return email
    }
    
    username := parts[0]
    domain := parts[1]
    
    if len(username) <= 2 {
        return email
    }
    
    maskedUsername := username[:1] + strings.Repeat("*", len(username)-2) + username[len(username)-1:]
    return maskedUsername + "@" + domain
}

// MaskPhone 掩码手机号码
func MaskPhone(phone string) string {
    if len(phone) < 7 {
        return phone
    }
    
    return phone[:3] + strings.Repeat("*", len(phone)-6) + phone[len(phone)-3:]
}
```

## 时间处理工具

```go
// time.go
package utils

import (
    "fmt"
    "time"
)

const (
    DateFormat     = "2006-01-02"
    TimeFormat     = "15:04:05"
    DateTimeFormat = "2006-01-02 15:04:05"
    RFC3339Format  = time.RFC3339
)

// FormatTime 格式化时间
func FormatTime(t time.Time, format string) string {
    return t.Format(format)
}

// ParseTime 解析时间字符串
func ParseTime(timeStr, format string) (time.Time, error) {
    return time.Parse(format, timeStr)
}

// GetCurrentTime 获取当前时间
func GetCurrentTime() time.Time {
    return time.Now()
}

// GetCurrentTimeString 获取当前时间字符串
func GetCurrentTimeString(format string) string {
    return time.Now().Format(format)
}

// GetTimestamp 获取时间戳
func GetTimestamp() int64 {
    return time.Now().Unix()
}

// GetMillisTimestamp 获取毫秒时间戳
func GetMillisTimestamp() int64 {
    return time.Now().UnixMilli()
}

// TimestampToTime 时间戳转时间
func TimestampToTime(timestamp int64) time.Time {
    return time.Unix(timestamp, 0)
}

// MillisTimestampToTime 毫秒时间戳转时间
func MillisTimestampToTime(timestamp int64) time.Time {
    return time.UnixMilli(timestamp)
}

// GetStartOfDay 获取一天的开始时间
func GetStartOfDay(t time.Time) time.Time {
    return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
}

// GetEndOfDay 获取一天的结束时间
func GetEndOfDay(t time.Time) time.Time {
    return time.Date(t.Year(), t.Month(), t.Day(), 23, 59, 59, 999999999, t.Location())
}

// GetStartOfWeek 获取一周的开始时间（周一）
func GetStartOfWeek(t time.Time) time.Time {
    weekday := int(t.Weekday())
    if weekday == 0 {
        weekday = 7 // 将周日调整为7
    }
    return GetStartOfDay(t.AddDate(0, 0, -(weekday-1)))
}

// GetEndOfWeek 获取一周的结束时间（周日）
func GetEndOfWeek(t time.Time) time.Time {
    return GetEndOfDay(GetStartOfWeek(t).AddDate(0, 0, 6))
}

// GetStartOfMonth 获取一月的开始时间
func GetStartOfMonth(t time.Time) time.Time {
    return time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, t.Location())
}

// GetEndOfMonth 获取一月的结束时间
func GetEndOfMonth(t time.Time) time.Time {
    return GetStartOfMonth(t).AddDate(0, 1, 0).Add(-time.Nanosecond)
}

// GetStartOfYear 获取一年的开始时间
func GetStartOfYear(t time.Time) time.Time {
    return time.Date(t.Year(), 1, 1, 0, 0, 0, 0, t.Location())
}

// GetEndOfYear 获取一年的结束时间
func GetEndOfYear(t time.Time) time.Time {
    return time.Date(t.Year(), 12, 31, 23, 59, 59, 999999999, t.Location())
}

// DiffDays 计算两个时间相差的天数
func DiffDays(t1, t2 time.Time) int {
    return int(t2.Sub(t1).Hours() / 24)
}

// DiffHours 计算两个时间相差的小时数
func DiffHours(t1, t2 time.Time) int {
    return int(t2.Sub(t1).Hours())
}

// DiffMinutes 计算两个时间相差的分钟数
func DiffMinutes(t1, t2 time.Time) int {
    return int(t2.Sub(t1).Minutes())
}

// IsToday 判断是否为今天
func IsToday(t time.Time) bool {
    now := time.Now()
    return t.Year() == now.Year() && t.Month() == now.Month() && t.Day() == now.Day()
}

// IsYesterday 判断是否为昨天
func IsYesterday(t time.Time) bool {
    yesterday := time.Now().AddDate(0, 0, -1)
    return t.Year() == yesterday.Year() && t.Month() == yesterday.Month() && t.Day() == yesterday.Day()
}

// IsTomorrow 判断是否为明天
func IsTomorrow(t time.Time) bool {
    tomorrow := time.Now().AddDate(0, 0, 1)
    return t.Year() == tomorrow.Year() && t.Month() == tomorrow.Month() && t.Day() == tomorrow.Day()
}

// IsWeekend 判断是否为周末
func IsWeekend(t time.Time) bool {
    weekday := t.Weekday()
    return weekday == time.Saturday || weekday == time.Sunday
}

// IsWorkday 判断是否为工作日
func IsWorkday(t time.Time) bool {
    return !IsWeekend(t)
}

// GetAge 根据生日计算年龄
func GetAge(birthday time.Time) int {
    now := time.Now()
    age := now.Year() - birthday.Year()
    
    if now.Month() < birthday.Month() || (now.Month() == birthday.Month() && now.Day() < birthday.Day()) {
        age--
    }
    
    return age
}

// FormatDuration 格式化时间间隔
func FormatDuration(d time.Duration) string {
    if d < time.Minute {
        return fmt.Sprintf("%.0f秒", d.Seconds())
    } else if d < time.Hour {
        return fmt.Sprintf("%.0f分钟", d.Minutes())
    } else if d < 24*time.Hour {
        return fmt.Sprintf("%.1f小时", d.Hours())
    } else {
        return fmt.Sprintf("%.1f天", d.Hours()/24)
    }
}

// TimeAgo 时间前描述（如：3分钟前）
func TimeAgo(t time.Time) string {
    now := time.Now()
    diff := now.Sub(t)
    
    if diff < time.Minute {
        return "刚刚"
    } else if diff < time.Hour {
        return fmt.Sprintf("%.0f分钟前", diff.Minutes())
    } else if diff < 24*time.Hour {
        return fmt.Sprintf("%.0f小时前", diff.Hours())
    } else if diff < 30*24*time.Hour {
        return fmt.Sprintf("%.0f天前", diff.Hours()/24)
    } else if diff < 365*24*time.Hour {
        return fmt.Sprintf("%.0f个月前", diff.Hours()/(24*30))
    } else {
        return fmt.Sprintf("%.0f年前", diff.Hours()/(24*365))
    }
}

// Sleep 休眠指定时间
func Sleep(duration time.Duration) {
    time.Sleep(duration)
}

// SleepSeconds 休眠指定秒数
func SleepSeconds(seconds int) {
    time.Sleep(time.Duration(seconds) * time.Second)
}

// SleepMilliseconds 休眠指定毫秒数
func SleepMilliseconds(milliseconds int) {
    time.Sleep(time.Duration(milliseconds) * time.Millisecond)
}
```

## 文件操作工具

```go
// file.go
package utils

import (
    "bufio"
    "fmt"
    "io"
    "io/ioutil"
    "os"
    "path/filepath"
    "strings"
)

// FileExists 检查文件是否存在
func FileExists(filename string) bool {
    _, err := os.Stat(filename)
    return !os.IsNotExist(err)
}

// DirExists 检查目录是否存在
func DirExists(dirname string) bool {
    info, err := os.Stat(dirname)
    return !os.IsNotExist(err) && info.IsDir()
}

// CreateDir 创建目录
func CreateDir(dirname string) error {
    return os.MkdirAll(dirname, 0755)
}

// CreateDirIfNotExists 如果目录不存在则创建
func CreateDirIfNotExists(dirname string) error {
    if !DirExists(dirname) {
        return CreateDir(dirname)
    }
    return nil
}

// ReadFile 读取文件内容
func ReadFile(filename string) ([]byte, error) {
    return ioutil.ReadFile(filename)
}

// ReadFileAsString 读取文件内容为字符串
func ReadFileAsString(filename string) (string, error) {
    data, err := ReadFile(filename)
    if err != nil {
        return "", err
    }
    return string(data), nil
}

// WriteFile 写入文件
func WriteFile(filename string, data []byte) error {
    return ioutil.WriteFile(filename, data, 0644)
}

// WriteStringToFile 写入字符串到文件
func WriteStringToFile(filename, content string) error {
    return WriteFile(filename, []byte(content))
}

// AppendToFile 追加内容到文件
func AppendToFile(filename, content string) error {
    file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        return err
    }
    defer file.Close()
    
    _, err = file.WriteString(content)
    return err
}

// CopyFile 复制文件
func CopyFile(src, dst string) error {
    sourceFile, err := os.Open(src)
    if err != nil {
        return err
    }
    defer sourceFile.Close()
    
    destFile, err := os.Create(dst)
    if err != nil {
        return err
    }
    defer destFile.Close()
    
    _, err = io.Copy(destFile, sourceFile)
    if err != nil {
        return err
    }
    
    return destFile.Sync()
}

// MoveFile 移动文件
func MoveFile(src, dst string) error {
    err := CopyFile(src, dst)
    if err != nil {
        return err
    }
    return os.Remove(src)
}

// DeleteFile 删除文件
func DeleteFile(filename string) error {
    return os.Remove(filename)
}

// DeleteDir 删除目录
func DeleteDir(dirname string) error {
    return os.RemoveAll(dirname)
}

// GetFileSize 获取文件大小
func GetFileSize(filename string) (int64, error) {
    info, err := os.Stat(filename)
    if err != nil {
        return 0, err
    }
    return info.Size(), nil
}

// GetFileExtension 获取文件扩展名
func GetFileExtension(filename string) string {
    return filepath.Ext(filename)
}

// GetFileName 获取文件名（不含扩展名）
func GetFileName(filename string) string {
    base := filepath.Base(filename)
    ext := filepath.Ext(base)
    return strings.TrimSuffix(base, ext)
}

// GetFileDir 获取文件所在目录
func GetFileDir(filename string) string {
    return filepath.Dir(filename)
}

// ListFiles 列出目录下的所有文件
func ListFiles(dirname string) ([]string, error) {
    var files []string
    
    err := filepath.Walk(dirname, func(path string, info os.FileInfo, err error) error {
        if err != nil {
            return err
        }
        if !info.IsDir() {
            files = append(files, path)
        }
        return nil
    })
    
    return files, err
}

// ListDirs 列出目录下的所有子目录
func ListDirs(dirname string) ([]string, error) {
    var dirs []string
    
    err := filepath.Walk(dirname, func(path string, info os.FileInfo, err error) error {
        if err != nil {
            return err
        }
        if info.IsDir() && path != dirname {
            dirs = append(dirs, path)
        }
        return nil
    })
    
    return dirs, err
}

// ReadLines 按行读取文件
func ReadLines(filename string) ([]string, error) {
    file, err := os.Open(filename)
    if err != nil {
        return nil, err
    }
    defer file.Close()
    
    var lines []string
    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        lines = append(lines, scanner.Text())
    }
    
    return lines, scanner.Err()
}

// WriteLines 按行写入文件
func WriteLines(filename string, lines []string) error {
    file, err := os.Create(filename)
    if err != nil {
        return err
    }
    defer file.Close()
    
    writer := bufio.NewWriter(file)
    for _, line := range lines {
        _, err := writer.WriteString(line + "\n")
        if err != nil {
            return err
        }
    }
    
    return writer.Flush()
}

// GetTempDir 获取临时目录
func GetTempDir() string {
    return os.TempDir()
}

// CreateTempFile 创建临时文件
func CreateTempFile(prefix string) (*os.File, error) {
    return ioutil.TempFile("", prefix)
}

// CreateTempDir 创建临时目录
func CreateTempDir(prefix string) (string, error) {
    return ioutil.TempDir("", prefix)
}

// FormatFileSize 格式化文件大小
func FormatFileSize(size int64) string {
    const unit = 1024
    if size < unit {
        return fmt.Sprintf("%d B", size)
    }
    
    div, exp := int64(unit), 0
    for n := size / unit; n >= unit; n /= unit {
        div *= unit
        exp++
    }
    
    return fmt.Sprintf("%.1f %cB", float64(size)/float64(div), "KMGTPE"[exp])
}
```

## 加密解密工具

```go
// crypto.go
package utils

import (
    "crypto/aes"
    "crypto/cipher"
    "crypto/md5"
    "crypto/rand"
    "crypto/sha1"
    "crypto/sha256"
    "crypto/sha512"
    "encoding/base64"
    "encoding/hex"
    "fmt"
    "io"
    
    "golang.org/x/crypto/bcrypt"
)

// MD5 计算MD5哈希
func MD5(data string) string {
    hash := md5.Sum([]byte(data))
    return hex.EncodeToString(hash[:])
}

// SHA1 计算SHA1哈希
func SHA1(data string) string {
    hash := sha1.Sum([]byte(data))
    return hex.EncodeToString(hash[:])
}

// SHA256 计算SHA256哈希
func SHA256(data string) string {
    hash := sha256.Sum256([]byte(data))
    return hex.EncodeToString(hash[:])
}

// SHA512 计算SHA512哈希
func SHA512(data string) string {
    hash := sha512.Sum512([]byte(data))
    return hex.EncodeToString(hash[:])
}

// Base64Encode Base64编码
func Base64Encode(data []byte) string {
    return base64.StdEncoding.EncodeToString(data)
}

// Base64Decode Base64解码
func Base64Decode(data string) ([]byte, error) {
    return base64.StdEncoding.DecodeString(data)
}

// Base64URLEncode Base64 URL安全编码
func Base64URLEncode(data []byte) string {
    return base64.URLEncoding.EncodeToString(data)
}

// Base64URLDecode Base64 URL安全解码
func Base64URLDecode(data string) ([]byte, error) {
    return base64.URLEncoding.DecodeString(data)
}

// HashPassword 使用bcrypt加密密码
func HashPassword(password string) (string, error) {
    bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    return string(bytes), err
}

// CheckPassword 验证密码
func CheckPassword(password, hash string) bool {
    err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
    return err == nil
}

// GenerateRandomBytes 生成随机字节
func GenerateRandomBytes(n int) ([]byte, error) {
    b := make([]byte, n)
    _, err := rand.Read(b)
    if err != nil {
        return nil, err
    }
    return b, nil
}

// GenerateRandomString 生成随机字符串
func GenerateRandomString(n int) (string, error) {
    bytes, err := GenerateRandomBytes(n)
    if err != nil {
        return "", err
    }
    return base64.URLEncoding.EncodeToString(bytes)[:n], nil
}

// AESEncrypt AES加密
func AESEncrypt(plaintext, key []byte) ([]byte, error) {
    block, err := aes.NewCipher(key)
    if err != nil {
        return nil, err
    }
    
    gcm, err := cipher.NewGCM(block)
    if err != nil {
        return nil, err
    }
    
    nonce := make([]byte, gcm.NonceSize())
    if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
        return nil, err
    }
    
    ciphertext := gcm.Seal(nonce, nonce, plaintext, nil)
    return ciphertext, nil
}

// AESDecrypt AES解密
func AESDecrypt(ciphertext, key []byte) ([]byte, error) {
    block, err := aes.NewCipher(key)
    if err != nil {
        return nil, err
    }
    
    gcm, err := cipher.NewGCM(block)
    if err != nil {
        return nil, err
    }
    
    nonceSize := gcm.NonceSize()
    if len(ciphertext) < nonceSize {
        return nil, fmt.Errorf("ciphertext too short")
    }
    
    nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
    plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
    if err != nil {
        return nil, err
    }
    
    return plaintext, nil
}

// AESEncryptString AES加密字符串
func AESEncryptString(plaintext, key string) (string, error) {
    keyBytes := []byte(key)
    if len(keyBytes) != 16 && len(keyBytes) != 24 && len(keyBytes) != 32 {
        return "", fmt.Errorf("key length must be 16, 24, or 32 bytes")
    }
    
    ciphertext, err := AESEncrypt([]byte(plaintext), keyBytes)
    if err != nil {
        return "", err
    }
    
    return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// AESDecryptString AES解密字符串
func AESDecryptString(ciphertext, key string) (string, error) {
    keyBytes := []byte(key)
    if len(keyBytes) != 16 && len(keyBytes) != 24 && len(keyBytes) != 32 {
        return "", fmt.Errorf("key length must be 16, 24, or 32 bytes")
    }
    
    ciphertextBytes, err := base64.StdEncoding.DecodeString(ciphertext)
    if err != nil {
        return "", err
    }
    
    plaintext, err := AESDecrypt(ciphertextBytes, keyBytes)
    if err != nil {
        return "", err
    }
    
    return string(plaintext), nil
}
```

## 数据验证工具

```go
// validator.go
package utils

import (
    "net"
    "net/mail"
    "regexp"
    "strconv"
    "strings"
    "unicode"
)

// IsEmail 验证邮箱格式
func IsEmail(email string) bool {
    _, err := mail.ParseAddress(email)
    return err == nil
}

// IsPhone 验证手机号码（中国大陆）
func IsPhone(phone string) bool {
    pattern := `^1[3-9]\d{9}$`
    matched, _ := regexp.MatchString(pattern, phone)
    return matched
}

// IsURL 验证URL格式
func IsURL(url string) bool {
    pattern := `^https?://[\w\-]+(\.[\w\-]+)+([\w\-\.,@?^=%&:/~\+#]*[\w\-\@?^=%&/~\+#])?$`
    matched, _ := regexp.MatchString(pattern, url)
    return matched
}

// IsIP 验证IP地址
func IsIP(ip string) bool {
    return net.ParseIP(ip) != nil
}

// IsIPv4 验证IPv4地址
func IsIPv4(ip string) bool {
    parsedIP := net.ParseIP(ip)
    return parsedIP != nil && parsedIP.To4() != nil
}

// IsIPv6 验证IPv6地址
func IsIPv6(ip string) bool {
    parsedIP := net.ParseIP(ip)
    return parsedIP != nil && parsedIP.To4() == nil
}

// IsIDCard 验证身份证号码（中国大陆）
func IsIDCard(idCard string) bool {
    if len(idCard) != 18 {
        return false
    }
    
    // 验证前17位是否为数字
    for i := 0; i < 17; i++ {
        if !unicode.IsDigit(rune(idCard[i])) {
            return false
        }
    }
    
    // 验证最后一位校验码
    weights := []int{7, 9, 10, 5, 8, 4, 2, 1, 6, 3, 7, 9, 10, 5, 8, 4, 2}
    checkCodes := []string{"1", "0", "X", "9", "8", "7", "6", "5", "4", "3", "2"}
    
    sum := 0
    for i := 0; i < 17; i++ {
        digit, _ := strconv.Atoi(string(idCard[i]))
        sum += digit * weights[i]
    }
    
    checkCode := checkCodes[sum%11]
    return string(idCard[17]) == checkCode
}

// IsNumeric 验证是否为数字
func IsNumeric(s string) bool {
    _, err := strconv.ParseFloat(s, 64)
    return err == nil
}

// IsInteger 验证是否为整数
func IsInteger(s string) bool {
    _, err := strconv.Atoi(s)
    return err == nil
}

// IsAlpha 验证是否只包含字母
func IsAlpha(s string) bool {
    for _, r := range s {
        if !unicode.IsLetter(r) {
            return false
        }
    }
    return len(s) > 0
}

// IsAlphaNumeric 验证是否只包含字母和数字
func IsAlphaNumeric(s string) bool {
    for _, r := range s {
        if !unicode.IsLetter(r) && !unicode.IsDigit(r) {
            return false
        }
    }
    return len(s) > 0
}

// IsStrongPassword 验证是否为强密码
func IsStrongPassword(password string) bool {
    if len(password) < 8 {
        return false
    }
    
    hasUpper := false
    hasLower := false
    hasDigit := false
    hasSpecial := false
    
    for _, r := range password {
        switch {
        case unicode.IsUpper(r):
            hasUpper = true
        case unicode.IsLower(r):
            hasLower = true
        case unicode.IsDigit(r):
            hasDigit = true
        case unicode.IsPunct(r) || unicode.IsSymbol(r):
            hasSpecial = true
        }
    }
    
    return hasUpper && hasLower && hasDigit && hasSpecial
}

// IsCreditCard 验证信用卡号码（Luhn算法）
func IsCreditCard(cardNumber string) bool {
    // 移除空格和连字符
    cardNumber = strings.ReplaceAll(cardNumber, " ", "")
    cardNumber = strings.ReplaceAll(cardNumber, "-", "")
    
    if len(cardNumber) < 13 || len(cardNumber) > 19 {
        return false
    }
    
    // 验证是否全为数字
    for _, r := range cardNumber {
        if !unicode.IsDigit(r) {
            return false
        }
    }
    
    // Luhn算法验证
    sum := 0
    alternate := false
    
    for i := len(cardNumber) - 1; i >= 0; i-- {
        digit, _ := strconv.Atoi(string(cardNumber[i]))
        
        if alternate {
            digit *= 2
            if digit > 9 {
                digit = digit%10 + digit/10
            }
        }
        
        sum += digit
        alternate = !alternate
    }
    
    return sum%10 == 0
}

// IsJSON 验证是否为有效的JSON格式
func IsJSON(s string) bool {
    var js interface{}
    return json.Unmarshal([]byte(s), &js) == nil
}

// IsBase64 验证是否为Base64编码
func IsBase64(s string) bool {
    _, err := base64.StdEncoding.DecodeString(s)
    return err == nil
}

// IsHexColor 验证是否为十六进制颜色值
func IsHexColor(color string) bool {
    pattern := `^#([A-Fa-f0-9]{6}|[A-Fa-f0-9]{3})$`
    matched, _ := regexp.MatchString(pattern, color)
    return matched
}

// IsMAC 验证MAC地址
func IsMAC(mac string) bool {
    _, err := net.ParseMAC(mac)
    return err == nil
}

// IsPort 验证端口号
func IsPort(port string) bool {
    p, err := strconv.Atoi(port)
    return err == nil && p >= 1 && p <= 65535
}

// IsLatitude 验证纬度
func IsLatitude(lat string) bool {
    latitude, err := strconv.ParseFloat(lat, 64)
    return err == nil && latitude >= -90 && latitude <= 90
}

// IsLongitude 验证经度
func IsLongitude(lng string) bool {
    longitude, err := strconv.ParseFloat(lng, 64)
    return err == nil && longitude >= -180 && longitude <= 180
}
```

## 分页工具

```go
// pagination.go
package utils

import "math"

// Pagination 分页结构
type Pagination struct {
    Page       int   `json:"page"`        // 当前页码
    PageSize   int   `json:"page_size"`   // 每页大小
    Total      int64 `json:"total"`       // 总记录数
    TotalPages int   `json:"total_pages"` // 总页数
    HasNext    bool  `json:"has_next"`    // 是否有下一页
    HasPrev    bool  `json:"has_prev"`    // 是否有上一页
    Offset     int   `json:"offset"`      // 偏移量
}

// NewPagination 创建分页对象
func NewPagination(page, pageSize int, total int64) *Pagination {
    if page < 1 {
        page = 1
    }
    if pageSize < 1 {
        pageSize = 10
    }
    
    totalPages := int(math.Ceil(float64(total) / float64(pageSize)))
    if totalPages < 1 {
        totalPages = 1
    }
    
    if page > totalPages {
        page = totalPages
    }
    
    offset := (page - 1) * pageSize
    
    return &Pagination{
        Page:       page,
        PageSize:   pageSize,
        Total:      total,
        TotalPages: totalPages,
        HasNext:    page < totalPages,
        HasPrev:    page > 1,
        Offset:     offset,
    }
}

// GetOffset 获取偏移量
func (p *Pagination) GetOffset() int {
    return p.Offset
}

// GetLimit 获取限制数量
func (p *Pagination) GetLimit() int {
    return p.PageSize
}

// GetNextPage 获取下一页页码
func (p *Pagination) GetNextPage() int {
    if p.HasNext {
        return p.Page + 1
    }
    return p.Page
}

// GetPrevPage 获取上一页页码
func (p *Pagination) GetPrevPage() int {
    if p.HasPrev {
        return p.Page - 1
    }
    return p.Page
}

// GetPageNumbers 获取页码列表（用于分页导航）
func (p *Pagination) GetPageNumbers(maxPages int) []int {
    if maxPages <= 0 {
        maxPages = 10
    }
    
    var pages []int
    
    if p.TotalPages <= maxPages {
        // 总页数不超过最大显示页数，显示所有页码
        for i := 1; i <= p.TotalPages; i++ {
            pages = append(pages, i)
        }
    } else {
        // 总页数超过最大显示页数，计算显示范围
        half := maxPages / 2
        start := p.Page - half
        end := p.Page + half
        
        if start < 1 {
            start = 1
            end = maxPages
        }
        
        if end > p.TotalPages {
            end = p.TotalPages
            start = p.TotalPages - maxPages + 1
            if start < 1 {
                start = 1
            }
        }
        
        for i := start; i <= end; i++ {
            pages = append(pages, i)
        }
    }
    
    return pages
}

// IsValidPage 验证页码是否有效
func (p *Pagination) IsValidPage(page int) bool {
    return page >= 1 && page <= p.TotalPages
}

// GetStartRecord 获取当前页第一条记录的序号
func (p *Pagination) GetStartRecord() int {
    if p.Total == 0 {
        return 0
    }
    return p.Offset + 1
}

// GetEndRecord 获取当前页最后一条记录的序号
func (p *Pagination) GetEndRecord() int {
    end := p.Offset + p.PageSize
    if end > int(p.Total) {
        end = int(p.Total)
    }
    return end
}

// GetRecordRange 获取当前页记录范围描述
func (p *Pagination) GetRecordRange() string {
    if p.Total == 0 {
        return "0 - 0 of 0"
    }
    return fmt.Sprintf("%d - %d of %d", p.GetStartRecord(), p.GetEndRecord(), p.Total)
}
```

## 使用示例

### 字符串处理示例

```go
package main

import (
    "fmt"
    "github.com/your-project/pkg/utils"
)

func main() {
    // 字符串操作
    fmt.Println(utils.IsEmpty(""))           // true
    fmt.Println(utils.Capitalize("hello"))   // "Hello"
    fmt.Println(utils.Reverse("hello"))      // "olleh"
    fmt.Println(utils.CamelCase("hello_world")) // "helloWorld"
    fmt.Println(utils.SnakeCase("HelloWorld"))  // "hello_world"
    
    // 敏感信息掩码
    email := "user@example.com"
    fmt.Println(utils.MaskEmail(email))      // "u***r@example.com"
    
    phone := "13812345678"
    fmt.Println(utils.MaskPhone(phone))      // "138****5678"
}
```

### 时间处理示例

```go
func main() {
    now := time.Now()
    
    // 时间格式化
    fmt.Println(utils.FormatTime(now, utils.DateTimeFormat))
    
    // 时间计算
    startOfDay := utils.GetStartOfDay(now)
    endOfDay := utils.GetEndOfDay(now)
    fmt.Printf("今天: %s - %s\n", 
        utils.FormatTime(startOfDay, utils.TimeFormat),
        utils.FormatTime(endOfDay, utils.TimeFormat))
    
    // 时间差计算
    yesterday := now.AddDate(0, 0, -1)
    fmt.Printf("相差天数: %d\n", utils.DiffDays(yesterday, now))
    
    // 时间描述
    fmt.Println(utils.TimeAgo(yesterday)) // "1天前"
}
```

### 文件操作示例

```go
func main() {
    filename := "test.txt"
    content := "Hello, World!"
    
    // 写入文件
    err := utils.WriteStringToFile(filename, content)
    if err != nil {
        fmt.Printf("写入文件失败: %v\n", err)
        return
    }
    
    // 读取文件
    data, err := utils.ReadFileAsString(filename)
    if err != nil {
        fmt.Printf("读取文件失败: %v\n", err)
        return
    }
    fmt.Printf("文件内容: %s\n", data)
    
    // 获取文件信息
    size, _ := utils.GetFileSize(filename)
    fmt.Printf("文件大小: %s\n", utils.FormatFileSize(size))
    
    // 清理
    utils.DeleteFile(filename)
}
```

### 加密解密示例

```go
func main() {
    password := "mypassword"
    
    // 密码加密
    hashedPassword, err := utils.HashPassword(password)
    if err != nil {
        fmt.Printf("密码加密失败: %v\n", err)
        return
    }
    
    // 密码验证
    isValid := utils.CheckPassword(password, hashedPassword)
    fmt.Printf("密码验证: %t\n", isValid)
    
    // AES加密
    key := "1234567890123456" // 16字节密钥
    plaintext := "Hello, World!"
    
    encrypted, err := utils.AESEncryptString(plaintext, key)
    if err != nil {
        fmt.Printf("AES加密失败: %v\n", err)
        return
    }
    fmt.Printf("加密结果: %s\n", encrypted)
    
    // AES解密
    decrypted, err := utils.AESDecryptString(encrypted, key)
    if err != nil {
        fmt.Printf("AES解密失败: %v\n", err)
        return
    }
    fmt.Printf("解密结果: %s\n", decrypted)
}
```

### 数据验证示例

```go
func main() {
    // 邮箱验证
    email := "user@example.com"
    fmt.Printf("邮箱格式: %t\n", utils.IsEmail(email))
    
    // 手机号验证
    phone := "13812345678"
    fmt.Printf("手机号格式: %t\n", utils.IsPhone(phone))
    
    // 强密码验证
    password := "MyPassword123!"
    fmt.Printf("强密码: %t\n", utils.IsStrongPassword(password))
    
    // URL验证
    url := "https://www.example.com"
    fmt.Printf("URL格式: %t\n", utils.IsURL(url))
    
    // IP地址验证
    ip := "192.168.1.1"
    fmt.Printf("IP地址: %t\n", utils.IsIP(ip))
}
```

### 分页示例

```go
func main() {
    // 创建分页对象
    total := int64(1000)
    page := 5
    pageSize := 20
    
    pagination := utils.NewPagination(page, pageSize, total)
    
    fmt.Printf("当前页: %d\n", pagination.Page)
    fmt.Printf("每页大小: %d\n", pagination.PageSize)
    fmt.Printf("总记录数: %d\n", pagination.Total)
    fmt.Printf("总页数: %d\n", pagination.TotalPages)
    fmt.Printf("是否有下一页: %t\n", pagination.HasNext)
    fmt.Printf("是否有上一页: %t\n", pagination.HasPrev)
    fmt.Printf("偏移量: %d\n", pagination.GetOffset())
    fmt.Printf("记录范围: %s\n", pagination.GetRecordRange())
    
    // 获取页码列表
    pageNumbers := pagination.GetPageNumbers(10)
    fmt.Printf("页码列表: %v\n", pageNumbers)
}
```

## 性能优化

### 1. 字符串操作优化

```go
// 使用 strings.Builder 进行字符串拼接
func ConcatStrings(strs []string) string {
    var builder strings.Builder
    for _, str := range strs {
        builder.WriteString(str)
    }
    return builder.String()
}

// 使用缓存避免重复计算
var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

func IsEmailOptimized(email string) bool {
    return emailRegex.MatchString(email)
}
```

### 2. 文件操作优化

```go
// 使用缓冲读写提高性能
func ReadLargeFile(filename string) error {
    file, err := os.Open(filename)
    if err != nil {
        return err
    }
    defer file.Close()
    
    reader := bufio.NewReader(file)
    buffer := make([]byte, 4096)
    
    for {
        n, err := reader.Read(buffer)
        if err != nil {
            if err == io.EOF {
                break
            }
            return err
        }
        
        // 处理读取的数据
        processData(buffer[:n])
    }
    
    return nil
}
```

### 3. 加密操作优化

```go
// 使用对象池减少内存分配
var hashPool = sync.Pool{
    New: func() interface{} {
        return sha256.New()
    },
}

func SHA256Optimized(data string) string {
    hasher := hashPool.Get().(hash.Hash)
    defer hashPool.Put(hasher)
    
    hasher.Reset()
    hasher.Write([]byte(data))
    return hex.EncodeToString(hasher.Sum(nil))
}
```

## 测试

### 单元测试示例

```go
// string_test.go
package utils

import (
    "testing"
    "github.com/stretchr/testify/assert"
)

func TestIsEmpty(t *testing.T) {
    tests := []struct {
        input    string
        expected bool
    }{
        {"", true},
        {" ", true},
        {"\t\n", true},
        {"hello", false},
        {" hello ", false},
    }
    
    for _, test := range tests {
        result := IsEmpty(test.input)
        assert.Equal(t, test.expected, result, "IsEmpty(%q) should be %t", test.input, test.expected)
    }
}

func TestCapitalize(t *testing.T) {
    tests := []struct {
        input    string
        expected string
    }{
        {"", ""},
        {"hello", "Hello"},
        {"HELLO", "HELLO"},
        {"hELLO", "HELLO"},
        {"中文", "中文"},
    }
    
    for _, test := range tests {
        result := Capitalize(test.input)
        assert.Equal(t, test.expected, result, "Capitalize(%q) should be %q", test.input, test.expected)
    }
}

func BenchmarkMD5(b *testing.B) {
    data := "Hello, World!"
    b.ResetTimer()
    
    for i := 0; i < b.N; i++ {
        MD5(data)
    }
}
```

## 最佳实践

### 1. 错误处理

```go
// 统一错误处理
func SafeOperation(operation func() error) error {
    defer func() {
        if r := recover(); r != nil {
            klogger.Error("Operation panic", zap.Any("error", r))
        }
    }()
    
    return operation()
}
```

### 2. 参数验证

```go
// 参数验证
func ValidateInput(input string) error {
    if IsEmpty(input) {
        return fmt.Errorf("input cannot be empty")
    }
    
    if len(input) > 1000 {
        return fmt.Errorf("input too long")
    }
    
    return nil
}
```

### 3. 资源管理

```go
// 确保资源释放
func ProcessFile(filename string) error {
    file, err := os.Open(filename)
    if err != nil {
        return err
    }
    defer file.Close() // 确保文件关闭
    
    // 处理文件
    return nil
}
```

## 相关文档

- [JWT 认证系统文档](JWT_AUTH.md)
- [中间件系统文档](MIDDLEWARE.md)
- [统一响应格式文档](RESPONSE.md)
- [Redis 缓存系统文档](REDIS.md)
- [Go 官方文档](https://golang.org/doc/)

---

**最佳实践**: 使用工具函数前进行参数验证；处理敏感数据时使用安全的加密算法；文件操作时注意资源释放；字符串操作时考虑性能优化；为工具函数编写完整的单元测试；遵循单一职责原则设计工具函数。