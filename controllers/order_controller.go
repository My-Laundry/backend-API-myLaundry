package controllers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
	"github.com/raihansyahrin/backend_laundry_app.git/config"
	"github.com/raihansyahrin/backend_laundry_app.git/models"
	"github.com/raihansyahrin/backend_laundry_app.git/response"
)

func GetOrderDetailForCustomer(c *gin.Context) {

	type UserResponse struct {
		ID       uint   `json:"id"`
		Username string `json:"username"`
		Email    string `json:"email"`
	}

	customerIDStr := c.Param("customer_id")
	customerID, err := strconv.Atoi(customerIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.DefaultResponse{
			Success: false,
			Message: "Invalid customer ID",
			Code:    http.StatusBadRequest,
		})
		return
	}

	// Fetch orders by customer ID
	var orders []models.Order
	if err := config.DB.Preload("Customer").Preload("Courier").Preload("Admin").Preload("Service").Preload("Address").Where("customer_id = ?", customerID).Find(&orders).Error; err != nil {
		c.JSON(http.StatusBadRequest, response.DefaultResponse{
			Success: false,
			Message: "Invalid customer ID or no orders found",
			Code:    http.StatusBadRequest,
		})
		return
	}

	if len(orders) == 0 {
		c.JSON(http.StatusNotFound, response.DefaultResponse{
			Success: false,
			Message: "No orders found for this customer",
			Code:    http.StatusNotFound,
		})
		return
	}

	var orderResponses []response.OrderResponse
	for _, order := range orders {
		orderResponse := response.OrderResponse{
			ID:        order.ID,
			Status:    order.Status,
			CreatedAt: order.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt: order.UpdatedAt.Format("2006-01-02 15:04:05"),
			Customer: response.UserResponse{
				ID:       order.Customer.ID,
				Username: order.Customer.Username,
				Email:    order.Customer.Email,
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
			Courier: response.UserResponse{
				ID:       order.Courier.ID,
				Username: order.Courier.Username,
				Email:    order.Courier.Email,
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
		}
		orderResponses = append(orderResponses, orderResponse)
	}

	c.JSON(http.StatusOK, response.DefaultResponse{
		Success: true,
		Message: "Orders retrieved successfully",
		Code:    http.StatusOK,
		Data:    orderResponses,
	})
}

func AcceptOrder(c *gin.Context) {
	orderID := c.Param("id")

	var order models.Order
	if err := config.DB.First(&order, orderID).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid order ID"})
		return
	}

	adminID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "User not authenticated"})
		return
	}

	order.Status = "Kurir On The Way"
	order.AdminID = adminID.(uint)

	if err := config.DB.Save(&order).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to accept order"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Order accepted successfully"})
}

func CourierArrived(c *gin.Context) {
	var body struct {
		OrderID uint    `json:"order_id" form:"order_id"`
		Weight  float64 `json:"weight" form:"weight"`
	}

	if err := c.ShouldBind(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid input format"})
		return
	}

	var order models.Order
	if err := config.DB.First(&order, body.OrderID).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid order ID"})
		return
	}

	// Update weight and status
	order.Weight = body.Weight
	order.Status = "arrived"

	if err := config.DB.Save(&order).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to update order status and weight"})
		return
	}

	// Send updated order details to admin and user
	// You might want to use a real-time messaging service or email for notifications

	c.JSON(http.StatusOK, gin.H{"message": "Courier arrived, weight updated successfully", "order": order})
}
func generateQRCode(order models.Order) (string, error) {
	client := resty.New()
	// Replace with your payment gateway API URL and parameters
	resp, err := client.R().
		SetBody(map[string]interface{}{
			"amount":      order.TotalPrice,
			"description": "Payment for order " + string(order.ID),
		}).
		Post("https://api.example.com/generate_qr")
	if err != nil {
		return "", err
	}

	var result map[string]interface{}
	err = json.Unmarshal(resp.Body(), &result)
	if err != nil {
		return "", err
	}

	qrCode, ok := result["qr_code"].(string)
	if !ok {
		return "", errors.New("failed to parse QR code from response")
	}

	return qrCode, nil
}

func ProcessPayment(c *gin.Context) {
	var body struct {
		OrderID uint   `json:"order_id" form:"order_id"`
		Method  string `json:"method" form:"method"` // "cash" or "qris"
	}

	if err := c.ShouldBind(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid input format"})
		return
	}

	var order models.Order
	if err := config.DB.First(&order, body.OrderID).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid order ID"})
		return
	}

	if body.Method == "cash" {
		order.Status = "completed"
		if err := config.DB.Save(&order).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to update order status"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Order marked as paid with cash"})
	} else if body.Method == "qris" {
		qrCode, err := generateQRCode(order)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to generate QR code", "error": err.Error()})
			return
		}
		order.Status = "waiting for payment confirmation"
		if err := config.DB.Save(&order).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to update order status"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "QRIS payment initiated", "qr_code": qrCode})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid payment method"})
	}
}

func UpdateOrderStatus(c *gin.Context) {
	var body struct {
		OrderID uint   `json:"order_id" form:"order_id"`
		Status  string `json:"status" form:"status"`
	}

	if err := c.ShouldBind(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid input format"})
		return
	}

	var order models.Order
	if err := config.DB.First(&order, body.OrderID).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid order ID"})
		return
	}

	// Check if the status update is valid
	if body.Status == "in progress" && order.Status != "arrived" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid status update"})
		return
	}

	order.Status = body.Status

	if err := config.DB.Save(&order).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to update order status"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Order status updated successfully"})
}

func CreateOrder(c *gin.Context) {
	var body struct {
		ServiceID uint    `json:"service_id" form:"service_id"`
		Weight    float64 `json:"weight" form:"weight"`
		AddressID uint    `json:"address_id" form:"address_id"`
	}

	if err := c.ShouldBind(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid input format"})
		return
	}

	// Ambil ID pengguna dari token JWT
	customerID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "User not authenticated"})
		return
	}

	// Pastikan pengguna memiliki role customer
	role, exists := c.Get("role")
	if !exists || role != "customer" {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "User is not a customer or role not found"})
		return
	}

	// Ambil courier yang tersedia (misalnya yang pertama ditemukan atau sesuai logika penugasan)
	var courier models.User
	if err := config.DB.Where("role = ?", "courier").First(&courier).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "No available courier found", "error": err.Error()})
		return
	}

	// Pastikan courier tidak memiliki order aktif
	var existingOrder models.Order
	if err := config.DB.Where("courier_id = ? AND status != ?", courier.ID, "completed").First(&existingOrder).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"message": "Courier already has an active order"})
		return
	}

	// Ambil service dari database
	var service models.Service
	if err := config.DB.First(&service, body.ServiceID).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid service ID"})
		return
	}

	// Ambil address dari database
	var address models.Address
	if err := config.DB.First(&address, body.AddressID).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid address ID"})
		return
	}

	totalPrice := body.Weight * service.Price

	// Buat order baru dengan customer_id dari pengguna yang sedang login dan courier_id dari courier yang tersedia
	order := models.Order{
		CustomerID: customerID.(uint), // Set customer_id sebagai ID pengguna yang login
		CourierID:  courier.ID,        // Set courier_id dari courier yang tersedia
		AdminID:    1,                 // Ganti dengan admin_id yang valid jika diperlukan
		ServiceID:  body.ServiceID,
		AddressID:  body.AddressID,
		Weight:     body.Weight,
		TotalPrice: totalPrice,
		Status:     "waiting for pickup",
	}

	if err := config.DB.Create(&order).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to create order", "error": err.Error()})
		return
	}

	// Preload entitas terkait sebelum mengirimkan respons
	if err := config.DB.Preload("Address").Preload("Customer").Preload("Courier").Preload("Admin").Preload("Service").First(&order, order.ID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to retrieve created order with associated data", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Order created successfully", "order": order})
}

func GetOrders(c *gin.Context) {
	var orders []models.Order
	if err := config.DB.Preload("Customer").Preload("Courier").Preload("Admin").Preload("Service").Preload("Address").Find(&orders).Error; err != nil {
		c.JSON(http.StatusInternalServerError, response.DefaultResponse{
			Success: false,
			Message: "Failed to retrieve orders",
			Code:    http.StatusInternalServerError,
		})
		return
	}

	var orderResponses []response.OrderResponse
	for _, order := range orders {
		orderResponse := response.OrderResponse{
			ID:        order.ID,
			Status:    order.Status,
			CreatedAt: order.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt: order.UpdatedAt.Format("2006-01-02 15:04:05"),
			Customer: response.UserResponse{
				ID:        order.Customer.ID,
				Username:  order.Customer.Username,
				Email:     order.Customer.Email,
				Role:      nil,
				Addresses: nil,
			},
			Courier: response.UserResponse{
				ID:        order.Courier.ID,
				Username:  order.Courier.Username,
				Email:     order.Courier.Email,
				Role:      nil,
				Addresses: nil},
			Admin: response.UserResponse{
				ID:        order.Admin.ID,
				Username:  order.Admin.Username,
				Email:     order.Admin.Email,
				Role:      nil,
				Addresses: nil},
			Service: response.ServiceResponse{
				ID:    order.Service.ID,
				Title: order.Service.Title,
				Price: uint(order.Service.Price),
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
		orderResponses = append(orderResponses, orderResponse)
	}

	c.JSON(http.StatusOK, response.DefaultResponse{
		Success: true,
		Message: "Successfully retrieved orders",
		Code:    http.StatusOK,
		Data:    orderResponses,
	})
}

// func UpdateOrderStatus(c *gin.Context) {
// 	var body struct {
// 		OrderID uint   `json:"order_id" form:"order_id"`
// 		Status  string `json:"status" form:"status"`
// 	}

// 	if err := c.ShouldBind(&body); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid input format"})
// 		return
// 	}

// 	var order models.Order
// 	if err := config.DB.First(&order, body.OrderID).Error; err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid order ID"})
// 		return
// 	}

// 	order.Status = body.Status

// 	if err := config.DB.Save(&order).Error; err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to update order status"})
// 		return
// 	}

// 	c.JSON(http.StatusOK, gin.H{"message": "Order status updated successfully"})
// }

func DeleteOrder(c *gin.Context) {
	var body struct {
		OrderID uint `json:"order_id"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid input format"})
		return
	}

	if err := config.DB.Delete(&models.Order{}, body.OrderID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to delete order"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Order deleted successfully"})
}
