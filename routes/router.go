package routes

import (
	"BangkitcellBe/controllers"

	"github.com/gin-gonic/gin"
)

func BrandsRouter(r *gin.Engine) {
	brandsGroup := r.Group("/brands")
	{
		brandsGroup.GET("/", controllers.GetAllBrands)
		brandsGroup.GET("/:id", controllers.GetBrandByID)
		brandsGroup.POST("/", controllers.CreateBrand)
		brandsGroup.PUT("/:id", controllers.UpdateBrand)
		brandsGroup.DELETE("/:id", controllers.DeleteBrand)
	}
}

func DeviceRouter(r *gin.Engine) {
	deviceGroup := r.Group("/device")
	{
		deviceGroup.GET("/", controllers.GetAllDevice)
		deviceGroup.GET("/:id", controllers.GetDeviceById)
		deviceGroup.POST("/", controllers.CreateDevice)
		deviceGroup.PUT("/:id", controllers.UpdateDevice)
		deviceGroup.DELETE("/:id", controllers.DeleteDevice)
	}
}

func TransactionRouter(r *gin.Engine) {
	transactionGroup := r.Group("/transactions")
	{
		transactionGroup.GET("/", controllers.GetAllTransaction)
		transactionGroup.GET("/:id", controllers.GetTransactionById)
		transactionGroup.POST("/", controllers.CreateTransaction)
		transactionGroup.PATCH("/:id/payment", controllers.UpdateTransaction)
		transactionGroup.DELETE("/:id", controllers.DeleteTransaction)
	}
}

func StatsRouter(r *gin.Engine) {

	g := r.Group("/stats")
	{
		g.GET("/", controllers.StatsIndex)
		g.GET("/report", controllers.StatsReport)
		g.GET("/omset", controllers.GetOmset)
	}
}

