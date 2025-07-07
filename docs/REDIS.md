# Redis 缓存系统文档

## 概述

本项目集成了完整的 Redis 缓存系统，基于 go-redis 客户端库构建，提供了键值对操作、哈希操作、列表操作、健康检查等功能。Redis 系统支持连接池管理、自动重连、性能监控和故障恢复。

## 系统架构

### Redis 客户端结构

```
Redis Client
├── Connection Pool          # 连接池管理
├── Key-Value Operations     # 键值对操作
├── Hash Operations          # 哈希操作
├── List Operations          # 列表操作
├── Set Operations           # 集合操作
├── Sorted Set Operations    # 有序集合操作
├── Pub/Sub Operations       # 发布订阅
├── Transaction Support      # 事务支持
├── Pipeline Support         # 管道支持
└── Health Check            # 健康检查
```

## 配置管理

### Redis 配置结构

```go
type RedisConfig struct {
    Host         string        `mapstructure:"host"`          // Redis 主机地址
    Port         int           `mapstructure:"port"`          // Redis 端口
    Password     string        `mapstructure:"password"`      // Redis 密码
    DB           int           `mapstructure:"db"`            // 数据库编号
    PoolSize     int           `mapstructure:"pool_size"`     // 连接池大小
    MinIdleConns int           `mapstructure:"min_idle_conns"` // 最小空闲连接数
    MaxRetries   int           `mapstructure:"max_retries"`   // 最大重试次数
    DialTimeout  time.Duration `mapstructure:"dial_timeout"`  // 连接超时
    ReadTimeout  time.Duration `mapstructure:"read_timeout"`  // 读取超时
    WriteTimeout time.Duration `mapstructure:"write_timeout"` // 写入超时
    IdleTimeout  time.Duration `mapstructure:"idle_timeout"`  // 空闲超时
}
```

### 配置示例

```yaml
# config/config.yaml
redis:
  host: "localhost"
  port: 6379
  password: ""
  db: 0
  pool_size: 10
  min_idle_conns: 5
  max_retries: 3
  dial_timeout: 5s
  read_timeout: 3s
  write_timeout: 3s
  idle_timeout: 300s
```

## 核心功能

### 客户端初始化

```go
var (
    client *redis.Client
    once   sync.Once
)

// Init 初始化 Redis 客户端
func Init(cfg *config.RedisConfig) error {
    var err error
    once.Do(func() {
        client = redis.NewClient(&redis.Options{
            Addr:         fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
            Password:     cfg.Password,
            DB:           cfg.DB,
            PoolSize:     cfg.PoolSize,
            MinIdleConns: cfg.MinIdleConns,
            MaxRetries:   cfg.MaxRetries,
            DialTimeout:  cfg.DialTimeout,
            ReadTimeout:  cfg.ReadTimeout,
            WriteTimeout: cfg.WriteTimeout,
            IdleTimeout:  cfg.IdleTimeout,
        })
        
        // 测试连接
        ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
        defer cancel()
        
        _, err = client.Ping(ctx).Result()
        if err != nil {
            klogger.Error("Failed to connect to Redis", zap.Error(err))
            return
        }
        
        klogger.Info("Redis client initialized successfully")
    })
    
    return err
}
```

### 客户端管理

```go
// GetClient 获取 Redis 客户端实例
func GetClient() *redis.Client {
    if client == nil {
        klogger.Error("Redis client not initialized")
        return nil
    }
    return client
}

// Close 关闭 Redis 客户端
func Close() error {
    if client != nil {
        err := client.Close()
        if err != nil {
            klogger.Error("Failed to close Redis client", zap.Error(err))
            return err
        }
        klogger.Info("Redis client closed successfully")
    }
    return nil
}
```

## 键值对操作

### 基础操作

```go
// Set 设置键值对
func Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
    if client == nil {
        return errors.New("Redis client not initialized")
    }
    
    err := client.Set(ctx, key, value, expiration).Err()
    if err != nil {
        klogger.Error("Failed to set key",
            zap.String("key", key),
            zap.Error(err),
        )
        return err
    }
    
    klogger.Debug("Key set successfully",
        zap.String("key", key),
        zap.Duration("expiration", expiration),
    )
    return nil
}

// Get 获取键值
func Get(ctx context.Context, key string) (string, error) {
    if client == nil {
        return "", errors.New("Redis client not initialized")
    }
    
    val, err := client.Get(ctx, key).Result()
    if err != nil {
        if err == redis.Nil {
            klogger.Debug("Key not found", zap.String("key", key))
            return "", nil
        }
        klogger.Error("Failed to get key",
            zap.String("key", key),
            zap.Error(err),
        )
        return "", err
    }
    
    return val, nil
}

// Del 删除键
func Del(ctx context.Context, keys ...string) error {
    if client == nil {
        return errors.New("Redis client not initialized")
    }
    
    err := client.Del(ctx, keys...).Err()
    if err != nil {
        klogger.Error("Failed to delete keys",
            zap.Strings("keys", keys),
            zap.Error(err),
        )
        return err
    }
    
    klogger.Debug("Keys deleted successfully", zap.Strings("keys", keys))
    return nil
}

// Exists 检查键是否存在
func Exists(ctx context.Context, keys ...string) (int64, error) {
    if client == nil {
        return 0, errors.New("Redis client not initialized")
    }
    
    count, err := client.Exists(ctx, keys...).Result()
    if err != nil {
        klogger.Error("Failed to check key existence",
            zap.Strings("keys", keys),
            zap.Error(err),
        )
        return 0, err
    }
    
    return count, nil
}
```

### 过期时间操作

```go
// Expire 设置键的过期时间
func Expire(ctx context.Context, key string, expiration time.Duration) error {
    if client == nil {
        return errors.New("Redis client not initialized")
    }
    
    err := client.Expire(ctx, key, expiration).Err()
    if err != nil {
        klogger.Error("Failed to set expiration",
            zap.String("key", key),
            zap.Duration("expiration", expiration),
            zap.Error(err),
        )
        return err
    }
    
    return nil
}

// TTL 获取键的剩余生存时间
func TTL(ctx context.Context, key string) (time.Duration, error) {
    if client == nil {
        return 0, errors.New("Redis client not initialized")
    }
    
    ttl, err := client.TTL(ctx, key).Result()
    if err != nil {
        klogger.Error("Failed to get TTL",
            zap.String("key", key),
            zap.Error(err),
        )
        return 0, err
    }
    
    return ttl, nil
}
```

## 哈希操作

### 哈希字段操作

```go
// HSet 设置哈希字段
func HSet(ctx context.Context, key string, values ...interface{}) error {
    if client == nil {
        return errors.New("Redis client not initialized")
    }
    
    err := client.HSet(ctx, key, values...).Err()
    if err != nil {
        klogger.Error("Failed to set hash field",
            zap.String("key", key),
            zap.Error(err),
        )
        return err
    }
    
    return nil
}

// HGet 获取哈希字段值
func HGet(ctx context.Context, key, field string) (string, error) {
    if client == nil {
        return "", errors.New("Redis client not initialized")
    }
    
    val, err := client.HGet(ctx, key, field).Result()
    if err != nil {
        if err == redis.Nil {
            return "", nil
        }
        klogger.Error("Failed to get hash field",
            zap.String("key", key),
            zap.String("field", field),
            zap.Error(err),
        )
        return "", err
    }
    
    return val, nil
}

// HGetAll 获取哈希的所有字段和值
func HGetAll(ctx context.Context, key string) (map[string]string, error) {
    if client == nil {
        return nil, errors.New("Redis client not initialized")
    }
    
    result, err := client.HGetAll(ctx, key).Result()
    if err != nil {
        klogger.Error("Failed to get all hash fields",
            zap.String("key", key),
            zap.Error(err),
        )
        return nil, err
    }
    
    return result, nil
}

// HDel 删除哈希字段
func HDel(ctx context.Context, key string, fields ...string) error {
    if client == nil {
        return errors.New("Redis client not initialized")
    }
    
    err := client.HDel(ctx, key, fields...).Err()
    if err != nil {
        klogger.Error("Failed to delete hash fields",
            zap.String("key", key),
            zap.Strings("fields", fields),
            zap.Error(err),
        )
        return err
    }
    
    return nil
}
```

## 列表操作

### 列表基础操作

```go
// LPush 从列表左侧推入元素
func LPush(ctx context.Context, key string, values ...interface{}) error {
    if client == nil {
        return errors.New("Redis client not initialized")
    }
    
    err := client.LPush(ctx, key, values...).Err()
    if err != nil {
        klogger.Error("Failed to left push to list",
            zap.String("key", key),
            zap.Error(err),
        )
        return err
    }
    
    return nil
}

// RPush 从列表右侧推入元素
func RPush(ctx context.Context, key string, values ...interface{}) error {
    if client == nil {
        return errors.New("Redis client not initialized")
    }
    
    err := client.RPush(ctx, key, values...).Err()
    if err != nil {
        klogger.Error("Failed to right push to list",
            zap.String("key", key),
            zap.Error(err),
        )
        return err
    }
    
    return nil
}

// LPop 从列表左侧弹出元素
func LPop(ctx context.Context, key string) (string, error) {
    if client == nil {
        return "", errors.New("Redis client not initialized")
    }
    
    val, err := client.LPop(ctx, key).Result()
    if err != nil {
        if err == redis.Nil {
            return "", nil
        }
        klogger.Error("Failed to left pop from list",
            zap.String("key", key),
            zap.Error(err),
        )
        return "", err
    }
    
    return val, nil
}

// RPop 从列表右侧弹出元素
func RPop(ctx context.Context, key string) (string, error) {
    if client == nil {
        return "", errors.New("Redis client not initialized")
    }
    
    val, err := client.RPop(ctx, key).Result()
    if err != nil {
        if err == redis.Nil {
            return "", nil
        }
        klogger.Error("Failed to right pop from list",
            zap.String("key", key),
            zap.Error(err),
        )
        return "", err
    }
    
    return val, nil
}

// LLen 获取列表长度
func LLen(ctx context.Context, key string) (int64, error) {
    if client == nil {
        return 0, errors.New("Redis client not initialized")
    }
    
    length, err := client.LLen(ctx, key).Result()
    if err != nil {
        klogger.Error("Failed to get list length",
            zap.String("key", key),
            zap.Error(err),
        )
        return 0, err
    }
    
    return length, nil
}
```

## 高级功能

### 分布式锁

```go
// DistributedLock 分布式锁结构
type DistributedLock struct {
    key        string
    value      string
    expiration time.Duration
    client     *redis.Client
}

// NewDistributedLock 创建分布式锁
func NewDistributedLock(key string, expiration time.Duration) *DistributedLock {
    return &DistributedLock{
        key:        key,
        value:      generateLockValue(),
        expiration: expiration,
        client:     GetClient(),
    }
}

// Lock 获取锁
func (l *DistributedLock) Lock(ctx context.Context) (bool, error) {
    if l.client == nil {
        return false, errors.New("Redis client not initialized")
    }
    
    // 使用 SET NX EX 命令实现原子性锁获取
    result, err := l.client.SetNX(ctx, l.key, l.value, l.expiration).Result()
    if err != nil {
        klogger.Error("Failed to acquire lock",
            zap.String("key", l.key),
            zap.Error(err),
        )
        return false, err
    }
    
    if result {
        klogger.Debug("Lock acquired", zap.String("key", l.key))
    } else {
        klogger.Debug("Lock already exists", zap.String("key", l.key))
    }
    
    return result, nil
}

// Unlock 释放锁
func (l *DistributedLock) Unlock(ctx context.Context) error {
    if l.client == nil {
        return errors.New("Redis client not initialized")
    }
    
    // 使用 Lua 脚本确保原子性释放
    script := `
        if redis.call("get", KEYS[1]) == ARGV[1] then
            return redis.call("del", KEYS[1])
        else
            return 0
        end
    `
    
    result, err := l.client.Eval(ctx, script, []string{l.key}, l.value).Result()
    if err != nil {
        klogger.Error("Failed to release lock",
            zap.String("key", l.key),
            zap.Error(err),
        )
        return err
    }
    
    if result.(int64) == 1 {
        klogger.Debug("Lock released", zap.String("key", l.key))
    } else {
        klogger.Warn("Lock not owned by this instance", zap.String("key", l.key))
    }
    
    return nil
}

// generateLockValue 生成锁值
func generateLockValue() string {
    return fmt.Sprintf("%d-%s", time.Now().UnixNano(), uuid.New().String())
}
```

### 缓存管理器

```go
// CacheManager 缓存管理器
type CacheManager struct {
    client     *redis.Client
    defaultTTL time.Duration
    prefix     string
}

// NewCacheManager 创建缓存管理器
func NewCacheManager(prefix string, defaultTTL time.Duration) *CacheManager {
    return &CacheManager{
        client:     GetClient(),
        defaultTTL: defaultTTL,
        prefix:     prefix,
    }
}

// buildKey 构建缓存键
func (cm *CacheManager) buildKey(key string) string {
    return fmt.Sprintf("%s:%s", cm.prefix, key)
}

// SetJSON 设置 JSON 对象到缓存
func (cm *CacheManager) SetJSON(ctx context.Context, key string, value interface{}, ttl ...time.Duration) error {
    data, err := json.Marshal(value)
    if err != nil {
        return fmt.Errorf("failed to marshal JSON: %w", err)
    }
    
    expiration := cm.defaultTTL
    if len(ttl) > 0 {
        expiration = ttl[0]
    }
    
    return Set(ctx, cm.buildKey(key), data, expiration)
}

// GetJSON 从缓存获取 JSON 对象
func (cm *CacheManager) GetJSON(ctx context.Context, key string, dest interface{}) error {
    data, err := Get(ctx, cm.buildKey(key))
    if err != nil {
        return err
    }
    
    if data == "" {
        return redis.Nil
    }
    
    return json.Unmarshal([]byte(data), dest)
}

// Delete 删除缓存
func (cm *CacheManager) Delete(ctx context.Context, keys ...string) error {
    fullKeys := make([]string, len(keys))
    for i, key := range keys {
        fullKeys[i] = cm.buildKey(key)
    }
    return Del(ctx, fullKeys...)
}

// Clear 清空指定前缀的所有缓存
func (cm *CacheManager) Clear(ctx context.Context) error {
    pattern := cm.prefix + ":*"
    keys, err := cm.client.Keys(ctx, pattern).Result()
    if err != nil {
        return err
    }
    
    if len(keys) > 0 {
        return cm.client.Del(ctx, keys...).Err()
    }
    
    return nil
}
```

### 发布订阅

```go
// Publisher 发布者
type Publisher struct {
    client *redis.Client
}

// NewPublisher 创建发布者
func NewPublisher() *Publisher {
    return &Publisher{
        client: GetClient(),
    }
}

// Publish 发布消息
func (p *Publisher) Publish(ctx context.Context, channel string, message interface{}) error {
    if p.client == nil {
        return errors.New("Redis client not initialized")
    }
    
    data, err := json.Marshal(message)
    if err != nil {
        return fmt.Errorf("failed to marshal message: %w", err)
    }
    
    err = p.client.Publish(ctx, channel, data).Err()
    if err != nil {
        klogger.Error("Failed to publish message",
            zap.String("channel", channel),
            zap.Error(err),
        )
        return err
    }
    
    klogger.Debug("Message published",
        zap.String("channel", channel),
        zap.ByteString("message", data),
    )
    return nil
}

// Subscriber 订阅者
type Subscriber struct {
    client *redis.Client
    pubsub *redis.PubSub
}

// NewSubscriber 创建订阅者
func NewSubscriber(channels ...string) *Subscriber {
    client := GetClient()
    if client == nil {
        return nil
    }
    
    pubsub := client.Subscribe(context.Background(), channels...)
    return &Subscriber{
        client: client,
        pubsub: pubsub,
    }
}

// Listen 监听消息
func (s *Subscriber) Listen(ctx context.Context, handler func(channel, message string) error) error {
    if s.pubsub == nil {
        return errors.New("PubSub not initialized")
    }
    
    defer s.pubsub.Close()
    
    ch := s.pubsub.Channel()
    for {
        select {
        case msg := <-ch:
            if msg == nil {
                return nil
            }
            
            err := handler(msg.Channel, msg.Payload)
            if err != nil {
                klogger.Error("Message handler error",
                    zap.String("channel", msg.Channel),
                    zap.Error(err),
                )
            }
            
        case <-ctx.Done():
            return ctx.Err()
        }
    }
}
```

## 健康检查

```go
// HealthCheck 健康检查
func HealthCheck(ctx context.Context) error {
    if client == nil {
        return errors.New("Redis client not initialized")
    }
    
    // 执行 PING 命令
    _, err := client.Ping(ctx).Result()
    if err != nil {
        klogger.Error("Redis health check failed", zap.Error(err))
        return err
    }
    
    // 检查连接池状态
    stats := client.PoolStats()
    klogger.Debug("Redis pool stats",
        zap.Uint32("hits", stats.Hits),
        zap.Uint32("misses", stats.Misses),
        zap.Uint32("timeouts", stats.Timeouts),
        zap.Uint32("total_conns", stats.TotalConns),
        zap.Uint32("idle_conns", stats.IdleConns),
        zap.Uint32("stale_conns", stats.StaleConns),
    )
    
    return nil
}

// GetStats 获取 Redis 统计信息
func GetStats() map[string]interface{} {
    if client == nil {
        return map[string]interface{}{
            "status": "not_initialized",
        }
    }
    
    stats := client.PoolStats()
    return map[string]interface{}{
        "status":      "connected",
        "hits":        stats.Hits,
        "misses":      stats.Misses,
        "timeouts":    stats.Timeouts,
        "total_conns": stats.TotalConns,
        "idle_conns":  stats.IdleConns,
        "stale_conns": stats.StaleConns,
    }
}
```

## 使用示例

### 基础使用

```go
func ExampleBasicUsage() {
    ctx := context.Background()
    
    // 设置键值对
    err := Set(ctx, "user:1001", "John Doe", time.Hour)
    if err != nil {
        log.Printf("Set error: %v", err)
        return
    }
    
    // 获取值
    value, err := Get(ctx, "user:1001")
    if err != nil {
        log.Printf("Get error: %v", err)
        return
    }
    fmt.Printf("User: %s\n", value)
    
    // 设置哈希
    err = HSet(ctx, "user:1001:profile", "name", "John Doe", "email", "john@example.com")
    if err != nil {
        log.Printf("HSet error: %v", err)
        return
    }
    
    // 获取哈希字段
    email, err := HGet(ctx, "user:1001:profile", "email")
    if err != nil {
        log.Printf("HGet error: %v", err)
        return
    }
    fmt.Printf("Email: %s\n", email)
}
```

### 缓存管理器使用

```go
type User struct {
    ID    uint   `json:"id"`
    Name  string `json:"name"`
    Email string `json:"email"`
}

func ExampleCacheManager() {
    ctx := context.Background()
    
    // 创建用户缓存管理器
    userCache := NewCacheManager("user", time.Hour)
    
    // 缓存用户对象
    user := &User{
        ID:    1001,
        Name:  "John Doe",
        Email: "john@example.com",
    }
    
    err := userCache.SetJSON(ctx, "1001", user)
    if err != nil {
        log.Printf("Cache set error: %v", err)
        return
    }
    
    // 从缓存获取用户对象
    var cachedUser User
    err = userCache.GetJSON(ctx, "1001", &cachedUser)
    if err != nil {
        if err == redis.Nil {
            log.Println("User not found in cache")
        } else {
            log.Printf("Cache get error: %v", err)
        }
        return
    }
    
    fmt.Printf("Cached user: %+v\n", cachedUser)
}
```

### 分布式锁使用

```go
func ExampleDistributedLock() {
    ctx := context.Background()
    
    // 创建分布式锁
    lock := NewDistributedLock("resource:1001", time.Minute)
    
    // 尝试获取锁
    acquired, err := lock.Lock(ctx)
    if err != nil {
        log.Printf("Lock error: %v", err)
        return
    }
    
    if !acquired {
        log.Println("Failed to acquire lock")
        return
    }
    
    // 执行临界区代码
    defer func() {
        err := lock.Unlock(ctx)
        if err != nil {
            log.Printf("Unlock error: %v", err)
        }
    }()
    
    // 模拟业务处理
    time.Sleep(5 * time.Second)
    log.Println("Critical section completed")
}
```

## 性能优化

### 连接池优化

```go
// OptimizeConnectionPool 优化连接池配置
func OptimizeConnectionPool(cfg *config.RedisConfig) {
    // 根据并发量调整连接池大小
    if cfg.PoolSize == 0 {
        cfg.PoolSize = runtime.NumCPU() * 2
    }
    
    // 设置最小空闲连接数
    if cfg.MinIdleConns == 0 {
        cfg.MinIdleConns = cfg.PoolSize / 2
    }
    
    // 设置合理的超时时间
    if cfg.DialTimeout == 0 {
        cfg.DialTimeout = 5 * time.Second
    }
    
    if cfg.ReadTimeout == 0 {
        cfg.ReadTimeout = 3 * time.Second
    }
    
    if cfg.WriteTimeout == 0 {
        cfg.WriteTimeout = 3 * time.Second
    }
}
```

### 管道操作

```go
// BatchOperations 批量操作示例
func BatchOperations(ctx context.Context, operations map[string]interface{}) error {
    if client == nil {
        return errors.New("Redis client not initialized")
    }
    
    // 使用管道批量执行操作
    pipe := client.Pipeline()
    
    for key, value := range operations {
        pipe.Set(ctx, key, value, time.Hour)
    }
    
    // 执行管道
    _, err := pipe.Exec(ctx)
    if err != nil {
        klogger.Error("Pipeline execution failed", zap.Error(err))
        return err
    }
    
    klogger.Info("Batch operations completed",
        zap.Int("count", len(operations)),
    )
    return nil
}
```

## 监控和指标

### 性能监控

```go
// MonitoringMiddleware Redis 监控中间件
func MonitoringMiddleware() {
    // 定期收集 Redis 指标
    ticker := time.NewTicker(30 * time.Second)
    defer ticker.Stop()
    
    for range ticker.C {
        if client == nil {
            continue
        }
        
        stats := client.PoolStats()
        
        // 记录连接池指标
        klogger.Info("Redis metrics",
            zap.Uint32("pool_hits", stats.Hits),
            zap.Uint32("pool_misses", stats.Misses),
            zap.Uint32("pool_timeouts", stats.Timeouts),
            zap.Uint32("total_connections", stats.TotalConns),
            zap.Uint32("idle_connections", stats.IdleConns),
        )
        
        // 检查连接池健康状态
        if stats.TotalConns > 0 {
            hitRate := float64(stats.Hits) / float64(stats.Hits+stats.Misses) * 100
            if hitRate < 80 {
                klogger.Warn("Low Redis connection pool hit rate",
                    zap.Float64("hit_rate", hitRate),
                )
            }
        }
    }
}
```

## 测试

### 单元测试

```go
func TestRedisOperations(t *testing.T) {
    // 设置测试环境
    cfg := &config.RedisConfig{
        Host:     "localhost",
        Port:     6379,
        DB:       1, // 使用测试数据库
        PoolSize: 5,
    }
    
    err := Init(cfg)
    assert.NoError(t, err)
    defer Close()
    
    ctx := context.Background()
    
    // 测试基础操作
    t.Run("Basic Operations", func(t *testing.T) {
        key := "test:key"
        value := "test value"
        
        // 测试设置
        err := Set(ctx, key, value, time.Minute)
        assert.NoError(t, err)
        
        // 测试获取
        result, err := Get(ctx, key)
        assert.NoError(t, err)
        assert.Equal(t, value, result)
        
        // 测试删除
        err = Del(ctx, key)
        assert.NoError(t, err)
        
        // 验证删除
        result, err = Get(ctx, key)
        assert.NoError(t, err)
        assert.Empty(t, result)
    })
    
    // 测试哈希操作
    t.Run("Hash Operations", func(t *testing.T) {
        key := "test:hash"
        
        err := HSet(ctx, key, "field1", "value1", "field2", "value2")
        assert.NoError(t, err)
        
        value, err := HGet(ctx, key, "field1")
        assert.NoError(t, err)
        assert.Equal(t, "value1", value)
        
        all, err := HGetAll(ctx, key)
        assert.NoError(t, err)
        assert.Equal(t, "value1", all["field1"])
        assert.Equal(t, "value2", all["field2"])
    })
}
```

## 故障排查

### 常见问题

1. **连接超时**
   ```go
   // 检查网络连接和防火墙设置
   // 调整 DialTimeout 配置
   cfg.DialTimeout = 10 * time.Second
   ```

2. **连接池耗尽**
   ```go
   // 增加连接池大小
   cfg.PoolSize = 20
   // 检查连接泄漏
   stats := client.PoolStats()
   if stats.TotalConns >= cfg.PoolSize {
       log.Println("Connection pool exhausted")
   }
   ```

3. **内存使用过高**
   ```go
   // 设置合理的过期时间
   Set(ctx, key, value, time.Hour) // 避免永不过期的键
   
   // 定期清理过期键
   client.FlushDB(ctx) // 谨慎使用
   ```

## 相关文档

- [数据库系统完整指南](DATABASE_GUIDE.md)
- [日志系统文档](LOGGING.md)
- [配置管理文档](../README.md#配置管理)
- [Redis 官方文档](https://redis.io/documentation)
- [go-redis 客户端文档](https://github.com/go-redis/redis)

---

**最佳实践**: 合理配置连接池参数；为所有缓存设置合理的过期时间；使用分布式锁处理并发问题；定期监控 Redis 性能指标；在测试环境使用独立的数据库；实现优雅的错误处理和重试机制。