package middleware

import (
	"GOLANG/github.com/HwuuPhuc0904/backend-api/global"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"strings" 
	"time"
)

// CORSMiddleware thiết lập middleware CORS từ cấu hình
func CORSMiddleware() gin.HandlerFunc {
	corsConfig := global.Config.CORS
	// Thiết lập giá trị mặc định nếu không có cấu hình	
	if len(corsConfig.AllowOrigins) == 0 {
		corsConfig.AllowOrigins = []string{"*"} // Hoặc cụ thể: []string{"http://localhost:3000", "https://your-ngrok-domain.ngrok-free.app"}
	}

	if len(corsConfig.AllowMethods) == 0 {
		corsConfig.AllowMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"}
	}

	// Đây là chỗ cần sửa đổi
	defaultAllowedHeaders := []string{"Origin", "Content-Length", "Content-Type", "Authorization"}
	if len(corsConfig.AllowHeaders) == 0 {
		corsConfig.AllowHeaders = defaultAllowedHeaders
	}
    // Thêm 'ngrok-skip-browser-warning' ào danh sách được phépv
    // Kiểm tra để tránh thêm trùng lặp nếu đã có sẵn (dù ít khả năng)
    headerAlreadyExists := false
    for _, header := range corsConfig.AllowHeaders {
        // So sánh không phân biệt chữ hoa chữ thường cho tên header
        if strings.EqualFold(header, "ngrok-skip-browser-warning") {
            headerAlreadyExists = true
            break
        }
    }
    if !headerAlreadyExists {
        corsConfig.AllowHeaders = append(corsConfig.AllowHeaders, "ngrok-skip-browser-warning")
    }


	maxAge := time.Duration(corsConfig.MaxAge) * time.Second
	if corsConfig.MaxAge == 0 {
		maxAge = 12 * time.Hour
	}

	return cors.New(cors.Config{
		AllowOrigins:     corsConfig.AllowOrigins,
		AllowMethods:     corsConfig.AllowMethods,
		AllowHeaders:     corsConfig.AllowHeaders, // Sử dụng danh sách đã được cập nhật
		ExposeHeaders:    corsConfig.ExposeHeaders,
		AllowCredentials: corsConfig.AllowCredentials,
		MaxAge:           maxAge,
	})
}