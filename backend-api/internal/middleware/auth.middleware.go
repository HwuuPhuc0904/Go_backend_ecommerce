package middleware

import (

	"GOLANG/github.com/HwuuPhuc0904/backend-api/pkg/response"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc{
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if token != "valid-token"{
			response.ErrorResponse(c, response.TokenInvalid,"invalid token")
			c.Abort()
			return
		}
		c.Next()
	}
}

