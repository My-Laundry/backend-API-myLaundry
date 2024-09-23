package customer_controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/raihansyahrin/backend_laundry_app.git/config"
	"github.com/raihansyahrin/backend_laundry_app.git/models"
	"github.com/raihansyahrin/backend_laundry_app.git/response"
)

type AddressController struct{}

// CreateAddress membuat alamat baru untuk pengguna tertentu
func (ac *AddressController) CreateAddress(c *gin.Context) {
	var address models.Address
	if err := c.ShouldBind(&address); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid input format"})
		return
	}

	// Mengambil customer_id dari konteks pengguna yang sedang login
	customerID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "User not authenticated"})
		return
	}

	// Mengisi customer_id
	address.CustomerID = customerID.(uint)

	if err := config.DB.Create(&address).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to create address"})
		return
	}

	// Buat instansi AddressResponse dan isi dengan data
	addressResponse := response.AddressResponse{
		ID:            address.ID,
		CustomerID:    address.CustomerID,
		ReceiverName:  address.ReceiverName,
		PhoneNumber:   address.PhoneNumber,
		HouseNumber:   address.HouseNumber,
		ResidenceName: address.ResidenceName,
		AddressNotes:  address.AddressNotes,
		StreetName:    address.StreetName,
		District:      address.District,
		SubDistrict:   address.SubDistrict,
		City:          address.City, // Gunakan nilai City dari model Address
		Area:          address.Area, // Gunakan nilai Area dari model Address
	}

	c.JSON(http.StatusOK, response.DefaultResponse{
		Success: true,
		Message: "Address created successfully",
		Code:    http.StatusOK,
		Data:    addressResponse,
	})
}

func (ac *AddressController) GetAddressesByUserID(c *gin.Context) {
	userID := c.Param("user_id")

	var addresses []models.Address
	if err := config.DB.Where("customer_id = ?", userID).Find(&addresses).Error; err != nil {
		c.JSON(http.StatusInternalServerError, response.DefaultResponse{
			Success: false,
			Message: "Failed to retrieve addresses",
			Code:    http.StatusInternalServerError,
		})
		return
	}

	var addressResponses []response.AddressResponse
	for _, address := range addresses {
		addressResponse := response.AddressResponse{
			ID:            address.ID,
			CustomerID:    address.CustomerID,
			ReceiverName:  address.ReceiverName,
			PhoneNumber:   address.PhoneNumber,
			HouseNumber:   address.HouseNumber,
			ResidenceName: address.ResidenceName,
			AddressNotes:  address.AddressNotes,
			StreetName:    address.StreetName,
			District:      address.District,
			SubDistrict:   address.SubDistrict,
			City:          "Bandung",
			Area:          "Bojongsoang",
		}
		if address.City != "" {
			addressResponse.City = address.City
		}
		if address.Area != "" {
			addressResponse.Area = address.Area
		}
		addressResponses = append(addressResponses, addressResponse)
	}

	c.JSON(http.StatusOK, response.DefaultResponse{
		Success: true,
		Message: "Successfully retrieved addresses",
		Code:    http.StatusOK,
		Data:    addressResponses,
	})
}

// UpdateAddress mengupdate alamat berdasarkan ID
func (ac *AddressController) UpdateAddress(c *gin.Context) {
	addressID := c.Param("id")

	var address models.Address
	if err := config.DB.First(&address, addressID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Address not found"})
		return
	}

	if err := c.ShouldBindJSON(&address); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid input format"})
		return
	}

	if err := config.DB.Save(&address).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to update address"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Address updated successfully", "data": address})
}

// DeleteAddress menghapus alamat berdasarkan ID
func (ac *AddressController) DeleteAddress(c *gin.Context) {
	addressID := c.Param("id")

	if err := config.DB.Delete(&models.Address{}, addressID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to delete address"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Address deleted successfully"})
}
