package initialize

import(
    "fmt"
    "path/filepath"
    "github.com/spf13/viper"
    "GOLANG/github.com/HwuuPhuc0904/backend-api/global"
)

func LoadConfig() {
    configPath := filepath.Join("/home/binperdock/GOLANG/github.com/HwuuPhuc0904/backend-api/configs")
    v := viper.New() // Đổi tên biến để tránh xung đột
    v.AddConfigPath(configPath)
    v.SetConfigName("local")
    v.SetConfigType("yaml")

    // Thêm đọc cấu hình và xử lý lỗi
    err := v.ReadInConfig()
    if err != nil {
        panic(fmt.Errorf("error reading config file: %w", err))
    }

    // Unmarshal sau khi đọc thành công
    if err := v.Unmarshal(&global.Config); err != nil {
        panic(fmt.Errorf("unable to decode into struct: %w", err))
    }
}