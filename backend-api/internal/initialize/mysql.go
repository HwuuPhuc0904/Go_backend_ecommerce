package initialize

import (
    "fmt"
    "time"
    
    "GOLANG/github.com/HwuuPhuc0904/backend-api/global"
    "gorm.io/driver/mysql"
    "gorm.io/gorm"
    "gorm.io/gorm/logger"
)

func InitMysql() {
    m := global.Config.MySQL
    
    // Sử dụng DSN (Data Source Name) từ cấu hình
    dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
        m.Username, 
        m.Password, // Mật khẩu mới từ cấu hình
        m.Host,
        m.Port,
        m.Database)
    
    fmt.Println("Connecting to MySQL database...")
    
    // Tạo kết nối database
    db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
        Logger: logger.Default.LogMode(logger.Info),
    })
    
    if err != nil {
        if global.Logger != nil {
            global.Logger.Error(fmt.Sprintf("Failed to connect to database: %v", err))
        }
        panic(fmt.Sprintf("Failed to connect to database: %v", err))
    }
    
    // Cấu hình connection pool
    sqlDB, _ := db.DB()
    sqlDB.SetMaxIdleConns(10)
    sqlDB.SetMaxOpenConns(100)
    sqlDB.SetConnMaxLifetime(time.Hour)
    
    // Lưu instance vào biến global
    global.DB = db
    
    if global.Logger != nil {
        global.Logger.Info("MySQL connected successfully")
    }
    fmt.Println("MySQL connected successfully")
}