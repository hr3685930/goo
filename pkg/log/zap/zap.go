package zap

import (
    log2 "goo/pkg/log"
    "go.uber.org/zap"
    "go.uber.org/zap/zapcore"
    "time"
)

type zaps struct {
    OutputPaths []string
}

func NewZaps(OutputPaths []string) *zaps {
    return &zaps{OutputPaths}
}

func (z *zaps) InitLog() {
    encoderConfig := zapcore.EncoderConfig{
        TimeKey:       "time",
        LevelKey:      "level",
        NameKey:       "logger",
        MessageKey:    "msg",
        StacktraceKey: "stacktrace",
        LineEnding:    zapcore.DefaultLineEnding,
        EncodeLevel:   zapcore.LowercaseLevelEncoder,
        EncodeTime: func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
            enc.AppendString(t.Format("2006-01-02 15:04:05"))
        },
        EncodeDuration: zapcore.SecondsDurationEncoder,
        EncodeCaller:   zapcore.FullCallerEncoder,
    }

    atom := zap.NewAtomicLevelAt(zap.DebugLevel)

    config := zap.Config{
        Level:            atom,
        Development:      true,
        Encoding:         "json",
        EncoderConfig:    encoderConfig,
        OutputPaths:      z.OutputPaths,
        ErrorOutputPaths: []string{"stderr"},
    }

    log2.Logger, _ = config.Build()
}
