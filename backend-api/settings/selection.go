package setting

type Config struct {
	MySQL MySQLSettings `mapstructure:"mysql"`
}


type MySQLSettings struct {
		Host     string `mapstructure:"host"`
		Port     string `mapstructure:"port"`
		Username string `mapstructure:"username"`
		Password string `mapstructure:"password"`
		Database string `mapstructure:"database"`
}

