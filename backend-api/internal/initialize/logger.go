package initialize

import (
    "fmt"
    "GOLANG/github.com/HwuuPhuc0904/backend-api/global"
    logger "GOLANG/github.com/HwuuPhuc0904/backend-api/pkg/logger"
    "go.uber.org/zap"
)

func InitLogger() {
	global.Logger = logger.NewLogger(global.Config.LoggerSetting)
    
    // Ghi log thử nghiệm với các mức độ khác nhau
    global.Logger.Info("Logger initialized successfully",
    zap.String("filepath", global.Config.LoggerSetting.FilePath))
    global.Logger.Debug("This is a debug message")
    global.Logger.Warn("This is a warning message")
    global.Logger.Error("This is an error message for testing")
    
    fmt.Println("Logger path:", global.Config.LoggerSetting.FilePath)
    fmt.Println("Logger has been initialized. Check log file.")
}