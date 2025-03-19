package controller

import (
	"GOLANG/github.com/HwuuPhuc0904/backend-api/global"
	"GOLANG/github.com/HwuuPhuc0904/backend-api/internal/service"
	"net/http"
	"github.com/gin-gonic/gin"
    "go.uber.org/zap"
	
	"strconv"
)

type UserController struct {
	UserService * service.UserService
}

func NewUserController() *UserController {
	return &UserController {
		UserService: service.NewUserService(),
	}
}


// GetUserByID lấy thông tin người dùng theo ID
func (uc *UserController) GetUserByID(c *gin.Context) {
    // Lấy id từ param
    idParam := c.Param("id")
    id, err := strconv.ParseUint(idParam, 10, 32)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "success": false,
            "message": "ID không hợp lệ",
        })
        return
    }

    // Gọi service để lấy thông tin user
    user, err := uc.UserService.GetUserByID(uint(id))
    if err != nil {
        global.Logger.Error("Không tìm thấy người dùng", zap.Error(err))
        c.JSON(http.StatusNotFound, gin.H{
            "success": false,
            "message": "Người dùng không tồn tại",
        })
        return
    }

    // Trả về kết quả thành công
    c.JSON(http.StatusOK, gin.H{
        "success": true,
        "data":    user,
    })
}