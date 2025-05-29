package routers
import(
	"github.com/gin-gonic/gin"
    "GOLANG/github.com/HwuuPhuc0904/backend-api/internal/middleware"
)

func SetupRouter() *gin.Engine {
    r := gin.Default() // Đã bao gồm Logger và Recovery middleware
    
    // use CORS middleware
    r.Use(middleware.CORSMiddleware())
    

    // API version group
    v1 := r.Group("/api/v1")
    
    // Đăng ký các routes
    RegisterUserRoutes(v1)
    RegisterProductRoutes(v1)
    RegisterOrderRoutes(v1)
    
    // Thêm route ping để kiểm tra server hoạt động
    r.GET("/ping", func(c *gin.Context) {
        c.JSON(200, gin.H{
            "message": "pong",
        })
    })
    
    return r
}