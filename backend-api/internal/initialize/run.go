package initialize


func Run() {
	LoadConfig()
	InitLogger()
	InitMysql()
	InitRedis()

	// Init Router
    r := InitRouter()
	r.Run(":8002")
}