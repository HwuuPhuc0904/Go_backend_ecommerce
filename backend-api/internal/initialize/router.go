package initialize

import (
    "GOLANG/github.com/HwuuPhuc0904/backend-api/internal/routers"
)

func InitRouters() {
	r := routers.SetupRouter()
	r.Run(":8080")

}

