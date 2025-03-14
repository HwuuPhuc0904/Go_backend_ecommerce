package response

import (
	"net/http"
	"github.com/gin-gonic/gin"
)

type ResponseData struct {
// Các thẻ `json:"code"`, `json:"message"`, `json:"data"`: Đây là struct tags trong Go, chỉ định cách các trường này được chuyển đổi thành JSON khi serialize.
// Khi struct được encode thành JSON, trường Code sẽ trở thành "code", Message thành "message", và Data thành "data".
	Code    int         `json:"code"` // ma code
	Message string      `json:"message"` // thong bao
	Data    interface{} `json:"data"` // du lieu 
} 


// sucess response
func SuccessResponse(c* gin.Context, code int, data interface{}) {
	c.JSON(http.StatusOK, ResponseData{
		Code: code,
		Message: msg[code],
		Data: data,
	})
}

// error response
func ErrorResponse(c * gin.Context, code int, message string){
	c.JSON(http.StatusOK, ResponseData{
		Code: code,
		Message: msg[code],
		Data: nil,
	})
}
