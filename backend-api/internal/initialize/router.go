package initialize

import (
    "GOLANG/github.com/HwuuPhuc0904/backend-api/internal/routers"
	"GOLANG/github.com/HwuuPhuc0904/backend-api/internal/middleware"
)

func InitRouters() {
	r := routers.SetupRouter()
	middleware.GenerateAPIDocumentation(r)
	r.Run(":8080")

}

