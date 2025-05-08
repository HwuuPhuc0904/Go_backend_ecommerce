package middleware

import (
    "GOLANG/github.com/HwuuPhuc0904/backend-api/global"
    "github.com/gin-contrib/cors"
    "github.com/gin-gonic/gin"
    "time"
)

// CORSMiddleware thiết lập middleware CORS từ cấu hình
func CORSMiddleware() gin.HandlerFunc {
    corsConfig := global.Config.CORS
    // Thiết lập giá trị mặc định nếu không có cấu hình
    if len(corsConfig.AllowOrigins) == 0 {
        corsConfig.AllowOrigins = []string{"*"}
    }
    
    if len(corsConfig.AllowMethods) == 0 {
        corsConfig.AllowMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"}
    }
    
    if len(corsConfig.AllowHeaders) == 0 {
        corsConfig.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type", "Authorization"}
    }
    
    maxAge := time.Duration(corsConfig.MaxAge) * time.Second
    if corsConfig.MaxAge == 0 {
        maxAge = 12 * time.Hour
    }
    
    return cors.New(cors.Config{
        AllowOrigins:     corsConfig.AllowOrigins,
        AllowMethods:     corsConfig.AllowMethods,
        AllowHeaders:     corsConfig.AllowHeaders,
        ExposeHeaders:    corsConfig.ExposeHeaders,
        AllowCredentials: corsConfig.AllowCredentials,
        MaxAge:           maxAge,
    })
}