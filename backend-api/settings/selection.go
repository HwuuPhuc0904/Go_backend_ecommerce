package setting

type Config struct {
	MySQL MySQLSettings `mapstructure:"mysql"`
	LoggerSetting LoggerSetting `mapstructure:"logger"`
}


type MySQLSettings struct {
		Host     string `mapstructure:"host"`
		Port     string `mapstructure:"port"`
		Username string `mapstructure:"username"`
		Password string `mapstructure:"password"`
		Database string `mapstructure:"database"`
}

type LoggerSetting struct {
	Level string `mapstructure:"level"`
	MaxSize int `mapstructure:"maxsize"`
	MaxBackup int `mapstructure:"maxbackup"`
	MaxAge int `mapstructure:"maxage"`
	Compress bool `mapstructure:"compress"`
	FilePath string `mapstructure:"filepath"`
}