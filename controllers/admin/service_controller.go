package admin_controllers

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/raihansyahrin/backend_laundry_app.git/config"
	"github.com/raihansyahrin/backend_laundry_app.git/models"
)

type ServiceController struct{}

// GetServices mengambil semua layanan
func (sc *ServiceController) GetServices(c *gin.Context) {
	var services []models.Service
	if err := config.DB.Find(&services).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to retrieve services"})
		return
	}

	// Mengubah format data response
	var serviceResponse []gin.H
	for _, service := range services {
		serviceResponse = append(serviceResponse, gin.H{
			"id":       service.ID,
			"title":    service.Title,
			"time":     service.Time,
			"price":    service.Price,
			"category": service.Category,
		})
	}

	c.JSON(http.StatusOK, gin.H{"data": serviceResponse})
}

// GetServiceByCategory mengambil layanan berdasarkan kategori
func (sc *ServiceController) GetServiceByCategory(c *gin.Context) {
	category := c.Param("category")

	// Ubah spasi menjadi underscore pada kategori
	categoryEndpoint := strings.ReplaceAll(strings.ToLower(category), " ", "_")

	var services []models.Service
	if err := config.DB.Where("category = ?", category).Find(&services).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to retrieve services by category"})
		return
	}

	var serviceResponse []gin.H
	for _, service := range services {
		serviceResponse = append(serviceResponse, gin.H{
			"id":       service.ID,
			"title":    service.Title,
			"time":     service.Time,
			"price":    service.Price,
			"category": service.Category,
		})
	}

	c.JSON(http.StatusOK, gin.H{"data": serviceResponse, "category": categoryEndpoint, "code": 200, "success": true})
}

// GetServiceByID mengambil layanan berdasarkan ID
func (sc *ServiceController) GetServiceByID(c *gin.Context) {
	id := c.Param("id")

	var service models.Service
	if err := config.DB.First(&service, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Service not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":       service.ID,
		"title":    service.Title,
		"time":     service.Time,
		"price":    service.Price,
		"category": service.Category,
	})
}

// CreateService membuat layanan baru
func (sc *ServiceController) CreateService(c *gin.Context) {
	var service models.Service
	if err := c.ShouldBind(&service); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid input format"})
		return
	}

	if err := config.DB.Create(&service).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to create service"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"success": true,
		"message": "Service created successfully",
		"data": gin.H{
			"id":       service.ID,
			"title":    service.Title,
			"time":     service.Time,
			"price":    service.Price,
			"category": service.Category,
		},
	})
}

// UpdateService mengupdate layanan berdasarkan ID
func (sc *ServiceController) UpdateService(c *gin.Context) {
	id := c.Param("id")

	var service models.Service
	if err := config.DB.First(&service, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Service not found"})
		return
	}

	var updatedService models.Service
	if err := c.ShouldBindJSON(&updatedService); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid input format"})
		return
	}

	service.Title = updatedService.Title
	service.Time = updatedService.Time
	service.Price = updatedService.Price
	service.Category = updatedService.Category // tambahkan update untuk category

	if err := config.DB.Save(&service).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to update service"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Service updated successfully", "service": service})
}

// DeleteService menghapus layanan berdasarkan ID
func (sc *ServiceController) DeleteService(c *gin.Context) {
	id := c.Param("id")

	if err := config.DB.Delete(&models.Service{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to delete service"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Service deleted successfully"})
}
