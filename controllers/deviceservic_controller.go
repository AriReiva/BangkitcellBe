package controllers

import(
	"BangkitcellBe/config"
	"BangkitcellBe/model"
	"net/http"
	"strconv"
	"BangkitcellBe/utils"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetAllDeviceService(c *gin.Context){
	var dsv []model.DeviceServiceVariant

	if err:= config.DB.First(&dsv).Error;
	err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	utils.RespondSuccess(c, dsv)
}


func GetDeviceServiceById(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		utils.RespondError(c, http.StatusBadRequest, err)
		return
	}
	var dsv model.DeviceServiceVariant
	if err := config.DB.First(&dsv, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.RespondError(c, http.StatusNotFound, err)
		} else {
			utils.RespondError(c, http.StatusInternalServerError, err)
		}
		return
	}
	utils.RespondSuccess(c, dsv)
}

func CreateDeviceService(c *gin.Context) {
	var dsv model.DeviceServiceVariant
	if err := c.ShouldBindJSON(&dsv); err != nil {
		utils.RespondError(c, http.StatusBadRequest, err)
		return
	}
	if err := config.DB.Create(&dsv).Error; err != nil {
		utils.RespondError(c, http.StatusInternalServerError, err)
		return
	}
	utils.RespondSuccess(c, dsv)
}

func UpdateDeviceService(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid brand ID"})
		return
	}
	var device model.Device
	if err := config.DB.First(&device, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Brand not found"})
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
	c.JSON(http.StatusOK, device)
}

func DeleteDeviceService(c *gin.Context){
	idParam := c.Param("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid brand ID"})
		return
	}
	if err := config.DB.Delete(&model.DeviceServiceVariant{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Brand deleted successfully"})
}
