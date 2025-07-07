package logger

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"gorm.io/gorm/logger"

	"github.com/cuiyuanxin/kunpeng/internal/config"
)

var (
	Logger     *zap.Logger
	GormLogger logger.Interface
	// 钩子注册表
	hookRegistry = make(map[string]HookFunc)
)

// HookFunc 钩子函数类型
type HookFunc func(entry zapcore.Entry, fields []zapcore.Field) error

// Hook 钩子接口
type Hook interface {
	Fire(entry zapcore.Entry, fields []zapcore.Field) error
	Levels() []zapcore.Level
}

// CustomHook 自定义钩子实现
type CustomHook struct {
	name     string
	levels   []zapcore.Level
	hookFunc HookFunc
	logger   *zap.Logger
}

// Fire 执行钩子
func (h *CustomHook) Fire(entry zapcore.Entry, fields []zapcore.Field) error {
	return h.hookFunc(entry, fields)
}

// Levels 返回钩子关注的日志级别
func (h *CustomHook) Levels() []zapcore.Level {
	return h.levels
}

// GormZapLogger GORM的Zap日志适配器
type GormZapLogger struct {
	zapLogger    *zap.Logger
	sqlLogger    *zap.Logger
	errorLogger  *zap.Logger
	logLevel     logger.LogLevel
	slowThreshold time.Duration
}

// NewGormZapLogger 创建GORM Zap日志适配器
func NewGormZapLogger(cfg *config.GormLogging, environment string) *GormZapLogger {
	if !cfg.Enabled {
		return &GormZapLogger{
			zapLogger:     zap.NewNop(),
			sqlLogger:     zap.NewNop(),
			errorLogger:   zap.NewNop(),
			logLevel:      logger.Silent,
			slowThreshold: 200 * time.Millisecond,
		}
	}

	// 解析慢查询阈值
	slowThreshold := 200 * time.Millisecond
	if cfg.SlowThreshold != "" {
		if duration, err := time.ParseDuration(cfg.SlowThreshold); err == nil {
			slowThreshold = duration
		}
	}

	// 确定输出方式
	output := determineOutput(cfg.AutoMode, cfg.ForceConsole, cfg.ForceFile, environment)

	// 创建SQL日志器
	sqlLogger := createSpecialLogger("sql", cfg.SQLFile, output, cfg.Level)
	// 创建错误日志器
	errorLogger := createSpecialLogger("gorm-error", cfg.ErrorFile, output, "error")
	// 创建通用GORM日志器
	gormLogger := createSpecialLogger("gorm", "", output, cfg.Level)

	return &GormZapLogger{
		zapLogger:     gormLogger,
		sqlLogger:     sqlLogger,
		errorLogger:   errorLogger,
		logLevel:      parseGormLogLevel(cfg.Level),
		slowThreshold: slowThreshold,
	}
}

// Init 初始化日志器
func Init(cfg *config.Logging) error {
	// 获取环境信息（从配置中获取，如果没有则默认为development）
	environment := "development" // 这里应该从应用配置中获取
	
	// 环境自适应配置处理
	if cfg.AutoMode {
		applyAutoModeConfig(cfg, environment)
	}
	
	// 初始化主日志器
	var err error
	if cfg.SeparateFiles {
		err = initSeparateFilesLogger(cfg)
	} else {
		err = initSingleFileLogger(cfg)
	}
	
	if err != nil {
		return err
	}
	
	// 初始化GORM日志器
	GormLogger = NewGormZapLogger(&cfg.Gorm, environment)
	
	// 初始化钩子
	initHooks(cfg, environment)
	
	return nil
}

// InitWithEnvironment 带环境参数的初始化方法
func InitWithEnvironment(cfg *config.Logging, environment string) error {
	// 环境自适应配置处理
	if cfg.AutoMode {
		applyAutoModeConfig(cfg, environment)
	}
	
	// 初始化主日志器
	var err error
	if cfg.SeparateFiles {
		err = initSeparateFilesLogger(cfg)
	} else {
		err = initSingleFileLogger(cfg)
	}
	
	if err != nil {
		return err
	}
	
	// 初始化GORM日志器
	GormLogger = NewGormZapLogger(&cfg.Gorm, environment)
	
	// 初始化钩子
	initHooks(cfg, environment)
	
	return nil
}

// initSingleFileLogger 初始化单文件日志器（原有逻辑）
func initSingleFileLogger(cfg *config.Logging) error {
	// 创建日志目录
	if cfg.Output == "file" && cfg.FilePath != "" {
		dir := filepath.Dir(cfg.FilePath)
		if err := os.MkdirAll(dir, 0755); err != nil {
			return err
		}
	}

	// 设置日志级别
	level := getLogLevel(cfg.Level)

	// 创建编码器配置
	encoderConfig := getEncoderConfig(cfg.Format)

	// 创建核心
	core := zapcore.NewCore(
		getEncoder(cfg.Format, encoderConfig),
		getWriteSyncer(cfg),
		level,
	)

	// 创建日志器
	Logger = zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))

	return nil
}

// initSeparateFilesLogger 初始化分级别文件日志器
func initSeparateFilesLogger(cfg *config.Logging) error {
	// 创建日志目录
	logDir := cfg.LogDir
	if logDir == "" {
		logDir = "logs"
	}
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return err
	}

	// 设置日志级别
	level := getLogLevel(cfg.Level)

	// 创建编码器配置
	encoderConfig := getEncoderConfig(cfg.Format)
	encoder := getEncoder(cfg.Format, encoderConfig)

	// 创建不同级别的写入器
	cores := []zapcore.Core{}

	// Debug级别日志
	if level <= zapcore.DebugLevel {
		debugWriter := getLevelWriteSyncer(cfg, logDir, "debug")
		debugCore := zapcore.NewCore(encoder, debugWriter, zapcore.LevelEnabler(zapcore.DebugLevel))
		cores = append(cores, debugCore)
	}

	// Info级别日志
	if level <= zapcore.InfoLevel {
		infoWriter := getMultiLevelWriteSyncer(cfg, logDir, "info", []zapcore.Level{zapcore.InfoLevel})
		infoCore := zapcore.NewCore(encoder, infoWriter, zapcore.LevelEnabler(zapcore.InfoLevel))
		cores = append(cores, infoCore)
	}

	// Warn级别日志
	if level <= zapcore.WarnLevel {
		warnWriter := getMultiLevelWriteSyncer(cfg, logDir, "warn", []zapcore.Level{zapcore.WarnLevel})
		warnCore := zapcore.NewCore(encoder, warnWriter, zapcore.LevelEnabler(zapcore.WarnLevel))
		cores = append(cores, warnCore)
	}

	// Error级别日志（包含Error、Panic、Fatal）
	if level <= zapcore.ErrorLevel {
		errorWriter := getMultiLevelWriteSyncer(cfg, logDir, "error", []zapcore.Level{zapcore.ErrorLevel, zapcore.PanicLevel, zapcore.FatalLevel})
		errorCore := zapcore.NewCore(encoder, errorWriter, zapcore.LevelEnabler(zapcore.ErrorLevel))
		cores = append(cores, errorCore)
	}

	// 如果输出到控制台，添加控制台输出
	if cfg.Output == "stdout" || cfg.Output == "both" {
		consoleCore := zapcore.NewCore(
			getEncoder("console", getEncoderConfig("console")),
			zapcore.AddSync(os.Stdout),
			level,
		)
		cores = append(cores, consoleCore)
	}

	// 合并所有核心
	core := zapcore.NewTee(cores...)

	// 创建日志器
	Logger = zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))

	return nil
}

// getLogLevel 获取日志级别
func getLogLevel(level string) zapcore.Level {
	switch level {
	case "debug":
		return zapcore.DebugLevel
	case "info":
		return zapcore.InfoLevel
	case "warn":
		return zapcore.WarnLevel
	case "error":
		return zapcore.ErrorLevel
	case "panic":
		return zapcore.PanicLevel
	case "fatal":
		return zapcore.FatalLevel
	default:
		return zapcore.InfoLevel
	}
}

// getEncoderConfig 获取编码器配置
func getEncoderConfig(format string) zapcore.EncoderConfig {
	config := zap.NewProductionEncoderConfig()
	config.TimeKey = "timestamp"
	config.EncodeTime = zapcore.ISO8601TimeEncoder
	config.EncodeLevel = zapcore.CapitalLevelEncoder
	config.EncodeCaller = zapcore.ShortCallerEncoder

	if format == "console" {
		config.EncodeLevel = zapcore.CapitalColorLevelEncoder
	}

	return config
}

// getEncoder 获取编码器
func getEncoder(format string, config zapcore.EncoderConfig) zapcore.Encoder {
	switch format {
	case "json":
		return zapcore.NewJSONEncoder(config)
	case "console":
		return zapcore.NewConsoleEncoder(config)
	default:
		return zapcore.NewJSONEncoder(config)
	}
}

// getWriteSyncer 获取写入同步器
func getWriteSyncer(cfg *config.Logging) zapcore.WriteSyncer {
	if cfg.Output == "stdout" {
		return zapcore.AddSync(os.Stdout)
	}

	// 文件输出，使用 lumberjack 进行日志轮转
	lumberjackLogger := &lumberjack.Logger{
		Filename:   cfg.FilePath,
		MaxSize:    cfg.MaxSize,
		MaxBackups: cfg.MaxBackups,
		MaxAge:     cfg.MaxAge,
		Compress:   cfg.Compress,
	}

	return zapcore.AddSync(lumberjackLogger)
}

// getLevelWriteSyncer 获取指定级别的写入同步器
func getLevelWriteSyncer(cfg *config.Logging, logDir, level string) zapcore.WriteSyncer {
	filename := filepath.Join(logDir, fmt.Sprintf("%s.log", level))
	lumberjackLogger := &lumberjack.Logger{
		Filename:   filename,
		MaxSize:    cfg.MaxSize,
		MaxBackups: cfg.MaxBackups,
		MaxAge:     cfg.MaxAge,
		Compress:   cfg.Compress,
	}
	return zapcore.AddSync(lumberjackLogger)
}

// getMultiLevelWriteSyncer 获取多级别共享的写入同步器
func getMultiLevelWriteSyncer(cfg *config.Logging, logDir, filename string, levels []zapcore.Level) zapcore.WriteSyncer {
	filePath := filepath.Join(logDir, fmt.Sprintf("%s.log", filename))
	lumberjackLogger := &lumberjack.Logger{
		Filename:   filePath,
		MaxSize:    cfg.MaxSize,
		MaxBackups: cfg.MaxBackups,
		MaxAge:     cfg.MaxAge,
		Compress:   cfg.Compress,
	}
	return zapcore.AddSync(lumberjackLogger)
}

// Debug 调试日志
func Debug(msg string, fields ...zap.Field) {
	Logger.Debug(msg, fields...)
}

// Info 信息日志
func Info(msg string, fields ...zap.Field) {
	Logger.Info(msg, fields...)
}

// Warn 警告日志
func Warn(msg string, fields ...zap.Field) {
	Logger.Warn(msg, fields...)
}

// Error 错误日志
func Error(msg string, fields ...zap.Field) {
	Logger.Error(msg, fields...)
}

// Panic panic日志
func Panic(msg string, fields ...zap.Field) {
	Logger.Panic(msg, fields...)
}

// Fatal 致命错误日志
func Fatal(msg string, fields ...zap.Field) {
	Logger.Fatal(msg, fields...)
}

// Sync 同步日志
func Sync() error {
	return Logger.Sync()
}

// ============= 环境自适应配置 =============

// applyAutoModeConfig 应用环境自适应配置
func applyAutoModeConfig(cfg *config.Logging, environment string) {
	// 强制配置优先级最高
	if cfg.ForceConsole != nil {
		if *cfg.ForceConsole {
			cfg.Output = "stdout"
		} else {
			cfg.Output = "file"
		}
		return
	}
	
	if cfg.ForceFile != nil {
		if *cfg.ForceFile {
			cfg.Output = "file"
		} else {
			cfg.Output = "stdout"
		}
		return
	}
	
	// 根据环境自动配置
	switch environment {
	case "production", "prod":
		cfg.Output = "file"
		if cfg.Format == "" {
			cfg.Format = "json"
		}
	case "development", "dev", "test":
		cfg.Output = "stdout"
		if cfg.Format == "" {
			cfg.Format = "console"
		}
	default:
		cfg.Output = "stdout"
	}
}

// determineOutput 确定输出方式
func determineOutput(autoMode bool, forceConsole, forceFile *bool, environment string) string {
	if forceConsole != nil && *forceConsole {
		return "stdout"
	}
	if forceFile != nil && *forceFile {
		return "file"
	}
	
	if autoMode {
		switch environment {
		case "production", "prod":
			return "file"
		default:
			return "stdout"
		}
	}
	
	return "stdout"
}

// ============= GORM日志适配器方法 =============

// LogMode 设置日志模式
func (l *GormZapLogger) LogMode(level logger.LogLevel) logger.Interface {
	newLogger := *l
	newLogger.logLevel = level
	return &newLogger
}

// Info 信息日志
func (l *GormZapLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	if l.logLevel >= logger.Info {
		l.zapLogger.Sugar().Infof(msg, data...)
	}
}

// Warn 警告日志
func (l *GormZapLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	if l.logLevel >= logger.Warn {
		l.zapLogger.Sugar().Warnf(msg, data...)
	}
}

// Error 错误日志
func (l *GormZapLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	if l.logLevel >= logger.Error {
		l.errorLogger.Sugar().Errorf(msg, data...)
	}
}

// Trace SQL追踪日志
func (l *GormZapLogger) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	if l.logLevel <= logger.Silent {
		return
	}
	
	elapsed := time.Since(begin)
	sql, rows := fc()
	
	fields := []zap.Field{
		zap.String("sql", sql),
		zap.Duration("elapsed", elapsed),
		zap.Int64("rows", rows),
	}
	
	switch {
	case err != nil && l.logLevel >= logger.Error:
		l.errorLogger.Error("SQL Error", append(fields, zap.Error(err))...)
	case elapsed > l.slowThreshold && l.slowThreshold != 0 && l.logLevel >= logger.Warn:
		l.sqlLogger.Warn("Slow SQL", append(fields, zap.Duration("threshold", l.slowThreshold))...)
	case l.logLevel == logger.Info:
		l.sqlLogger.Info("SQL", fields...)
	}
}

// parseGormLogLevel 解析GORM日志级别
func parseGormLogLevel(level string) logger.LogLevel {
	switch level {
	case "silent":
		return logger.Silent
	case "error":
		return logger.Error
	case "warn":
		return logger.Warn
	case "info":
		return logger.Info
	default:
		return logger.Info
	}
}

// createSpecialLogger 创建特殊用途的日志器
func createSpecialLogger(name, filePath, output, level string) *zap.Logger {
	logLevel := getLogLevel(level)
	encoderConfig := getEncoderConfig("json")
	
	var writeSyncer zapcore.WriteSyncer
	if output == "file" && filePath != "" {
		// 确保目录存在
		dir := filepath.Dir(filePath)
		os.MkdirAll(dir, 0755)
		
		lumberjackLogger := &lumberjack.Logger{
			Filename:   filePath,
			MaxSize:    100,
			MaxBackups: 3,
			MaxAge:     28,
			Compress:   true,
		}
		writeSyncer = zapcore.AddSync(lumberjackLogger)
	} else {
		writeSyncer = zapcore.AddSync(os.Stdout)
		encoderConfig = getEncoderConfig("console")
	}
	
	core := zapcore.NewCore(
		getEncoder("json", encoderConfig),
		writeSyncer,
		logLevel,
	)
	
	return zap.New(core, zap.AddCaller())
}

// ============= 钩子系统 =============

// RegisterHook 注册钩子
func RegisterHook(name string, hookFunc HookFunc) {
	hookRegistry[name] = hookFunc
}

// UnregisterHook 注销钩子
func UnregisterHook(name string) {
	delete(hookRegistry, name)
}

// initHooks 初始化钩子
func initHooks(cfg *config.Logging, environment string) {
	// 初始化链路追踪钩子
	if cfg.Hooks.Tracing.Enabled {
		initTracingHook(&cfg.Hooks.Tracing, environment)
	}
	
	// 初始化gRPC钩子
	if cfg.Hooks.GRPC.Enabled {
		initGRPCHook(&cfg.Hooks.GRPC, environment)
	}
	
	// 初始化自定义钩子
	if cfg.Hooks.Custom.Enabled {
		initCustomHook(&cfg.Hooks.Custom, environment)
	}
}

// initTracingHook 初始化链路追踪钩子
func initTracingHook(cfg *config.LoggingHook, environment string) {
	hookFunc := func(entry zapcore.Entry, fields []zapcore.Field) error {
		// 这里可以集成OpenTelemetry或Jaeger等链路追踪系统
		// 示例：添加trace_id和span_id到日志中
		for i, field := range fields {
			if field.Key == "trace_id" || field.Key == "span_id" {
				// 处理链路追踪相关字段
				_ = i // 避免未使用变量警告
			}
		}
		return nil
	}
	
	RegisterHook("tracing", hookFunc)
}

// initGRPCHook 初始化gRPC钩子
func initGRPCHook(cfg *config.LoggingHook, environment string) {
	hookFunc := func(entry zapcore.Entry, fields []zapcore.Field) error {
		// 这里可以集成gRPC日志记录
		// 示例：记录gRPC请求和响应信息
		for i, field := range fields {
			if field.Key == "grpc_method" || field.Key == "grpc_code" {
				// 处理gRPC相关字段
				_ = i // 避免未使用变量警告
			}
		}
		return nil
	}
	
	RegisterHook("grpc", hookFunc)
}

// initCustomHook 初始化自定义钩子
func initCustomHook(cfg *config.LoggingHook, environment string) {
	hookFunc := func(entry zapcore.Entry, fields []zapcore.Field) error {
		// 自定义钩子逻辑
		// 用户可以在这里添加自己的日志处理逻辑
		return nil
	}
	
	RegisterHook("custom", hookFunc)
}

// ============= 便捷方法 =============

// GetGormLogger 获取GORM日志器
func GetGormLogger() logger.Interface {
	return GormLogger
}

// WithTraceID 添加链路追踪ID
func WithTraceID(traceID string) zap.Field {
	return zap.String("trace_id", traceID)
}

// WithSpanID 添加Span ID
func WithSpanID(spanID string) zap.Field {
	return zap.String("span_id", spanID)
}

// WithGRPCMethod 添加gRPC方法
func WithGRPCMethod(method string) zap.Field {
	return zap.String("grpc_method", method)
}

// WithGRPCCode 添加gRPC状态码
func WithGRPCCode(code int) zap.Field {
	return zap.Int("grpc_code", code)
}

// WithUserID 添加用户ID
func WithUserID(userID string) zap.Field {
	return zap.String("user_id", userID)
}

// WithRequestID 添加请求ID
func WithRequestID(requestID string) zap.Field {
	return zap.String("request_id", requestID)
}
