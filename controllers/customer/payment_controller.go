package customer_controller

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
	"github.com/raihansyahrin/backend_laundry_app.git/config"
	"github.com/raihansyahrin/backend_laundry_app.git/models"
)

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
