package logger

import (
	"os"
	"path/filepath"
	"time"

	"github.com/cuiyuanxin/kunpeng/pkg/config"
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var logger *zap.Logger
var sqlLogger *zap.Logger

// Init 初始化日志
func Init() error {
	logConfig := config.GetLogConfig()

	// 初始化应用日志
	if err := initAppLogger(logConfig); err != nil {
		return err
	}

	// 初始化SQL日志
	if err := initSqlLogger(logConfig); err != nil {
		return err
	}

	return nil
}

// initAppLogger 初始化应用日志
func initAppLogger(logConfig config.LogConfig) error {
	// 设置日志级别
	level := getLogLevel(logConfig.Level)

	// 配置编码器
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		FunctionKey:    zapcore.OmitKey,
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeTime:     timeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	// 配置输出
	var cores []zapcore.Core

	// 根据环境选择输出方式
	if config.IsDevelopment() {
		// 开发环境：只输出到控制台
		consoleEncoder := zapcore.NewConsoleEncoder(encoderConfig)
		consoleCore := zapcore.NewCore(
			consoleEncoder,
			zapcore.AddSync(os.Stdout),
			level,
		)
		cores = append(cores, consoleCore)
	} else {
		// 生产环境：只输出到文件
		// 确保日志目录存在
		logDir := filepath.Dir(logConfig.Filename)
		if err := os.MkdirAll(logDir, 0755); err != nil {
			return err
		}

		fileEncoder := zapcore.NewJSONEncoder(encoderConfig)
		fileWriter := zapcore.AddSync(&lumberjack.Logger{
			Filename:   logConfig.Filename,
			MaxSize:    logConfig.MaxSize,
			MaxBackups: logConfig.MaxBackups,
			MaxAge:     logConfig.MaxAge,
			Compress:   logConfig.Compress,
		})
		fileCore := zapcore.NewCore(
			fileEncoder,
			fileWriter,
			level,
		)
		cores = append(cores, fileCore)
	}

	// 创建Logger
	core := zapcore.NewTee(cores...)
	logger = zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))

	return nil
}

// initSqlLogger 初始化SQL日志
func initSqlLogger(logConfig config.LogConfig) error {
	// 确保SQL日志目录存在
	sqlLogDir := filepath.Dir(logConfig.SqlFilename)
	if err := os.MkdirAll(sqlLogDir, 0755); err != nil {
		return err
	}

	// SQL日志编码器配置
	sqlEncoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		MessageKey:     "msg",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeTime:     timeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
	}

	// SQL日志文件输出
	sqlFileEncoder := zapcore.NewJSONEncoder(sqlEncoderConfig)
	sqlFileWriter := zapcore.AddSync(&lumberjack.Logger{
		Filename:   logConfig.SqlFilename,
		MaxSize:    logConfig.SqlMaxSize,
		MaxBackups: logConfig.SqlMaxBackups,
		MaxAge:     logConfig.SqlMaxAge,
		Compress:   logConfig.SqlCompress,
	})

	// 创建SQL日志核心
	sqlCore := zapcore.NewCore(
		sqlFileEncoder,
		sqlFileWriter,
		zapcore.DebugLevel, // SQL日志记录所有级别
	)

	// 如果是开发环境，同时输出到控制台
	var cores []zapcore.Core
	cores = append(cores, sqlCore)

	if config.IsDevelopment() {
		consoleEncoder := zapcore.NewConsoleEncoder(sqlEncoderConfig)
		consoleCore := zapcore.NewCore(
			consoleEncoder,
			zapcore.AddSync(os.Stdout),
			zapcore.DebugLevel,
		)
		cores = append(cores, consoleCore)
	}

	// 创建SQL Logger
	sqlLoggerCore := zapcore.NewTee(cores...)
	sqlLogger = zap.New(sqlLoggerCore)

	return nil
}

// 自定义时间格式
func timeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("2006-01-02 15:04:05.000"))
}

// 获取日志级别
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
	case "dpanic":
		return zapcore.DPanicLevel
	case "panic":
		return zapcore.PanicLevel
	case "fatal":
		return zapcore.FatalLevel
	default:
		return zapcore.InfoLevel
	}
}

// Debug 调试日志
func Debug(msg string, fields ...zap.Field) {
	logger.Debug(msg, fields...)
}

// Info 信息日志
func Info(msg string, fields ...zap.Field) {
	logger.Info(msg, fields...)
}

// Warn 警告日志
func Warn(msg string, fields ...zap.Field) {
	logger.Warn(msg, fields...)
}

// Error 错误日志
func Error(msg string, err error, fields ...zap.Field) {
	if err != nil {
		fields = append(fields, zap.Error(err))
	}
	logger.Error(msg, fields...)
}

// DPanic 开发环境宕机日志
func DPanic(msg string, fields ...zap.Field) {
	logger.DPanic(msg, fields...)
}

// Panic 宕机日志
func Panic(msg string, fields ...zap.Field) {
	logger.Panic(msg, fields...)
}

// Fatal 致命错误日志
func Fatal(msg string, fields ...zap.Field) {
	logger.Fatal(msg, fields...)
}

// With 创建子日志
func With(fields ...zap.Field) *zap.Logger {
	return logger.With(fields...)
}

// Sync 同步日志
func Sync() error {
	if err := logger.Sync(); err != nil {
		return err
	}
	if sqlLogger != nil {
		return sqlLogger.Sync()
	}
	return nil
}

// GetLogger 获取日志实例
func GetLogger() *zap.Logger {
	return logger
}

// Field 创建字段
func Field(key string, value interface{}) zap.Field {
	return zap.Any(key, value)
}

// String 创建字符串字段
func String(key string, value string) zap.Field {
	return zap.String(key, value)
}

// Int 创建整数字段
func Int(key string, value int) zap.Field {
	return zap.Int(key, value)
}

// Bool 创建布尔字段
func Bool(key string, value bool) zap.Field {
	return zap.Bool(key, value)
}

// SQL日志相关方法

// SqlDebug SQL调试日志
func SqlDebug(msg string, fields ...zap.Field) {
	sqlLogger.Debug(msg, fields...)
}

// SqlInfo SQL信息日志
func SqlInfo(msg string, fields ...zap.Field) {
	sqlLogger.Info(msg, fields...)
}

// SqlWarn SQL警告日志
func SqlWarn(msg string, fields ...zap.Field) {
	sqlLogger.Warn(msg, fields...)
}

// SqlError SQL错误日志
func SqlError(msg string, err error, fields ...zap.Field) {
	if err != nil {
		fields = append(fields, zap.Error(err))
	}
	sqlLogger.Error(msg, fields...)
}

// GetSqlLogger 获取SQL日志实例
func GetSqlLogger() *zap.Logger {
	return sqlLogger
}
