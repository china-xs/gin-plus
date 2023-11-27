// Package utils
// @file      : p.db.go
// @author    : china.gdxs@gmail.com
// @time      : 2023/11/9 16:18
// @Description:
package utils

import (
	"context"
	"fmt"
	"github.com/spf13/viper"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
	"time"
)

const format = "2006-01-02 15:04:05"

// LogOpts is log configuration struct
type LogOpts struct {
	Filename    string `yaml:"filename"`    // 文件名称
	MaxSize     int    `yaml:"maxSize"`     // 最大文件
	MaxBackups  int    `yaml:"maxBackups"`  // 最大备份数
	MaxAge      int    `yaml:"maxAge"`      //保留时长天 days
	Level       string `yaml:"level"`       // 日志登记 对应zap.level
	Stdout      bool   `yaml:"stdout"`      // 是否在终端输出
	Compress    bool   `yaml:"compress"`    // 是否压缩文件
	LocalTime   bool   `yaml:"localTime"`   // 使用本地时间
	IsContainer bool   `yaml:"isContainer"` // 是否是容器, 容器存在公用app.log 情况可以按容器ip 区分
}

// NewLog for init zap log library
func NewLog(v *viper.Viper) (l *zap.Logger, fc func(), err error) {
	var (
		level  = zap.NewAtomicLevel()
		logger *zap.Logger
		o      = new(LogOpts)
	)
	if err = v.UnmarshalKey("log", o); err != nil {
		return
	}
	// 按照ip 处理 待实现
	if o.IsContainer {
		fmt.Println(`debug`)
	}
	if level, err = zap.ParseAtomicLevel(o.Level); err != nil {
		return
	}
	write := &lumberjack.Logger{ // concurrent-safed
		Filename:   o.Filename,   // 文件路径
		MaxSize:    o.MaxSize,    // MaxSize 兆字节
		MaxBackups: o.MaxBackups, // 最多保留 300 个备份
		MaxAge:     o.MaxAge,     // 最大时间，默认单位 day
		LocalTime:  true,         // 使用本地时间
		Compress:   o.Compress,   // 是否压缩 disabled by default
	}

	fw := zapcore.AddSync(write)
	cw := zapcore.Lock(os.Stdout)

	// file core 采用jsonEncoder
	cores := make([]zapcore.Core, 0, 2)
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:       "time",
		LevelKey:      "level",
		NameKey:       "logger", // used by logger.Named(key); optional; useless
		MessageKey:    "msg",
		StacktraceKey: "stacktrace", // use by zap.AddStacktrace; optional; useless
		LineEnding:    zapcore.DefaultLineEnding,
		EncodeLevel:   zapcore.LowercaseLevelEncoder, // 小写编码器
		EncodeTime: func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendString(t.Format(format))
		},
		EncodeDuration: zapcore.MillisDurationEncoder,
	}
	je := zapcore.NewJSONEncoder(encoderConfig)
	cores = append(cores, zapcore.NewCore(je, fw, level))

	if o.Stdout {
		ce := zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig())
		cores = append(cores, zapcore.NewCore(ce, cw, level))
	}
	core := zapcore.NewTee(cores...)
	logger = zap.New(core)
	zap.ReplaceGlobals(logger)
	fc = func() {
		logger.Sync() // 缓存
		write.Close() // os close
	}
	return logger, fc, err
}

func WithCtx(ctx context.Context, log *zap.Logger) *zap.Logger {
	var traceId, spanId string
	if span := trace.SpanContextFromContext(ctx); span.HasTraceID() {
		traceId = span.TraceID().String()
	}
	if span := trace.SpanContextFromContext(ctx); span.HasSpanID() {
		spanId = span.SpanID().String()
	}
	var fields []zap.Field
	fields = append(fields,
		zap.String(_tranceId, traceId),
		zap.String(_spanId, spanId),
	)
	return log.With(fields...)
}

type Log struct {
	lg *zap.Logger
}

func NewWLog(lg *zap.Logger) *Log {
	return &Log{lg: lg}
}

func (this Log) WithCtx(ctx context.Context) *zap.Logger {
	return WithCtx(ctx, this.lg)
}
