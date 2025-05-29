package repo

import (
    "GOLANG/github.com/HwuuPhuc0904/backend-api/global"
    "GOLANG/github.com/HwuuPhuc0904/backend-api/internal/models"
    "errors"
    "fmt"
    "go.uber.org/zap"
    "gorm.io/gorm"
)

type AddressRepo struct{}

func NewAddressRepo() *AddressRepo {
    return &AddressRepo{}
}

// CreateAddress creates a new address for a user
func (ar *AddressRepo) CreateAddress(address *models.Address) error {
    // Start a transaction for consistency
    tx := global.DB.Begin()
    defer func() {
        if r := recover(); r != nil {
            tx.Rollback()
        }
    }()

    // If this address is set as default, unset any other default addresses of the same type
    if address.IsDefault {
        if err := ar.unsetDefaultAddressesWithTx(tx, address.UserID, address.AddressType); err != nil {
            tx.Rollback()
            global.Logger.Error("Failed to unset default addresses", zap.Error(err))
            return err
        }
    }

    // Create the new address
    if err := tx.Create(address).Error; err != nil {
        tx.Rollback()
        global.Logger.Error("Failed to create address", zap.Error(err))
        return err
    }

    return tx.Commit().Error
}

// GetAddressesByUserID returns all addresses for a user
func (ar *AddressRepo) GetAddressesByUserID(userID uint) ([]models.Address, error) {
    var addresses []models.Address
    
    if err := global.DB.Where("user_id = ?", userID).
        Order("is_default DESC, created_at DESC").
        Find(&addresses).Error; err != nil {
        global.Logger.Error("Failed to get user addresses", zap.Error(err))
        return nil, err
    }
    
    return addresses, nil
}

// GetAddressByID returns a specific address if it belongs to the user
func (ar *AddressRepo) GetAddressByID(addressID, userID uint) (*models.Address, error) {
    var address models.Address
    
    if err := global.DB.Where("id = ? AND user_id = ?", addressID, userID).
        First(&address).Error; err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, errors.New("address not found or does not belong to user")
        }
        global.Logger.Error("Failed to get address", zap.Error(err))
        return nil, err
    }
    
    return &address, nil
}

// UpdateAddress updates an existing address
func (ar *AddressRepo) UpdateAddress(address *models.Address) error {
    // Start a transaction for consistency
    tx := global.DB.Begin()
    defer func() {
        if r := recover(); r != nil {
            tx.Rollback()
        }
    }()

    // If this address is being set as default, unset any other default addresses of the same type
    if address.IsDefault {
        if err := ar.unsetDefaultAddressesWithTx(tx, address.UserID, address.AddressType);err != nil {
            tx.Rollback()
            global.Logger.Error("Failed to unset default addresses", zap.Error(err))
            return err
        }
    }

    // Update the address
    if err := tx.Save(address).Error; err != nil {
        tx.Rollback()
        global.Logger.Error("Failed to update address", zap.Error(err))
        return err
    }

    return tx.Commit().Error
}

// DeleteAddress removes an address if it belongs to the user
func (ar *AddressRepo) DeleteAddress(addressID, userID uint) error {
    result := global.DB.Where("id = ? AND user_id = ?", addressID, userID).Delete(&models.Address{})
    
    if result.Error != nil {
        global.Logger.Error("Failed to delete address", zap.Error(result.Error))
        return result.Error
    }
    
    if result.RowsAffected == 0 {
        return errors.New("address not found or does not belong to user")
    }
    
    return nil
}

// SetDefaultAddress marks an address as default for a specific type
func (ar *AddressRepo) SetDefaultAddress(addressID, userID uint, addressType string) error {
    // Start a transaction for consistency
    tx := global.DB.Begin()
    defer func() {
        if r := recover(); r != nil {
            tx.Rollback()
        }
    }()

    // Unset any existing default addresses of this type
    if err := ar.unsetDefaultAddressesWithTx(tx, userID, addressType); err != nil {
        tx.Rollback()
        return err
    }

    // Set the new default
    if err := tx.Model(&models.Address{}).
        Where("id = ? AND user_id = ?", addressID, userID).
        Update("is_default", true).Error; err != nil {
        tx.Rollback()
        global.Logger.Error("Failed to set default address", zap.Error(err))
        return err
    }

    return tx.Commit().Error
}

// unsetDefaultAddressesWithTx is a helper function to unset default addresses
// It uses a transaction to ensure consistency
func (ar *AddressRepo) unsetDefaultAddressesWithTx(tx *gorm.DB, userID uint, addressType string) error {
    query := tx.Model(&models.Address{}).
        Where("user_id = ? AND is_default = ?", userID, true)
    
    // Apply filter based on address type
    switch addressType {
    case "shipping":
        query = query.Where("address_type IN (?, ?)", "shipping", "both") // Sửa từ 'type' thành 'address_type'
    case "billing":
        query = query.Where("address_type IN (?, ?)", "billing", "both") // Sửa từ 'type' thành 'address_type'
    case "both":
        query = query.Where("address_type IN (?, ?, ?)", "shipping", "billing", "both") // Sửa từ 'type' thành 'address_type'
    default:
        return fmt.Errorf("invalid address type: %s", addressType)
    }
    
    return query.Update("is_default", false).Error
}




func (ar *AddressRepo) GetAddressesByUserIDAndType(userID uint, addressType string) ([]models.Address, error) {
    var addresses []models.Address
    err := global.DB.Where("user_id = ? AND address_type = ?", userID, addressType).Find(&addresses).Error
    return addresses, err
}

func (ar *AddressRepo) ClearDefaultAddress(userID uint, addressType string) error {
    // Example using GORM (adjust according to your ORM/database logic)
    result := global.DB.Model(&models.Address{}).
        Where("user_id = ? AND address_type = ? AND is_default = ?", userID, addressType, true).
        Update("is_default", false)
    return result.Error
}