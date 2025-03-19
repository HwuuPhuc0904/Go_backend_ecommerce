package initialize

import (
    "GOLANG/github.com/HwuuPhuc0904/backend-api/global"
    "fmt"
    "go.uber.org/zap"
    "gorm.io/driver/mysql"
    "gorm.io/gorm"
    "gorm.io/gorm/logger"
    "log"
    "os"
    "time"
)

// InitMysql khởi tạo kết nối MySQL
func InitMysql() {
    m := global.Config.MySQL
    
    fmt.Printf("MySQL Configuration - Host: %s, Port: %s, User: %s, Database: %s\n", 
        m.Host, m.Port, m.Username, m.Database)
    
    fmt.Println("Connecting to MySQL database...")
    
    // Cấu hình GORM logger
    newLogger := logger.New(
        log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
        logger.Config{
            SlowThreshold: time.Second,   // Thời gian query chậm
            LogLevel:      logger.Silent, // Log level
            Colorful:      false,        // Disable color
        },
    )
    
    // Tạo DSN (Data Source Name)
    dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
        m.Username, 
        m.Password,
        m.Host,
        m.Port,
        m.Database)
    
    // Kết nối đến MySQL
    db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
        Logger: newLogger,
    })
    
    if err != nil {
        global.Logger.Error("Failed to connect to database", zap.Error(err))
        panic(fmt.Sprintf("Failed to connect to database: %v", err))
    }
    
    // Thiết lập các tham số kết nối
    sqlDB, err := db.DB()
    if err != nil {
        global.Logger.Error("Failed to get database connection", zap.Error(err))
        panic(fmt.Sprintf("Failed to get database connection: %v", err))
    }
    
    // Cài đặt connection pool
    sqlDB.SetMaxIdleConns(10)
    sqlDB.SetMaxOpenConns(100)
    sqlDB.SetConnMaxLifetime(time.Hour)
    
    // Kiểm tra kết nối
    err = sqlDB.Ping()
    if err != nil {
        global.Logger.Error("Failed to ping database", zap.Error(err))
        panic(fmt.Sprintf("Failed to ping database: %v", err))
    }
    
    global.DB = db
    global.Logger.Info("Connected to MySQL database successfully")
}