package main

import (
	"os"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func main() {
	encoder := getEncoderLog()
	writerSync := getWriterSync()
	core := zapcore.NewCore(encoder, writerSync, zapcore.InfoLevel)

	logger := zap.New(core)


	logger.Info("Infor Log", zap.Int("line",1))
	logger.Info("Infor Log", zap.Int("line",2))
}

// format log
func getEncoderLog() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()

	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	encoderConfig.TimeKey = "timestamp"

	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder

	encoderConfig.EncodeCaller = zapcore.ShortCallerEncoder


	return zapcore.NewJSONEncoder(encoderConfig)

}

func getWriterSync() zapcore.WriteSyncer {
	file, _ := os.OpenFile("./log/log.txt", os.O_CREATE|os.O_WRONLY|os.O_APPEND, os.ModePerm)
	syncFile := zapcore.AddSync(file)
	syncConsole := zapcore.AddSync(os.Stderr)
	
	return zapcore.NewMultiWriteSyncer(syncFile, syncConsole)
}

