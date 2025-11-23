package controllers

import (
	"BangkitcellBe/config"
	"BangkitcellBe/model"
	"net/http"
	"strconv"
	"BangkitcellBe/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetAllDevice(c *gin.Context) {
	var devices []model.Device

	if err := config.DB.Find(&devices).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":   devices,
	})
}

func GetDeviceById(c *gin.Context) {
    idParam := c.Param("id")
    id, err := strconv.ParseInt(idParam, 10, 64)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
        return
    }

    var device model.Device
    if err := config.DB.Preload("Brand").First(&device, id).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Device not found"})
        return
    }

    // Load pivots
    var pivots []model.DeviceServiceVariant
    config.DB.Where("device_id = ?", device.ID).Find(&pivots)
    device.Pivots = pivots

    // Load services
    var services []model.Service
    config.DB.Joins(
        "JOIN device_service_variants dsv ON dsv.service_id = services.id",
    ).Where("dsv.device_id = ?", device.ID).
        Find(&services)

    device.Services = services

    // inject pivot ke service (loop)
    for i := range device.Services {
        for j := range device.Pivots {
            if device.Pivots[j].ServiceID == uint(device.Services[i].ID) {
                device.Services[i].Pivot = &device.Pivots[j]
            }
        }
    }

    utils.RespondSuccess(c, device)
}


func CreateDevice(c *gin.Context) {
	var device model.Device
	if err := c.ShouldBindJSON(&device); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := config.DB.Create(&device).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	utils.RespondSuccess(c, device)
}

func UpdateDevice(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid device ID"})
		return
	}
	var device model.Device
	if err := config.DB.First(&device, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Device not found"})
		}
		return
	}
	if err := c.ShouldBindJSON(&device); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := config.DB.Save(&device).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	utils.RespondSuccess(c, device)
}

func DeleteDevice(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid device ID"})
		return
	}
	if err := config.DB.Delete(&model.Device{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	utils.RespondSuccess(c, gin.H{"message": "Device deleted successfully"})
}