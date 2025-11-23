package controllers

import (
	"BangkitcellBe/config"
	"BangkitcellBe/model"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// GetAllDeviceServiceVariant - Get all device service variants
func GetAllDeviceServiceVariant(c *gin.Context) {
	var dsv []model.DeviceServiceVariant

	// Preload dengan relationship yang benar
	if err := config.DB.
		Preload("Device", func(db *gorm.DB) *gorm.DB {
			return db.Preload("Brand")
		}).
		Preload("Service").
		Find(&dsv).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    dsv,
	})
}

// GetDeviceServiceVariantById - Get device service variant by ID
func GetDeviceServiceVariantById(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid device service variant ID",
		})
		return
	}

	var dsv model.DeviceServiceVariant
	if err := config.DB.
		Preload("Device", func(db *gorm.DB) *gorm.DB {
			return db.Preload("Brand")
		}).
		Preload("Service").
		First(&dsv, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"success": false,
				"error":   "Device service variant not found",
			})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"error":   err.Error(),
			})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    dsv,
	})
}

// CreateDeviceServiceVariant - Create new device service variant
func CreateDeviceServiceVariant(c *gin.Context) {
	type Request struct {
		DeviceID  uint    `json:"device_id" binding:"required"`
		ServiceID uint    `json:"service_id" binding:"required"`
		TipePart  *string `json:"tipe_part,omitempty"`
		HargaMin  float64 `json:"harga_min" binding:"required"`
		HargaMax  float64 `json:"harga_max" binding:"required"`
		Catatan   *string `json:"catatan,omitempty"`
	}

	var req Request
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	dsv := model.DeviceServiceVariant{
		DeviceID:  req.DeviceID,
		ServiceID: req.ServiceID,
		TipePart:  req.TipePart,
		HargaMin:  req.HargaMin,
		HargaMax:  req.HargaMax,
		Catatan:   req.Catatan,
	}

	if err := config.DB.Create(&dsv).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	// Reload dengan relationships
	if err := config.DB.
		Preload("Device", func(db *gorm.DB) *gorm.DB {
			return db.Preload("Brand")
		}).
		Preload("Service").
		First(&dsv, dsv.ID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"data":    dsv,
	})
}

// UpdateDeviceServiceVariant - Update device service variant
func UpdateDeviceServiceVariant(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid device service variant ID",
		})
		return
	}

	type Request struct {
		DeviceID  uint    `json:"device_id" binding:"required"`
		ServiceID uint    `json:"service_id" binding:"required"`
		TipePart  *string `json:"tipe_part,omitempty"`
		HargaMin  float64 `json:"harga_min" binding:"required"`
		HargaMax  float64 `json:"harga_max" binding:"required"`
		Catatan   *string `json:"catatan,omitempty"`
	}

	var req Request
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	var dsv model.DeviceServiceVariant
	if err := config.DB.First(&dsv, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"success": false,
				"error":   "Device service variant not found",
			})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"error":   err.Error(),
			})
		}
		return
	}

	// Update fields
	dsv.DeviceID = req.DeviceID
	dsv.ServiceID = req.ServiceID
	dsv.TipePart = req.TipePart
	dsv.HargaMin = req.HargaMin
	dsv.HargaMax = req.HargaMax
	dsv.Catatan = req.Catatan

	if err := config.DB.Save(&dsv).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	// Reload dengan relationships
	if err := config.DB.
		Preload("Device", func(db *gorm.DB) *gorm.DB {
			return db.Preload("Brand")
		}).
		Preload("Service").
		First(&dsv, dsv.ID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    dsv,
	})
}

// DeleteDeviceServiceVariant - Delete device service variant
func DeleteDeviceServiceVariant(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid device service variant ID",
		})
		return
	}

	result := config.DB.Delete(&model.DeviceServiceVariant{}, id)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   result.Error.Error(),
		})
		return
	}

	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   "Device service variant not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Device service variant deleted successfully",
	})
}