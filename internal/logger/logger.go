package logger

import (
	"fmt"
	"os"
	"path/filepath"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"

	"github.com/cuiyuanxin/kunpeng/internal/config"
)

var Logger *zap.Logger

// Init 初始化日志器
func Init(cfg *config.Logging) error {
	if cfg.SeparateFiles {
		return initSeparateFilesLogger(cfg)
	}
	return initSingleFileLogger(cfg)
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
