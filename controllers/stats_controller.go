package controllers

import (
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"BangkitcellBe/config"
	"BangkitcellBe/model"
	"BangkitcellBe/utils"
)

type StatsController struct {
	DB *gorm.DB
}

func NewStatsController() *StatsController {
	return &StatsController{DB: config.DB}
}

////////////////////////////////////////////////////
// ================ INDEX =========================
////////////////////////////////////////////////////
func StatsIndex(c *gin.Context) {
    ctrl := NewStatsController()
    db := ctrl.DB

    var totalServices, totalBrands, totalDevices, totalTransactions, totalUsers int64
    var totalRevenue float64

    // Error handling sederhana
    if err := db.Model(&model.Service{}).Count(&totalServices).Error; err != nil {
        c.JSON(500, gin.H{"error": "Failed to count services: " + err.Error()})
        return
    }

    if err := db.Model(&model.Brand{}).Count(&totalBrands).Error; err != nil {
        c.JSON(500, gin.H{"error": "Failed to count brands: " + err.Error()})
        return
    }

    if err := db.Model(&model.Device{}).Count(&totalDevices).Error; err != nil {
        c.JSON(500, gin.H{"error": "Failed to count devices: " + err.Error()})
        return
    }

    if err := db.Model(&model.Transaction{}).Where("status = ?", "pending").Count(&totalTransactions).Error; err != nil {
        c.JSON(500, gin.H{"error": "Failed to count transactions: " + err.Error()})
        return
    }

    if err := db.Model(&model.User{}).Count(&totalUsers).Error; err != nil {
        c.JSON(500, gin.H{"error": "Failed to count users: " + err.Error()})
        return
    }

    if err := db.Model(&model.Transaction{}).
        Where("status = ?", "success").
        Select("COALESCE(sum(total),0)").Scan(&totalRevenue).Error; err != nil {
        c.JSON(500, gin.H{"error": "Failed to calculate revenue: " + err.Error()})
        return
    }

    // Sales Data
    var salesData []struct {
        Date         time.Time `json:"date"`
        TotalRevenue float64   `json:"total_revenue"`
    }
    if err := db.Raw(`
        SELECT DATE(created_at) AS date, SUM(total) AS total_revenue
        FROM transactions
        GROUP BY DATE(created_at)
        ORDER BY DATE(created_at)
    `).Scan(&salesData).Error; err != nil {
        c.JSON(500, gin.H{"error": "Failed to get sales data: " + err.Error()})
        return
    }

    // Device distribution
    var deviceDistribution []struct {
        Name  string `json:"name"`
        Value int64  `json:"value"`
    }
    if err := db.Raw(`
        SELECT brands.nama AS name, COUNT(devices.id) AS value
        FROM devices
        JOIN brands ON devices.brand_id = brands.id
        GROUP BY brands.id, brands.nama
    `).Scan(&deviceDistribution).Error; err != nil {
        c.JSON(500, gin.H{"error": "Failed to get device distribution: " + err.Error()})
        return
    }

    // Recent activities (1 month)
    var recentActivities []model.Transaction
    if err := db.Where("created_at >= ?", time.Now().AddDate(0, -1, 0)).
        Order("created_at DESC").
        Find(&recentActivities).Error; err != nil {
        c.JSON(500, gin.H{"error": "Failed to get recent activities: " + err.Error()})
        return
    }

    c.JSON(200, gin.H{
        "total_services":      totalServices,
        "total_brands":        totalBrands,
        "total_devices":       totalDevices,
        "total_transactions":  totalTransactions,
        "total_users":         totalUsers,
        "total_revenue":       totalRevenue,
        "sales_data":          salesData,
        "device_distribution": deviceDistribution,
        "recent_activities":   recentActivities,
    })
}

////////////////////////////////////////////////////
// ================ REPORT ========================
////////////////////////////////////////////////////
func StatsReport(c *gin.Context) {
	ctrl := NewStatsController()
	db := ctrl.DB

	var totalRevenue float64
	var totalTransactions int64

	db.Model(&model.Transaction{}).
		Where("status = ?", "success").
		Select("COALESCE(sum(total),0)").Scan(&totalRevenue)

	db.Model(&model.Transaction{}).
		Where("status = ?", "success").
		Count(&totalTransactions)

	// Sales per month
	var salesData []struct {
		Year         int    `json:"year"`
		MonthNumber  int    `json:"month_number"`
		Month        string `json:"month"`
		Revenue      float64 `json:"revenue"`
		Transactions int64  `json:"transactions"`
	}
	db.Raw(`
		SELECT 
			YEAR(created_at) AS year,
			MONTH(created_at) AS month_number,
			DATE_FORMAT(created_at, "%b %Y") AS month,
			SUM(total) AS revenue,
			COUNT(*) AS transactions
		FROM transactions
		WHERE status = "success"
		GROUP BY year, month_number, month
		ORDER BY year, month_number
	`).Scan(&salesData)

	// Service Performance
	var servicePerformance []struct {
		Service      string  `json:"service"`
		Revenue      float64 `json:"revenue"`
		Transactions int64   `json:"transactions"`
	}
	db.Raw(`
		SELECT 
			services.nama AS service,
			SUM(transaction_details.harga) AS revenue,
			COUNT(transaction_details.id) AS transactions
		FROM transaction_details
		JOIN device_service_variants 
			ON device_service_variants.id = transaction_details.device_service_variant_id
		JOIN services 
			ON services.id = device_service_variants.service_id
		JOIN transactions 
			ON transactions.id = transaction_details.transaction_id
		WHERE transactions.status = "success"
		GROUP BY services.nama
		ORDER BY revenue DESC
		LIMIT 10
	`).Scan(&servicePerformance)

	utils.RespondSuccess(c, gin.H{
		"salesData":         salesData,
		"servicePerformance": servicePerformance,
		"totalRevenue":      totalRevenue,
		"totalTransactions": totalTransactions,
	})
}

////////////////////////////////////////////////////
// ================ GET OMSET ========================
////////////////////////////////////////////////////
func GetOmset(c *gin.Context) {
	ctrl := NewStatsController()
	db := ctrl.DB

	var omset float64
	oneMonthAgo := time.Now().AddDate(0, -1, 0)

	db.Model(&model.Transaction{}).
		Where("status = ?", "success").
		Where("DATE(created_at) = DATE(?)", oneMonthAgo).
		Select("COALESCE(sum(total),0)").Scan(&omset)

	utils.RespondSuccess(c, gin.H{"omset": omset})
}
