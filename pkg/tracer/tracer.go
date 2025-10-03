package tracer

import (
	"context"
	"fmt"
	"sync"

	"github.com/cuiyuanxin/kunpeng/pkg/config"
	"github.com/cuiyuanxin/kunpeng/pkg/logger"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

// TraceKey 上下文中的追踪键
type TraceKey string

const (
	// TraceIDKey 追踪ID键
	TraceIDKey TraceKey = "trace_id"
	// SpanIDKey 跨度ID键
	SpanIDKey TraceKey = "span_id"
	// ParentSpanIDKey 父跨度ID键
	ParentSpanIDKey TraceKey = "parent_span_id"
)

var (
	tracer *Tracer
	once   sync.Once
)

// Tracer 追踪器
type Tracer struct {
	enabled bool
}

// Init 初始化追踪器
func Init() {
	once.Do(func() {
		tracer = &Tracer{
			enabled: config.GetAppConfig().TraceEnable,
		}
		logger.Info("链路追踪初始化成功", zap.Bool("enabled", tracer.enabled))
	})
}

// GetTracer 获取追踪器
func GetTracer() *Tracer {
	if tracer == nil {
		Init()
	} else {
		// 更新追踪器状态以反映当前配置
		tracer.enabled = config.GetAppConfig().TraceEnable
	}
	return tracer
}

// IsEnabled 是否启用追踪
func (t *Tracer) IsEnabled() bool {
	return t.enabled
}

// NewContext 创建带追踪信息的上下文
func (t *Tracer) NewContext(ctx context.Context) context.Context {
	if !t.enabled {
		return ctx
	}

	// 生成追踪ID
	traceID := uuid.New().String()
	spanID := uuid.New().String()

	// 设置追踪信息
	ctx = context.WithValue(ctx, TraceIDKey, traceID)
	ctx = context.WithValue(ctx, SpanIDKey, spanID)

	return ctx
}

// NewChildContext 创建子追踪上下文
func (t *Tracer) NewChildContext(ctx context.Context) context.Context {
	if !t.enabled {
		return ctx
	}

	// 获取父追踪信息
	traceID := GetTraceID(ctx)
	parentSpanID := GetSpanID(ctx)

	// 如果没有父追踪信息，则创建新的追踪上下文
	if traceID == "" {
		return t.NewContext(ctx)
	}

	// 生成新的跨度ID
	spanID := uuid.New().String()

	// 设置追踪信息
	ctx = context.WithValue(ctx, TraceIDKey, traceID)
	ctx = context.WithValue(ctx, SpanIDKey, spanID)
	ctx = context.WithValue(ctx, ParentSpanIDKey, parentSpanID)

	return ctx
}

// GetTraceID 获取追踪ID
func GetTraceID(ctx context.Context) string {
	if ctx == nil {
		return ""
	}
	if traceID, ok := ctx.Value(TraceIDKey).(string); ok {
		return traceID
	}
	return ""
}

// GetSpanID 获取跨度ID
func GetSpanID(ctx context.Context) string {
	if ctx == nil {
		return ""
	}
	if spanID, ok := ctx.Value(SpanIDKey).(string); ok {
		return spanID
	}
	return ""
}

// GetParentSpanID 获取父跨度ID
func GetParentSpanID(ctx context.Context) string {
	if ctx == nil {
		return ""
	}
	if parentSpanID, ok := ctx.Value(ParentSpanIDKey).(string); ok {
		return parentSpanID
	}
	return ""
}

// LogFields 获取追踪日志字段
func LogFields(ctx context.Context) []zap.Field {
	fields := []zap.Field{}

	if traceID := GetTraceID(ctx); traceID != "" {
		fields = append(fields, zap.String("trace_id", traceID))
	}

	if spanID := GetSpanID(ctx); spanID != "" {
		fields = append(fields, zap.String("span_id", spanID))
	}

	if parentSpanID := GetParentSpanID(ctx); parentSpanID != "" {
		fields = append(fields, zap.String("parent_span_id", parentSpanID))
	}

	return fields
}

// FormatTraceInfo 格式化追踪信息
func FormatTraceInfo(ctx context.Context) string {
	traceID := GetTraceID(ctx)
	spanID := GetSpanID(ctx)
	parentSpanID := GetParentSpanID(ctx)

	if traceID == "" {
		return ""
	}

	if parentSpanID == "" {
		return fmt.Sprintf("trace_id=%s span_id=%s", traceID, spanID)
	}

	return fmt.Sprintf("trace_id=%s span_id=%s parent_span_id=%s", traceID, spanID, parentSpanID)
}
