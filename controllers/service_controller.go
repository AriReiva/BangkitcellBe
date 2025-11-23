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

// GET ALL SERVICE (dengan optional search)
func GetAllService(c *gin.Context) {
	var services []model.Service
	search := c.Query("search")

	query := config.DB.Model(&model.Service{})

	if search != "" {
		query = query.Where("nama LIKE ?", "%"+search+"%")
	}

	if err := query.Find(&services).Error; err != nil {
		utils.RespondError(c, http.StatusInternalServerError, err)
		return
	}

	utils.RespondSuccess(c, services)
}

// GET SERVICE BY ID (dengan relasi Variants dan Device)
func GetServiceById(c *gin.Context) {
	idParam := c.Param("id")

	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		utils.RespondError(c, http.StatusBadRequest, err)
		return
	}

	var service model.Service

	if err := config.DB.
		Preload("Variants.Device").
		First(&service, id).Error; err != nil {

		if err == gorm.ErrRecordNotFound {
			utils.RespondError(c, http.StatusNotFound, err)
		} else {
			utils.RespondError(c, http.StatusInternalServerError, err)
		}
		return
	}

	utils.RespondSuccess(c, service)
}

// CREATE SERVICE
func CreateService(c *gin.Context) {
	var service model.Service

	if err := c.ShouldBindJSON(&service); err != nil {
		utils.RespondError(c, http.StatusBadRequest, err)
		return
	}

	if err := config.DB.Create(&service).Error; err != nil {
		utils.RespondError(c, http.StatusInternalServerError, err)
		return
	}

	utils.RespondSuccess(c, service)
}

// UPDATE SERVICE
func UpdateService(c *gin.Context) {
	idParam := c.Param("id")

	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		utils.RespondError(c, http.StatusBadRequest, err)
		return
	}

	var service model.Service

	if err := config.DB.First(&service, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.RespondError(c, http.StatusNotFound, err)
		}
		return
	}

	if err := c.ShouldBindJSON(&service); err != nil {
		utils.RespondError(c, http.StatusBadRequest, err)
		return
	}

	if err := config.DB.Save(&service).Error; err != nil {
		utils.RespondError(c, http.StatusInternalServerError, err)
		return
	}

	utils.RespondSuccess(c, service)
}

// DELETE SERVICE
func DeleteService(c *gin.Context) {
	idParam := c.Param("id")

	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		utils.RespondError(c, http.StatusBadRequest, err)
		return
	}

	if err := config.DB.Delete(&model.Service{}, id).Error; err != nil {
		utils.RespondError(c, http.StatusInternalServerError, err)
		return
	}

	utils.RespondSuccess(c, "Service deleted successfully")
}
