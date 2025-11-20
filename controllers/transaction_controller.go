package controllers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"BangkitcellBe/config"
	"BangkitcellBe/model"
	"BangkitcellBe/utils"
)

type CreateTransactionRequest struct {
	IdOperator       uint     `json:"id_operator" binding:"required"`
	CustomerName     string   `json:"customer_name" binding:"required"`
	CustomerPhone    *string  `json:"customer_phone"`
	Keluhan          *string  `json:"keluhan"`
	MetodePembayaran string   `json:"metode_pembayaran"`
	JumlahBayar      *float64 `json:"jumlah_bayar"`
	Kembalian        *float64 `json:"kembalian"`
	QrisReference    *string  `json:"qris_reference"`
	Details          []struct {
		DeviceServiceVariantID int64   `json:"device_service_variant_id" binding:"required"`
		Harga                  float64 `json:"harga" binding:"required"`
	} `json:"details" binding:"required,min=1"`
}

type UpdateStatusRequest struct {
	Status           *string  `json:"status"`
	JumlahBayar      *float64 `json:"jumlah_bayar"`
	Kembalian        *float64 `json:"kembalian"`
	QrisReference    *string  `json:"qris_reference"`
	MetodePembayaran *string  `json:"metode_pembayaran"`
}

type TransactionController struct {
	DB *gorm.DB
}

func NewTransactionController() *TransactionController {
	return &TransactionController{DB: config.DB}
}

////////////////////////////////////////////////////
// ================ GET ALL =======================
////////////////////////////////////////////////////
func GetAllTransaction(c *gin.Context) {
	ctrl := NewTransactionController()
	var trxs []model.Transaction
	query := ctrl.DB.Preload("Details.Variant.Device").Preload("Details.Variant.Service")

	if status := c.Query("status"); status != "" {
		query = query.Where("status = ?", status)
	}

	if err := query.Find(&trxs).Error; err != nil {
		utils.RespondError(c, http.StatusInternalServerError, err)
		return
	}

	utils.RespondSuccess(c, trxs)
}

////////////////////////////////////////////////////
// =============== GET BY ID ======================
////////////////////////////////////////////////////
func GetTransactionById(c *gin.Context) {
	ctrl := NewTransactionController()

	id := c.Param("id")
	var trx model.Transaction

	if err := ctrl.DB.Preload("Details.Variant.Device").
		Preload("Details.Variant.Service").
		First(&trx, id).Error; err != nil {
		utils.RespondError(c, http.StatusNotFound, err)
		return
	}

	utils.RespondSuccess(c, trx)
}

////////////////////////////////////////////////////
// =============== CREATE =========================
////////////////////////////////////////////////////
func CreateTransaction(c *gin.Context) {
	ctrl := NewTransactionController()
	var req CreateTransactionRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondError(c, http.StatusBadRequest, err)
		return
	}

	trx := model.Transaction{
		IDOperator:       int64(req.IdOperator),
		CustomerName:     req.CustomerName,
		CustomerPhone:    utils.StrPtrOrDefault(req.CustomerPhone),
		MetodePembayaran: req.MetodePembayaran,
		JumlahBayar:      utils.FloatPtrOrDefault(req.JumlahBayar),
		Kembalian:        utils.FloatPtrOrDefault(req.Kembalian),
		QrisReference:    utils.StrPtrOrDefault(req.QrisReference),
		Status:           "pending",
		Total:            0,
		CreatedAt:        time.Now(),
	}

	err := ctrl.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&trx).Error; err != nil {
			return err
		}

		total := 0.0

		for _, d := range req.Details {
			detail := model.TransactionDetail{
				TransactionID:          trx.ID,
				DeviceServiceVariantID: d.DeviceServiceVariantID,
				Harga:                  d.Harga,
			}

			if err := tx.Create(&detail).Error; err != nil {
				return err
			}

			total += d.Harga
		}

		trx.Total = total

		return tx.Save(&trx).Error
	})

	if err != nil {
		utils.RespondError(c, http.StatusInternalServerError, err)
		return
	}

	_ = ctrl.DB.Preload("Details.Variant.Device").
		Preload("Details.Variant.Service").
		First(&trx, trx.ID)

	utils.RespondSuccess(c, trx)
}

////////////////////////////////////////////////////
// =============== UPDATE =========================
////////////////////////////////////////////////////
func UpdateTransaction(c *gin.Context) {
	ctrl := NewTransactionController()

	id := c.Param("id")
	var trx model.Transaction

	if err := ctrl.DB.First(&trx, id).Error; err != nil {
		utils.RespondError(c, http.StatusNotFound, err)
		return
	}

	var req UpdateStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondError(c, http.StatusBadRequest, err)
		return
	}

	updates := map[string]interface{}{}

	if req.Status != nil {
		updates["status"] = *req.Status
	}
	if req.JumlahBayar != nil {
		updates["jumlah_bayar"] = *req.JumlahBayar
	}
	if req.Kembalian != nil {
		updates["kembalian"] = *req.Kembalian
	}
	if req.QrisReference != nil {
		updates["qris_reference"] = *req.QrisReference
	}
	if req.MetodePembayaran != nil {
		updates["metode_pembayaran"] = *req.MetodePembayaran
	}

	if err := ctrl.DB.Model(&trx).Updates(updates).Error; err != nil {
		utils.RespondError(c, http.StatusInternalServerError, err)
		return
	}

	utils.RespondSuccess(c, trx)
}

////////////////////////////////////////////////////
// =============== DELETE =========================
////////////////////////////////////////////////////
func DeleteTransaction(c *gin.Context) {
	ctrl := NewTransactionController()
	id := c.Param("id")

	if err := ctrl.DB.Delete(&model.Transaction{}, id).Error; err != nil {
		utils.RespondError(c, http.StatusInternalServerError, err)
		return
	}

	utils.RespondMessage(c, "Transaction deleted successfully.")
}
