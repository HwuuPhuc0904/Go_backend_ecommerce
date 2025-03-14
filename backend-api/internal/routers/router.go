package routers

import (
	"github.com/gin-gonic/gin"
	c "GOLANG/github.com/HwuuPhuc0904/backend-api/internal/controller"
)

func NewRouter() *gin.Engine {
	r := gin.Default()

	v1 := r.Group("v1/2024")
	{
		v1.GET("/ping", c.Pong)
		// uu điểm dễ dàng quản lý 
		v1.GET("/ping/1", c.NewUserController().GetUserByID) 
		// v1.PUT("/ping", Pong)
		// v1.PATCH("/ping", Pong)
		// v1.DELETE("/ping", Pong)
		// v1.POST("/ping", Pong)
	}

	// v2 := r.Group("v2/2024")
	// {
	// 	v2.GET("/ping/:name", Pong)
	// 	v2.PUT("/ping", Pong)
	// 	v2.PATCH("/ping", Pong)
	// 	v2.DELETE("/ping", Pong)
	// 	v2.POST("/ping", Pong)
	// }
	return r
}

