package controller

import (
	"GOLANG/github.com/HwuuPhuc0904/backend-api/global"
	"GOLANG/github.com/HwuuPhuc0904/backend-api/internal/models"
	"GOLANG/github.com/HwuuPhuc0904/backend-api/internal/service"
	"net/http"
	"strconv"
    "time"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type UserController struct {
    userService *service.UserService
}

func NewUserController() *UserController {
    return &UserController{
        userService: service.NewUserService(),
    }
}

// Resgister user
func (uc * UserController) RegisterUser(c * gin.Context) {
    var user models.User

    if err := c.ShouldBindJSON(&user); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid data" + err.Error()})
        return
    }

    now := time.Now()
    user.CreatedAt = now
    user.UpdatedAt = now


    if err := uc.userService.CreateUser(&user); err != nil {
        global.Logger.Error("Failed to create user", zap.Error(err))
        
        if err.Error() == "email already exists" {
            c.JSON(http.StatusConflict, gin.H{"error": "email already exists"})
            return
        }
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
        return
    }

    user.Password = ""
    c.JSON(http.StatusCreated, gin.H{"msg": "User registered successfully","user": user})
}   

// login
func (uc *UserController) Login(c *gin.Context) {
    var loginData struct {
        Email   string `json:"email" binding:"required"`
        Password string `json:"password" binding:"required"`
    }

    if err := c.ShouldBindJSON(&loginData); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid data" + err.Error()})
        return
    }

    // check email and password before call authenticate
    if loginData.Email == "" || loginData.Password == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid email or password"})
        return 
    }    

    user, token, err := uc.userService.AuthenticateUser(loginData.Email, loginData.Password)

    if err != nil {
        global.Logger.Info("Failed to login user", zap.String("email", loginData.Email), zap.Error(err))
        
        switch err.Error() {
        case "user not found":
            c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
            return
        case "invalid password":
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid password"})
            return
        default:
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to login user"})
            return
        }
        
    }
    if user == nil {
        global.Logger.Error("User is nil after successful authentication")
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Intenal server error"})
        return
    }

    user.Password = ""
    c.JSON(http.StatusOK, gin.H{"msg": "User logged in successfully", "user": user, "token": token})
}

// get user by id
func (uc *UserController) GetUserByID(c *gin.Context) {
    userID := c.Param("id")
    id, err := strconv.ParseUint(userID, 10, 64)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
        return
    }

    user, err := uc.userService.GetUserByID(uint(id))

    if err != nil {
        global.Logger.Error("Failed to get user by ID", zap.Uint64("id", id), zap.Error(err))
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve user: " + err.Error()})
        return
    }

    user.Password = ""
    c.JSON(http.StatusOK, gin.H{"msg": "User retrieved successfully", "user": user})   
}

// Update profile
func (uc *UserController) UpdateProfileUser(c *gin.Context) {
    userID , exists := c.Get("userID") 
    if !exists {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "User not Unauthorized"})
        return
    }

    var userData models.User
    if err := c.ShouldBindJSON(&userData); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid data" + err.Error()})
        return
    }    

    userData.ID = userID.(uint)
    if err := uc.userService.UpdateUser(&userData); err != nil {
        global.Logger.Error("Failed to update user", zap.Uint("id", userData.ID), zap.Error(err))

        if err.Error() == "user not found" {
            c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
            return
        }


        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
        return
    }   

    user, _ := uc.userService.GetUserByID(userData.ID)
    user.Password = ""

    c.JSON(http.StatusOK, gin.H{"msg": "User updated successfully", "user": user})
}

// Change pasword
func (uc *UserController) ChangePassword(c * gin.Context) {
    userID, exists := c.Get("userID")

    if !exists {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "User not Unauthorized"})
        return
    }

    var passwordData struct {
        CurrentPassword string `json:"current_password" binding:"required"`
        NewPassword     string `json:"new_password" binding:"required"`
    }

    if err := c.ShouldBindJSON(&passwordData); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid data" + err.Error()})
        return
    }

    if err := uc.userService.ChangePassword(userID.(uint), passwordData.CurrentPassword, passwordData.NewPassword); err != nil {
        global.Logger.Error("Failed to change password", zap.Uint("id", userID.(uint)), zap.Error(err))
        if err.Error() == "invalid password" {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid password"})
            return
        }
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to change password"})
    }


    c.JSON(http.StatusOK, gin.H{"msg": "Password changed successfully"})
}

// get all user for admin
func (uc *UserController) GetAllUser(c *gin.Context) {
    pageStr := c.DefaultQuery("page", "1")
    limitStr := c.DefaultQuery("limit", "10")

    page, _ := strconv.Atoi(pageStr)
    limit, _ := strconv.Atoi(limitStr)

    users, total, err := uc.userService.GetAllUsers(page, limit)

    if err != nil {
        global.Logger.Error("Failed to get all users", zap.Error(err))
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get all users" + err.Error()})
        return
    }

    for i := range users {
        users[i].Password = ""
    }

    c.JSON(http.StatusOK, gin.H{
        "msg": "Users retrieved successfully", 
        "users": users, 
        "total": total,
        "page": page,
        "limit": limit,
    })
}

func (uc * UserController) UpdateUser(c * gin.Context) {
    idStr := c.Param("id")
    id, err := strconv.ParseUint(idStr, 10, 64) 
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
        return
    }

    var userData models.User

    if err := c.ShouldBindJSON(&userData); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid data" + err.Error()})
        return
    }

    userData.ID = uint(id)

    if err := uc.userService.UpdateUser(&userData); err != nil {
        global.Logger.Error("Failed to update user", zap.Uint("id", userData.ID), zap.Error(err))
        
        if err.Error() == "user not found" {
            c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
            return
        }
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
        return
    }

    user, _ := uc.userService.GetUserByID(userData.ID)
    user.Password = ""

    c.JSON(http.StatusOK, gin.H{"msg": "User updated successfully", "user": user})
}

func (c *UserController) ForgotPassword(ctx *gin.Context) {
    var data struct {
        Email string `json:"email" binding:"required,email"`
    }

    if err := ctx.ShouldBindJSON(&data); err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{
            "error": "Invalid email: " + err.Error(),
        })
        return
    }

    // Trong triển khai thực tế, bạn sẽ gửi email với link reset password
    // Ở đây tạm thời chỉ kiểm tra email tồn tại hay không

    _, err := c.userService.GetUserByEmail(data.Email)
    if err != nil {
        // Không trả về lỗi cụ thể để tránh rò rỉ thông tin
        ctx.JSON(http.StatusOK, gin.H{
            "message": "If your email is registered, you will receive password reset instructions",
        })
        return
    }

    // Thực tế sẽ tạo token reset password và gửi email

    ctx.JSON(http.StatusOK, gin.H{
        "message": "If your email is registered, you will receive password reset instructions",
    })
}

// ResetPassword đặt lại mật khẩu
func (c *UserController) ResetPassword(ctx *gin.Context) {
    // Trong triển khai thực tế, bạn sẽ xác thực token và đặt lại mật khẩu
    ctx.JSON(http.StatusOK, gin.H{
        "message": "Password has been reset successfully",
    })
}

func (c *UserController) GetProfile(ctx *gin.Context) {
    userID, exists := ctx.Get("userID")
    if !exists {
        ctx.JSON(http.StatusUnauthorized, gin.H{
            "error": "User not authenticated",
        })
        return
    }

    user, err := c.userService.GetUserByID(userID.(uint))
    if err != nil {
        global.Logger.Error("Failed to get user profile", zap.Error(err), zap.Any("userID", userID))
        
        if err.Error() == "user not found" {
            ctx.JSON(http.StatusNotFound, gin.H{
                "error": "User not found",
            })
            return
        }
        
        ctx.JSON(http.StatusInternalServerError, gin.H{
            "error": "Failed to retrieve user profile: " + err.Error(),
        })
        return
    }

    // Loại bỏ mật khẩu từ phản hồi
    user.Password = ""

    ctx.JSON(http.StatusOK, gin.H{
        "user": user,
    })
}

func (c *UserController) DeleteUser(ctx *gin.Context) {
    idStr := ctx.Param("id")
    id, err := strconv.ParseUint(idStr, 10, 32)
    if err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{
            "error": "Invalid user ID",
        })
        return
    }

    if err := c.userService.DeleteUser(uint(id)); err != nil {
        global.Logger.Error("Failed to delete user", zap.Error(err), zap.String("id", idStr))
        
        if err.Error() == "user not found" {
            ctx.JSON(http.StatusNotFound, gin.H{
                "error": "User not found",
            })
            return
        }
        
        ctx.JSON(http.StatusInternalServerError, gin.H{
            "error": "Failed to delete user: " + err.Error(),
        })
        return
    }

    ctx.JSON(http.StatusOK, gin.H{
        "message": "User deleted successfully",
    })
}
