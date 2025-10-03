package database

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/cuiyuanxin/kunpeng/pkg/config"
	"github.com/cuiyuanxin/kunpeng/pkg/logger"
	"go.uber.org/zap"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

// GormLogger 自定义GORM日志记录器
type GormLogger struct {
	SlowThreshold         time.Duration
	SourceField           string
	SkipErrRecordNotFound bool
}

// NewGormLogger 创建GORM日志记录器
func NewGormLogger() gormlogger.Interface {
	return &GormLogger{
		SlowThreshold:         time.Second, // 慢查询阈值
		SkipErrRecordNotFound: true,        // 忽略记录未找到错误
	}
}

// LogMode 设置日志级别
func (l *GormLogger) LogMode(level gormlogger.LogLevel) gormlogger.Interface {
	return l
}

// Info 记录信息日志
func (l *GormLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	logger.SqlInfo(fmt.Sprintf(msg, data...))
}

// Warn 记录警告日志
func (l *GormLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	logger.SqlWarn(fmt.Sprintf(msg, data...))
}

// Error 记录错误日志
func (l *GormLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	logger.SqlError(fmt.Sprintf(msg, data...), nil)
}

// Trace 记录SQL执行日志
func (l *GormLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	// 如果不显示SQL，则直接返回
	if !config.GetDatabaseConfig().ShowSql {
		return
	}

	// 计算执行时间
	elapsed := time.Since(begin)
	// 获取SQL和影响行数
	sql, rows := fc()

	// 构建日志字段
	fields := []zap.Field{
		zap.String("sql", sql),
		zap.Int64("rows", rows),
		zap.Duration("elapsed", elapsed),
	}

	// 判断是否为慢查询
	if elapsed > l.SlowThreshold {
		logger.SqlWarn("GORM 慢查询", fields...)
		return
	}

	// 判断是否有错误
	if err != nil && (!errors.Is(err, gorm.ErrRecordNotFound) || !l.SkipErrRecordNotFound) {
		fields = append(fields, zap.Error(err))
		logger.SqlError("GORM 查询错误", err, fields...)
		return
	}

	// 记录正常查询
	logger.SqlDebug("GORM 查询", fields...)
}
