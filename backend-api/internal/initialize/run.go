package initialize

import(
	"GOLANG/github.com/HwuuPhuc0904/backend-api/global"
    "GOLANG/github.com/HwuuPhuc0904/backend-api/migrations"
    "go.uber.org/zap"
)

func Run() {
	LoadConfig()
	InitLogger()
	InitMysql()

	if err := migrations.MigrateDB(); err != nil {
        global.Logger.Error("Migration failed", zap.Error(err))
        panic(err)
    }
    
	InitRedis()
	InitRouters()
}