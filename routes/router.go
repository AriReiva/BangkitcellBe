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
