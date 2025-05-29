package controller

import (
    "GOLANG/github.com/HwuuPhuc0904/backend-api/global"
    "GOLANG/github.com/HwuuPhuc0904/backend-api/internal/models"
    "GOLANG/github.com/HwuuPhuc0904/backend-api/internal/service"
    "github.com/gin-gonic/gin"
    "go.uber.org/zap"
    "net/http"
    "strconv"
)

type AddressController struct {
    addressService *service.AddressService
}

func NewAddressController() *AddressController {
    return &AddressController{
        addressService: service.NewAddressService(),
    }
}

// GetAddresses trả về danh sách địa chỉ của người dùng
func (ac *AddressController) GetAddresses(c *gin.Context) {
    userID, exists := c.Get("userID")
    if !exists {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
        return
    }

    addresses, err := ac.addressService.GetAddressesByUserID(userID.(uint))
    if err != nil {
        global.Logger.Error("Failed to get addresses", zap.Error(err))
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get addresses"})
        return
    }

    c.JSON(http.StatusOK, addresses)
}

// CreateAddress tạo địa chỉ mới cho người dùng
func (ac *AddressController) CreateAddress(c *gin.Context) {
    userID, exists := c.Get("userID")
    if !exists {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
        return
    }

    var address models.Address
    if err := c.ShouldBindJSON(&address); err != nil {
        global.Logger.Error("Invalid address data", zap.Error(err))
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    address.UserID = userID.(uint)
    if err := ac.addressService.CreateAddress(&address); err != nil {
        global.Logger.Error("Failed to create address", zap.Error(err))
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusCreated, address)
}

// GetAddress trả về thông tin một địa chỉ cụ thể
func (ac *AddressController) GetAddress(c *gin.Context) {
    userID, exists := c.Get("userID")
    if !exists {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
        return
    }

    addressID, err := strconv.ParseUint(c.Param("id"), 10, 64)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid address ID"})
        return
    }

    address, err := ac.addressService.GetAddressByID(uint(addressID), userID.(uint))
    if err != nil {
        global.Logger.Error("Failed to get address", zap.Error(err), zap.Uint64("addressID", addressID))
        c.JSON(http.StatusNotFound, gin.H{"error": "Address not found"})
        return
    }

    c.JSON(http.StatusOK, address)
}

// UpdateAddress cập nhật thông tin địa chỉ
func (ac *AddressController) UpdateAddress(c *gin.Context) {
    userID, exists := c.Get("userID")
    if !exists {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
        return
    }

    addressID, err := strconv.ParseUint(c.Param("id"), 10, 64)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid address ID"})
        return
    }

    var address models.Address
    if err := c.ShouldBindJSON(&address); err != nil {
        global.Logger.Error("Invalid address data", zap.Error(err))
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    address.ID = uint(addressID)
    address.UserID = userID.(uint)

    if err := ac.addressService.UpdateAddress(address.ID, address.UserID, &address); err != nil {
        global.Logger.Error("Failed to update address", zap.Error(err), zap.Uint64("addressID", addressID))
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Address updated successfully"})
}

// DeleteAddress xóa địa chỉ
func (ac *AddressController) DeleteAddress(c *gin.Context) {
    userID, exists := c.Get("userID")
    if !exists {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
        return
    }

    addressID, err := strconv.ParseUint(c.Param("id"), 10, 64)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid address ID"})
        return
    }

    if err := ac.addressService.DeleteAddress(uint(addressID), userID.(uint)); err != nil {
        global.Logger.Error("Failed to delete address", zap.Error(err), zap.Uint64("addressID", addressID))
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Address deleted successfully"})
}

// SetDefaultAddress thiết lập địa chỉ mặc định
func (ac *AddressController) SetDefaultAddress(c *gin.Context) {
    userID, exists := c.Get("userID")
    if !exists {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
        return
    }

    addressID, err := strconv.ParseUint(c.Param("id"), 10, 64)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid address ID"})
        return
    }

    addressType := c.DefaultQuery("type", "both")
    if addressType != "shipping" && addressType != "billing" && addressType != "both" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid address type. Must be 'shipping', 'billing', or 'both'"})
        return
    }

    if err := ac.addressService.SetDefaultAddress(uint(addressID), userID.(uint)); err != nil {
        global.Logger.Error("Failed to set default address", zap.Error(err), zap.Uint64("addressID", addressID))
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Default address set successfully"})
}

// GetProvinces returns a list of provinces/cities
func (ac *AddressController) GetProvinces(c *gin.Context) {
    provinces, err := ac.addressService.GetProvinces()
    if err != nil {
        global.Logger.Error("Failed to get provinces", zap.Error(err))
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get provinces"})
        return
    }

    c.JSON(http.StatusOK, provinces)
}

// GetDistricts returns districts for a specific province
func (ac *AddressController) GetDistricts(c *gin.Context) {
    provinceCode, err := strconv.Atoi(c.Param("provinceCode"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid province code"})
        return
    }

    districts, err := ac.addressService.GetDistricts(strconv.Itoa(provinceCode))
    if err != nil {
        global.Logger.Error("Failed to get districts", zap.Error(err))
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get districts"})
        return
    }

    c.JSON(http.StatusOK, districts)
}

