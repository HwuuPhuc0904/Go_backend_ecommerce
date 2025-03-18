package logger

import (
	"os"
	setting "GOLANG/github.com/HwuuPhuc0904/backend-api/settings"

	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type LoggerZap struct {
	*zap.Logger
}

func NewLogger(loggerConfig setting.LoggerSetting) *LoggerZap {
	
	loglevel := loggerConfig.Level
	var level zapcore.Level

	switch loglevel {
	case "debug":
		level = zapcore.DebugLevel
	case "info":
		level = zapcore.InfoLevel
	case "warn":
		level = zapcore.WarnLevel
	case "error":
		level = zapcore.ErrorLevel
	default:
		level = zapcore.InfoLevel
	}
	encoder := GetEncoderLog()
	hook := lumberjack.Logger{
		Filename:   loggerConfig.FilePath,
		MaxSize:    loggerConfig.MaxSize,
		MaxBackups: loggerConfig.MaxBackup,
		MaxAge:     loggerConfig.MaxAge,
		Compress:   loggerConfig.Compress,
	}
	core := zapcore.NewCore(
		encoder,
		zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout),zapcore.AddSync(&hook)),
		level,
	)
	return &LoggerZap{zap.New(core, zap.AddCaller(), zap.AddStacktrace(zap.ErrorLevel))}
}

func GetEncoderLog() zapcore.Encoder {
	zapcodeConfig := zap.NewProductionEncoderConfig()
	zapcodeConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	zapcodeConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	zapcodeConfig.EncodeCaller = zapcore.ShortCallerEncoder
	return zapcore.NewJSONEncoder(zapcodeConfig)
} 