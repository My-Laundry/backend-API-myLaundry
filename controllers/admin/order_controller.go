package admin_controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/raihansyahrin/backend_laundry_app.git/config"
	"github.com/raihansyahrin/backend_laundry_app.git/models"
	"github.com/raihansyahrin/backend_laundry_app.git/response"
)

func OrderComplete(c *gin.Context) {
	var body struct {
		OrderID uint `json:"order_id" form:"order_id"`
	}

	if err := c.ShouldBind(&body); err != nil {
		c.JSON(http.StatusBadRequest, response.DefaultResponse{
			Code:    http.StatusBadRequest,
			Success: false,
			Message: "Invalid input format",
			Data:    nil,
		})
		return
	}

	var order models.Order
	if err := config.DB.Preload("Service").Preload("Courier").Preload("Customer").Preload("Admin").Preload("Address").First(&order, body.OrderID).Error; err != nil {
		c.JSON(http.StatusBadRequest, response.DefaultResponse{
			Code:    http.StatusBadRequest,
			Success: false,
			Message: "Invalid order ID",
			Data:    nil,
		})
		return
	}

	adminID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, response.DefaultResponse{
			Code:    http.StatusUnauthorized,
			Success: false,
			Message: "User not authenticated",
			Data:    nil,
		})
		return
	}

	// Validasi apakah user memiliki role sebagai admin
	userRole, exists := c.Get("role")
	if !exists || userRole != "admin" {
		c.JSON(http.StatusUnauthorized, response.DefaultResponse{
			Code:    http.StatusUnauthorized,
			Success: false,
			Message: "User is not authorized as an admin",
			Data:    nil,
		})
		return
	}

	// Ubah status pesanan menjadi 'done'
	order.Status = "done"

	// Type assertion to get uint value from adminID
	adminIDUint, ok := adminID.(uint)
	if !ok {
		c.JSON(http.StatusUnauthorized, response.DefaultResponse{
			Code:    http.StatusUnauthorized,
			Success: false,
			Message: "Invalid admin ID type",
			Data:    nil,
		})
		return
	}

	order.AdminID = &adminIDUint

	if err := config.DB.Save(&order).Error; err != nil {
		c.JSON(http.StatusInternalServerError, response.DefaultResponse{
			Code:    http.StatusInternalServerError,
			Success: false,
			Message: "Failed to update order status",
			Data:    nil,
		})
		return
	}

	orderResponse := response.OrderResponse{
		ID:         order.ID,
		Status:     order.Status,
		CreatedAt:  order.CreatedAt.String(),
		UpdatedAt:  order.UpdatedAt.String(),
		TotalPrice: order.TotalPrice,
		Weight:     order.Weight,
		Quantity:   order.Quantity,
		Customer: response.UserResponse{
			ID:       order.Customer.ID,
			Username: order.Customer.Username,
			Email:    order.Customer.Email,
		},
		Admin: response.UserResponse{
			ID:       order.Admin.ID,
			Username: order.Admin.Username,
			Email:    order.Admin.Email,
		},
		Service: response.ServiceResponse{
			ID:    order.Service.ID,
			Title: order.Service.Title,
			Price: uint(order.Service.Price),
		},
		Courier: response.UserResponse{
			ID:       order.Courier.ID,
			Username: order.Courier.Username,
			Email:    order.Courier.Email,
		},
		Address: response.AddressResponse{
			ID:            order.Address.ID,
			CustomerID:    order.Address.CustomerID,
			ReceiverName:  order.Address.ReceiverName,
			PhoneNumber:   order.Address.PhoneNumber,
			HouseNumber:   order.Address.HouseNumber,
			ResidenceName: order.Address.ResidenceName,
			AddressNotes:  order.Address.AddressNotes,
			StreetName:    order.Address.StreetName,
			District:      order.Address.District,
			SubDistrict:   order.Address.SubDistrict,
			City:          order.Address.City,
			Area:          order.Address.Area,
		},
	}

	c.JSON(http.StatusOK, response.DefaultResponse{
		Code:    http.StatusOK,
		Success: true,
		Message: "Order complete",
		Data:    orderResponse,
	})
}


// ACCEPT BY ADMIN
// func AcceptOrder(c *gin.Context) {
// 	orderID := c.Param("id")

// 	var body struct {
// 		CourierID uint `json:"courier_id" form:"courier_id"`
// 	}

// 	if err := c.ShouldBind(&body); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid input format"})
// 		return
// 	}

// 	var order models.Order
// 	if err := config.DB.First(&order, orderID).Error; err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid order ID"})
// 		return
// 	}

// 	adminID, exists := c.Get("user_id")
// 	if !exists {
// 		c.JSON(http.StatusUnauthorized, gin.H{"message": "User not authenticated"})
// 		return
// 	}

// 	// // Validasi apakah user memiliki role sebagai kurir
// 	// userRole, exists := c.Get("role")
// 	// if !exists || userRole != "courier" {
// 	// 	c.JSON(http.StatusUnauthorized, gin.H{"message": "User is not authorized as a courier"})
// 	// 	return
// 	// }

// 	// Validasi apakah order sudah diterima sebelumnya
// 	if order.Status != "waiting for admin approval" {
// 		c.JSON(http.StatusConflict, gin.H{"message": "Order has already been accepted or is in progress"})
// 		return
// 	}

// 	// Konversi adminID ke uint
// 	adminIDUint, ok := adminID.(uint)
// 	if !ok {
// 		c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid admin ID type"})
// 		return
// 	}

// 	// Validasi apakah orderID terkait dengan kurir yang sedang masuk
// 	if order.CourierID != nil && *order.CourierID != adminIDUint {
// 		c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized access to this order"})
// 		return
// 	}

// 	// Tetapkan courier_id dan ubah status order
// 	order.CourierID = &body.CourierID
// 	order.Status = "Kurir On The Way"
// 	order.AdminID = adminIDUint

// 	if err := config.DB.Save(&order).Error; err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to accept order"})
// 		return
// 	}

// 	c.JSON(http.StatusOK, gin.H{"message": "Order accepted successfully"})
// }