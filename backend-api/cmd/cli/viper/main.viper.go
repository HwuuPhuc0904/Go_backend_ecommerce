package main

import (
	"fmt"
    "path/filepath"
	"github.com/spf13/viper"
)

type Config struct {
	Server struct {
		Port int `mapstructure:"port"`
	} `mapstructure:"server"`
	Databases struct {
		User string `mapstructure:"user"`
		Password string `mapstructure:"password"`
		Host string `mapstructure:"host"`
	} `mapstructure:"databases"`
}


func main() {
	configPath := filepath.Join("/home/binperdock/GOLANG/github.com/HwuuPhuc0904/backend-api/configs")
	viper := viper.New()
	viper.AddConfigPath(configPath)
	viper.SetConfigName("production")
	viper.SetConfigType("yaml")

	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("error reading config file: %w", err))
	}
	//read server
	fmt.Println(viper.GetInt("server.port"))

	var config Config 
	//viper.Unmarshal(&config) là một phương thức của thư viện Viper dùng để chuyển đổi dữ liệu cấu hình đã đọc từ file (như YAML, JSON, TOML...) thành một struct trong Go.
	if err := viper.Unmarshal(&config); err != nil {
		panic(fmt.Errorf("unable to decode into struct, %w", err))
	}
	println("ConfigPort :",config.Server.Port)

}