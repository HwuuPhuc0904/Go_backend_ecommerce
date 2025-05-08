package setting

type Config struct {
	MySQL MySQLSettings `mapstructure:"mysql"`
	LoggerSetting LoggerSetting `mapstructure:"logger"`
	CORS CORSSetting `mapstructure:"cors"`
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

type CORSSetting struct {
    AllowOrigins     []string `mapstructure:"allow_origins"`
    AllowMethods     []string `mapstructure:"allow_methods"`
    AllowHeaders     []string `mapstructure:"allow_headers"`
    ExposeHeaders    []string `mapstructure:"expose_headers"`
    AllowCredentials bool     `mapstructure:"allow_credentials"`
    MaxAge           int      `mapstructure:"max_age"`
}
