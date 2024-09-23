package admin_controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/raihansyahrin/backend_laundry_app.git/config"
	"github.com/raihansyahrin/backend_laundry_app.git/models"
	"github.com/raihansyahrin/backend_laundry_app.git/response"
	"golang.org/x/crypto/bcrypt"
)

// GetAdmins retrieves all admins
func GetAdmins(c *gin.Context) {
	var admins []models.User
	if err := config.DB.Where("role = ?", "admin").Find(&admins).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to retrieve admins"})
		return
	}

	var adminResponses []response.UserResponse
	for _, admin := range admins {
		adminResponse := response.UserResponse{
			ID:       admin.ID,
			Username: admin.Username,
			Email:    admin.Email,
		}
		adminResponses = append(adminResponses, adminResponse)
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Successfully retrieved admins",
		"code":    http.StatusOK,
		"data":    adminResponses,
	})
}

// GetAdmin retrieves a single admin based on ID
func GetAdmin(c *gin.Context) {
	id := c.Param("id")

	var admin models.User
	if err := config.DB.Where("role = ? AND id = ?", "admin", id).First(&admin).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Admin not found"})
		return
	}

	adminResponse := response.UserResponse{
		ID:       admin.ID,
		Username: admin.Username,
		Email:    admin.Email,
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Successfully retrieved admin profile",
		"code":    http.StatusOK,
		"data":    adminResponse,
	})
}

// UpdateAdmin updates an admin's profile based on ID
func UpdateAdmin(c *gin.Context) {
	id := c.Param("id")

	var body struct {
		Username string `json:"username" form:"username"`
		Email    string `json:"email" form:"email"`
		Password string `json:"password" form:"password"`
	}

	if err := c.ShouldBind(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid input format"})
		return
	}

	var admin models.User
	if err := config.DB.Where("role = ? AND id = ?", "admin", id).First(&admin).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Admin not found"})
		return
	}

	admin.Username = body.Username
	admin.Email = body.Email
	if body.Password != "" {
		hash, err := bcrypt.GenerateFromPassword([]byte(body.Password), bcrypt.DefaultCost)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Error hashing password"})
			return
		}
		admin.Password = string(hash)
	}

	if err := config.DB.Save(&admin).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to update admin"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Admin updated successfully"})
}

// DeleteAdmin deletes an admin based on ID
func DeleteAdmin(c *gin.Context) {
	id := c.Param("id")

	if err := config.DB.Where("role = ? AND id = ?", "admin", id).Delete(&models.User{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to delete admin"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Admin deleted successfully"})
}
