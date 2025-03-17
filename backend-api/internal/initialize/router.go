package initialize

import (
	c "GOLANG/github.com/HwuuPhuc0904/backend-api/internal/controller"
	"GOLANG/github.com/HwuuPhuc0904/backend-api/internal/middleware"
	"fmt"

	"github.com/gin-gonic/gin"
)

func AA() gin.HandlerFunc {
	return func(c* gin.Context){
		fmt.Println("Before --> AA")
		c.Next()
		fmt.Println("After --> AA")
	}
}
func BB() gin.HandlerFunc {
	return func(c* gin.Context){
		fmt.Println("Before --> BB")
		c.Next()
		fmt.Println("After --> BB")
	}
}
func CC(c * gin.Context)  {
		fmt.Println("Before --> AA")
		c.Next()
		fmt.Println("After --> AA")
}



func InitRouter() *gin.Engine {
	r := gin.Default()

	// use middleware
	r.Use(middleware.AuthMiddleware(), BB() , CC)

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

