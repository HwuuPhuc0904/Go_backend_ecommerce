package middleware

import (
    "GOLANG/github.com/HwuuPhuc0904/backend-api/global"
    "GOLANG/github.com/HwuuPhuc0904/backend-api/pkg/utils"
    "net/http"
    "strings"

    "github.com/gin-gonic/gin"
    "go.uber.org/zap"
)

// AuthMiddleware kiểm tra người dùng đã đăng nhập hay chưa
func AuthMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        authHeader := c.GetHeader("Authorization")
        if authHeader == "" {
            c.JSON(http.StatusUnauthorized, gin.H{
                "error": "Authorization header is required",
            })
            c.Abort()
            return
        }

        // Kiểm tra định dạng Bearer token
        parts := strings.Split(authHeader, " ")
        if len(parts) != 2 || parts[0] != "Bearer" {
            c.JSON(http.StatusUnauthorized, gin.H{
                "error": "Authorization header format must be Bearer {token}",
            })
            c.Abort()
            return
        }

        tokenString := parts[1]
        claims, err := utils.ParseJWT(tokenString)
        if err != nil {
            global.Logger.Info("Invalid token", zap.Error(err))
            c.JSON(http.StatusUnauthorized, gin.H{
                "error": "Invalid or expired token",
            })
            c.Abort()
            return
        }

        // Lưu thông tin người dùng vào context
        c.Set("userID", claims.UserID)
        c.Set("userEmail", claims.Email)
        c.Set("userRole", claims.Role)

        c.Next()
    }
}

// AdminMiddleware kiểm tra người dùng có quyền admin hay không
func AdminMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        role, exists := c.Get("userRole")
        if !exists {
            c.JSON(http.StatusUnauthorized, gin.H{
                "error": "User not authenticated",
            })
            c.Abort()
            return
        }

        if role != "admin" {
            c.JSON(http.StatusForbidden, gin.H{
                "error": "Permission denied: admin role required",
            })
            c.Abort()
            return
        }

        c.Next()
    }
}