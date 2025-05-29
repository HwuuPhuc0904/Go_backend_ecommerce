package service

import (
	"GOLANG/github.com/HwuuPhuc0904/backend-api/global" // Giả sử đường dẫn import này đúng
	"GOLANG/github.com/HwuuPhuc0904/backend-api/internal/models"
	"GOLANG/github.com/HwuuPhuc0904/backend-api/internal/repo" // Giả sử đường dẫn import này đúng
	"errors"
	"fmt"
	"go.uber.org/zap"
	"time"
)

// LocationData, ProvinceData, DistrictData, WardData giữ nguyên như trước
// (vẫn dùng string cho Code và các mã liên quan)
type LocationData struct {
	Code     string `json:"code"`
	Name     string `json:"name"`
	Division string `json:"division_type"`
	CodeName string `json:"codename"`
}

type ProvinceData struct {
	LocationData
	PhoneCode string `json:"phone_code"`
}

type DistrictData struct {
	LocationData
	ProvinceCode string `json:"province_code"`
}

type WardData struct {
	LocationData
	DistrictCode string `json:"district_code"`
}


type AddressService struct {
	addressRepo *repo.AddressRepo
}

func NewAddressService() *AddressService {
	return &AddressService{
		addressRepo: repo.NewAddressRepo(),
	}
}

// CreateAddress creates a new address for a user
func (as *AddressService) CreateAddress(address *models.Address) error {
	// Validate required fields
	if address.Phone == "" {
		return errors.New("số điện thoại là bắt buộc (phone is required)")
	}
	if address.StreetAddress == "" {
		return errors.New("địa chỉ đường phố là bắt buộc (street address is required)")
	}

	// Mã và Tên cho Tỉnh/Huyện/Xã đều là bắt buộc từ client
	if address.ProvinceCode == "" {
		return errors.New("mã tỉnh/thành phố là bắt buộc (province code is required)")
	}
	if address.ProvinceName == "" {
		return errors.New("tên tỉnh/thành phố là bắt buộc (province name is required)")
	}
	if address.DistrictCode == "" {
		return errors.New("mã quận/huyện là bắt buộc (district code is required)")
	}
	if address.DistrictName == "" {
		return errors.New("tên quận/huyện là bắt buộc (district name is required)")
	}
	if address.WardCode == "" {
		return errors.New("mã phường/xã là bắt buộc (ward code is required)")
	}
	if address.WardName == "" {
		return errors.New("tên phường/xã là bắt buộc (ward name is required)")
	}
    if address.Country == "" { // Country cũng là not null
        address.Country = "Việt Nam" // Gán giá trị mặc định nếu client không gửi
    }


	// Validate location data (code and name consistency)
	if err := as.validateLocationData(address); err != nil {
		global.Logger.Error("Dữ liệu địa phương không hợp lệ (Invalid location data)", zap.Error(err))
		return err
	}

	// Set timestamps
	now := time.Now()
	address.CreatedAt = now
	address.UpdatedAt = now

	existingAddresses, err := as.addressRepo.GetAddressesByUserID(address.UserID)
	if err != nil {
		global.Logger.Error("Không thể lấy danh sách địa chỉ của người dùng (Failed to get user addresses)", zap.Error(err))
		return err
	}

	if len(existingAddresses) == 0 {
		address.IsDefault = true
	} else {
		if address.IsDefault {
			if err := as.addressRepo.ClearDefaultAddress(address.UserID, address.AddressType); err != nil {
				global.Logger.Error("Không thể bỏ đánh dấu địa chỉ mặc định cũ (Failed to clear old default address)", zap.Error(err))
				// Cân nhắc trả lỗi hoặc tiếp tục
			}
		}
	}

	return as.addressRepo.CreateAddress(address)
}

// GetAddressesByUserID, GetAddressByID giữ nguyên

func (as *AddressService) GetAddressesByUserID(userID uint) ([]models.Address, error) {
	return as.addressRepo.GetAddressesByUserID(userID)
}

func (as *AddressService) GetAddressByID(addressID, userID uint) (*models.Address, error) {
	return as.addressRepo.GetAddressByID(addressID, userID)
}


// UpdateAddress updates an existing address
func (as *AddressService) UpdateAddress(addressID uint, userID uint, updatedAddressInput *models.Address) error {
	existing, err := as.addressRepo.GetAddressByID(addressID, userID)
	if err != nil {
		return err // Address not found or does not belong to user
	}

	// Biến để theo dõi xem có cần xác thực lại dữ liệu vị trí không
	locationFieldsChanged := false

	// Cập nhật các trường cơ bản
	if updatedAddressInput.RecipientName != "" {
		existing.RecipientName = updatedAddressInput.RecipientName
	}
	// ... (các trường khác tương tự: CompanyName, Phone, AlternatePhone, Email, StreetAddress, PostalCode, DeliveryNotes, Country, Latitude, Longitude) ...
    if updatedAddressInput.CompanyName != "" {
		existing.CompanyName = updatedAddressInput.CompanyName
	}
	if updatedAddressInput.Phone != "" {
		existing.Phone = updatedAddressInput.Phone
	}
    if updatedAddressInput.AlternatePhone != "" {
		existing.AlternatePhone = updatedAddressInput.AlternatePhone
	}
    if updatedAddressInput.Email != "" {
		existing.Email = updatedAddressInput.Email
	}
	if updatedAddressInput.StreetAddress != "" {
		existing.StreetAddress = updatedAddressInput.StreetAddress
	}
    if updatedAddressInput.PostalCode != "" {
		existing.PostalCode = updatedAddressInput.PostalCode
	}
    if updatedAddressInput.DeliveryNotes != "" {
		existing.DeliveryNotes = updatedAddressInput.DeliveryNotes
	}
    if updatedAddressInput.Country != "" {
        existing.Country = updatedAddressInput.Country
    }
    if updatedAddressInput.Latitude != 0 { // Cân nhắc dùng con trỏ cho các trường có thể có giá trị 0 hợp lệ
        existing.Latitude = updatedAddressInput.Latitude
    }
    if updatedAddressInput.Longitude != 0 {
        existing.Longitude = updatedAddressInput.Longitude
    }


	// Xử lý cập nhật thông tin vị trí (Province, District, Ward)
	// Nếu bất kỳ mã hoặc tên nào được cung cấp để cập nhật, chúng ta cần đảm bảo tất cả đều có mặt và hợp lệ.
	if updatedAddressInput.ProvinceCode != "" || updatedAddressInput.ProvinceName != "" ||
		updatedAddressInput.DistrictCode != "" || updatedAddressInput.DistrictName != "" ||
		updatedAddressInput.WardCode != "" || updatedAddressInput.WardName != "" {

		locationFieldsChanged = true

		// Ưu tiên dữ liệu mới nếu có, nếu không giữ lại dữ liệu cũ để validate
		if updatedAddressInput.ProvinceCode != "" { existing.ProvinceCode = updatedAddressInput.ProvinceCode }
		if updatedAddressInput.ProvinceName != "" { existing.ProvinceName = updatedAddressInput.ProvinceName }
		if updatedAddressInput.DistrictCode != "" { existing.DistrictCode = updatedAddressInput.DistrictCode }
		if updatedAddressInput.DistrictName != "" { existing.DistrictName = updatedAddressInput.DistrictName }
		if updatedAddressInput.WardCode != "" { existing.WardCode = updatedAddressInput.WardCode }
		if updatedAddressInput.WardName != "" { existing.WardName = updatedAddressInput.WardName }

		// Kiểm tra lại các trường bắt buộc sau khi cập nhật
		if existing.ProvinceCode == "" || existing.ProvinceName == "" ||
			existing.DistrictCode == "" || existing.DistrictName == "" ||
			existing.WardCode == "" || existing.WardName == "" {
			return errors.New("khi cập nhật thông tin địa phương, cả mã và tên cho tỉnh, huyện, xã đều là bắt buộc")
		}
	}

	if updatedAddressInput.AddressType != "" {
		if existing.AddressType != updatedAddressInput.AddressType && updatedAddressInput.IsDefault {
            // Nếu loại địa chỉ thay đổi VÀ đang được đặt làm mặc định,
            // cần xóa mặc định cho loại cũ và loại mới của người dùng.
			if err := as.addressRepo.ClearDefaultAddress(userID, existing.AddressType); err != nil {
				 global.Logger.Error("Không thể bỏ đánh dấu địa chỉ mặc định cũ (loại cũ)", zap.Error(err), zap.String("oldAddressType", existing.AddressType))
			}
        }
		existing.AddressType = updatedAddressInput.AddressType
	}


	// Validate location data if any location field was intended for update
	if locationFieldsChanged {
		if err := as.validateLocationData(existing); err != nil {
			global.Logger.Error("Dữ liệu địa phương cập nhật không hợp lệ (Updated location data is invalid)", zap.Error(err))
			return err
		}
	}

	// Handle IsDefault
	if updatedAddressInput.IsDefault != existing.IsDefault { // Chỉ xử lý nếu có sự thay đổi về IsDefault
        if updatedAddressInput.IsDefault { // Đang muốn đặt làm mặc định
            if err := as.addressRepo.ClearDefaultAddress(userID, existing.AddressType); err != nil {
                global.Logger.Error("Không thể bỏ đánh dấu địa chỉ mặc định cũ khi cập nhật (Failed to clear old default address on update)", zap.Error(err))
                // Quyết định có nên trả lỗi cứng không
            }
            existing.IsDefault = true
        } else { // Đang muốn bỏ mặc định
            existing.IsDefault = false
            // Cân nhắc: nếu không còn địa chỉ mặc định nào khác cùng loại, có nên cảnh báo/bắt buộc đặt cái khác?
        }
    }


	existing.UpdatedAt = time.Now()
	return as.addressRepo.UpdateAddress(existing)
}


// DeleteAddress, SetDefaultAddress giữ nguyên
func (as *AddressService) DeleteAddress(addressID, userID uint) error {
	addressToDelete, err := as.addressRepo.GetAddressByID(addressID, userID)
	if err != nil {
		return err
	}

	if addressToDelete.IsDefault {
		otherAddresses, err := as.addressRepo.GetAddressesByUserIDAndType(userID, addressToDelete.AddressType)
		if err != nil {
			global.Logger.Error("Không thể lấy địa chỉ khác để đặt làm mặc định", zap.Error(err))
		} else {
			newDefaultSet := false
			for _, addr := range otherAddresses {
				if addr.ID != addressToDelete.ID {
					// Gọi SetDefaultAddress của service thay vì repo trực tiếp để đảm bảo logic nhất quán
					if err := as.SetDefaultAddress(addr.ID, userID); err != nil { // Sửa lại: SetDefaultAddress của service
						global.Logger.Error("Không thể đặt địa chỉ mới làm mặc định sau khi xóa", zap.Error(err), zap.Uint("newDefaultAddressID", addr.ID))
					} else {
						newDefaultSet = true
						global.Logger.Info("Địa chỉ mới đã được đặt làm mặc định sau khi xóa", zap.Uint("newDefaultAddressID", addr.ID))
						break
					}
				}
			}
			if !newDefaultSet && len(otherAddresses) > 1 {
                 global.Logger.Warn("Không thể tự động đặt địa chỉ mặc định mới sau khi xóa", zap.Uint("deletedAddressID", addressID))
            }
		}
	}

	return as.addressRepo.DeleteAddress(addressID, userID)
}

func (as *AddressService) SetDefaultAddress(addressID, userID uint) error {
	addressToSetDefault, err := as.addressRepo.GetAddressByID(addressID, userID)
	if err != nil {
		return err
	}

	if err := as.addressRepo.ClearDefaultAddress(userID, addressToSetDefault.AddressType); err != nil {
		global.Logger.Error("Không thể bỏ đánh dấu địa chỉ mặc định cũ", zap.Error(err))
		return err
	}
	// Sử dụng repo để cập nhật trạng thái IsDefault
	return as.addressRepo.SetDefaultAddress(addressID, userID, addressToSetDefault.AddressType)

}


// GetProvinces, GetDistricts, GetWards giữ nguyên (dùng để front-end lấy danh sách)
// Chúng vẫn trả về danh sách với mã và tên để front-end có thể hiển thị cho người dùng chọn.

func (as *AddressService) GetProvinces() ([]ProvinceData, error) {
	// Trong thực tế, điều này sẽ gọi một API hoặc DB
	provinces := []ProvinceData{
		{
			LocationData: LocationData{Code: "01", Name: "Thành phố Hà Nội", Division: "thành phố trung ương", CodeName: "thanh_pho_ha_noi"},
			PhoneCode: "024",
		},
		{
			LocationData: LocationData{Code: "79", Name: "Thành phố Hồ Chí Minh", Division: "thành phố trung ương", CodeName: "thanh_pho_ho_chi_minh"},
			PhoneCode: "028",
		},
	}
	return provinces, nil
}

func (as *AddressService) GetDistricts(provinceCode string) ([]DistrictData, error) {
	if provinceCode == "01" { // Hà Nội
		return []DistrictData{
			{LocationData: LocationData{Code: "001", Name: "Quận Ba Đình", Division: "quận", CodeName: "quan_ba_dinh"}, ProvinceCode: "01"},
			{LocationData: LocationData{Code: "002", Name: "Quận Hoàn Kiếm", Division: "quận", CodeName: "quan_hoan_kiem"}, ProvinceCode: "01"},
		}, nil
	}
	if provinceCode == "79" { // TP. Hồ Chí Minh
		return []DistrictData{
			{LocationData: LocationData{Code: "760", Name: "Quận 1", Division: "quận", CodeName: "quan_1"}, ProvinceCode: "79"},
			{LocationData: LocationData{Code: "761", Name: "Quận 12", Division: "quận", CodeName: "quan_12"}, ProvinceCode: "79"}, // Mã 761 là ví dụ
		}, nil
	}
	return []DistrictData{}, nil
}

func (as *AddressService) GetWards(districtCode string) ([]WardData, error) {
	if districtCode == "001" { // Ba Đình - Hà Nội
		return []WardData{
			{LocationData: LocationData{Code: "00001", Name: "Phường Phúc Xá", Division: "phường", CodeName: "phuong_phuc_xa"}, DistrictCode: "001"},
			{LocationData: LocationData{Code: "00004", Name: "Phường Trúc Bạch", Division: "phường", CodeName: "phuong_truc_bach"}, DistrictCode: "001"},
		}, nil
	}
    if districtCode == "760" { // Quận 1 - TP.HCM
        return []WardData{
            {LocationData: LocationData{Code: "26734", Name: "Phường Bến Nghé", Division: "phường", CodeName: "phuong_ben_nghe"}, DistrictCode: "760"},
        }, nil
    }
	return []WardData{}, nil
}


// validateLocationData validates the consistency of provided location codes and names.
func (as *AddressService) validateLocationData(address *models.Address) error {
	// Validate Province
	if address.ProvinceCode != "" { // Mã tỉnh được cung cấp
		provinces, err := as.GetProvinces() // Lấy danh sách chuẩn
		if err != nil {
			return fmt.Errorf("không thể lấy danh sách tỉnh/thành phố để xác thực: %w", err)
		}
		foundProvince := false
		for _, pData := range provinces {
			if pData.Code == address.ProvinceCode {
				foundProvince = true
				if pData.Name != address.ProvinceName {
					return fmt.Errorf("tên tỉnh/thành phố '%s' không khớp với tên chuẩn '%s' cho mã '%s'", address.ProvinceName, pData.Name, address.ProvinceCode)
				}
				break
			}
		}
		if !foundProvince {
			return fmt.Errorf("mã tỉnh/thành phố không hợp lệ: '%s'", address.ProvinceCode)
		}
	}

	// Validate District
	if address.DistrictCode != "" && address.ProvinceCode != "" {
		districts, err := as.GetDistricts(address.ProvinceCode)
		if err != nil {
			return fmt.Errorf("không thể lấy danh sách quận/huyện cho tỉnh '%s' để xác thực: %w", address.ProvinceCode, err)
		}
		foundDistrict := false
		for _, dData := range districts {
			if dData.Code == address.DistrictCode {
				foundDistrict = true
				if dData.Name != address.DistrictName {
					return fmt.Errorf("tên quận/huyện '%s' không khớp với tên chuẩn '%s' cho mã '%s' (thuộc tỉnh '%s')", address.DistrictName, dData.Name, address.DistrictCode, address.ProvinceCode)
				}
				break
			}
		}
		if !foundDistrict {
			return fmt.Errorf("mã quận/huyện '%s' không hợp lệ hoặc không thuộc tỉnh '%s'", address.DistrictCode, address.ProvinceCode)
		}
	}

	// Validate Ward
	if address.WardCode != "" && address.DistrictCode != "" {
		wards, err := as.GetWards(address.DistrictCode)
		if err != nil {
			return fmt.Errorf("không thể lấy danh sách phường/xã cho quận/huyện '%s' để xác thực: %w", address.DistrictCode, err)
		}
		foundWard := false
		for _, wData := range wards {
			if wData.Code == address.WardCode {
				foundWard = true
				if wData.Name != address.WardName {
					return fmt.Errorf("tên phường/xã '%s' không khớp với tên chuẩn '%s' cho mã '%s' (thuộc quận/huyện '%s')", address.WardName, wData.Name, address.WardCode, address.DistrictCode)
				}
				break
			}
		}
		if !foundWard {
			return fmt.Errorf("mã phường/xã '%s' không hợp lệ hoặc không thuộc quận/huyện '%s'", address.WardCode, address.DistrictCode)
		}
	}
	return nil
}