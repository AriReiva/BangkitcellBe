package controllers

import (
	"BangkitcellBe/config"
	"BangkitcellBe/model"
	"BangkitcellBe/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetAllDevice(c *gin.Context){
	var device [] model.Device

	if err := config.DB.Preload("Brand").Find(&device).Error;
	err != nil {
		utils.RespondError(c, http.StatusInternalServerError, err)
		return
	}

	utils.RespondSuccess(c, device)
}

func GetDeviceById(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		utils.RespondError(c, http.StatusBadRequest, err)
		return
	}

	var device model.Device

	// Preload relasi
	if err := config.DB.
		Preload("Services").
		Preload("Pivots").
		Preload("Brand").
		First(&device, id).
		Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.RespondError(c, http.StatusNotFound, err)
		} else {
			utils.RespondError(c, http.StatusInternalServerError, err)
		}
		return
	}
	for i := range device.Services {
		for j := range device.Pivots {
			if device.Pivots[j].ServiceID == device.Services[i].ID {
				device.Services[i].Pivot = &device.Pivots[j]
			}
		}
	}

	utils.RespondSuccess(c, device)
}


func CreateDevice(c *gin.Context) {
	var device model.Device
	if err := c.ShouldBindJSON(&device); err != nil {
		utils.RespondError(c, http.StatusBadRequest, err)
		return
	}
	if err := config.DB.Create(&device).Error; err != nil {
		utils.RespondError(c, http.StatusInternalServerError, err)
		return
	}
	utils.RespondSuccess(c, device)
}

func UpdateDevice(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		utils.RespondError(c, http.StatusBadRequest, err)
		return
	}
	var device model.Device
	if err := config.DB.First(&device, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.RespondError(c, http.StatusNotFound, err)
		}
		return
	}
	if err := c.ShouldBindJSON(&device); err != nil {
		utils.RespondError(c, http.StatusBadRequest, err)
		return
	}
	if err := config.DB.Save(&device).Error; err != nil {
		utils.RespondError(c, http.StatusInternalServerError, err)
		return
	}
	utils.RespondSuccess(c, device)
}

func DeleteDevice(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		utils.RespondError(c, http.StatusBadRequest, err)
		return
	}
	if err := config.DB.Delete(&model.Device{}, id).Error; err != nil {
		utils.RespondError(c, http.StatusInternalServerError, err)
		return
	}
	utils.RespondSuccess(c, "Device deleted successfully")
}