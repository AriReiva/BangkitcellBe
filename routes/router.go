package routes

import (
	"BangkitcellBe/controllers"
	"BangkitcellBe/middleware"

	"github.com/gin-gonic/gin"
)

func BrandsRouter(r *gin.RouterGroup) {
	brandsGroup := r.Group("/brands")
	{
		brandsGroup.GET("/", controllers.GetAllBrands)
		brandsGroup.GET("/:id", controllers.GetBrandByID)
		brandsGroup.POST("/", controllers.CreateBrand)
		brandsGroup.PUT("/:id", controllers.UpdateBrand)
		brandsGroup.DELETE("/:id", controllers.DeleteBrand)
	}
}

func DeviceRouter(r *gin.RouterGroup) {
	deviceGroup := r.Group("/devices")
	{
		deviceGroup.GET("/", controllers.GetAllDevice)
		deviceGroup.GET("/:id", controllers.GetDeviceById)
		deviceGroup.POST("/", controllers.CreateDevice)
		deviceGroup.PUT("/:id", controllers.UpdateDevice)
		deviceGroup.DELETE("/:id", controllers.DeleteDevice)
	}
}

func TransactionRouter(r *gin.RouterGroup) {
	transactionGroup := r.Group("/transactions")
	{
		transactionGroup.GET("/", controllers.GetAllTransaction)
		transactionGroup.GET("/:id", controllers.GetTransactionById)
		transactionGroup.POST("/", controllers.CreateTransaction)
		transactionGroup.PATCH("/:id/payment", controllers.UpdateTransaction)
		transactionGroup.DELETE("/:id", controllers.DeleteTransaction)
	}
}

func UserRouter(r *gin.RouterGroup) {
	transactionGroup := r.Group("/users")
	{
		transactionGroup.GET("/", controllers.GetAllUser)
		transactionGroup.GET("/:id", controllers.GetUserById)
		transactionGroup.POST("/", controllers.CreateUser)
		transactionGroup.PATCH("/:id/payment", controllers.UpdateUser)
		transactionGroup.PUT("/:id", controllers.UpdateUser)
		transactionGroup.DELETE("/:id", controllers.DeleteUser)
	}
}


func AuthRouter(r *gin.RouterGroup) {
	AuthGroup := r.Group("/auth")
	{
		AuthGroup.POST("/register", controllers.RegisterUser)
		AuthGroup.POST("/login", controllers.LoginUser)

		auth := AuthGroup.Group("/").Use(middleware.AuthMiddleware())
		auth.POST("/logout", controllers.LogoutUser)
		auth.GET("/me", func(c *gin.Context){
			c.JSON(200, c.MustGet("user"))
		})
	}
}

func StatsRouter(r *gin.RouterGroup) {
	g := r.Group("/stats")
	{
		g.GET("/", controllers.StatsIndex)
		g.GET("/report", controllers.StatsReport)
		g.GET("/omset", controllers.GetOmset)
	}
}

func DeviceServiceVariantRouter(r *gin.RouterGroup){
	dsvGroup := r.Group("/variants")
	{
		dsvGroup.GET("/", controllers.GetAllDeviceServiceVariant)
		// dsvGroup.GET("/:id", controllers.GetDeviceServiceByIdVariant)
		dsvGroup.POST("/", controllers.CreateDeviceServiceVariant)
		dsvGroup.PUT("/:id", controllers.UpdateDeviceServiceVariant)
		dsvGroup.DELETE("/:id", controllers.DeleteDeviceServiceVariant)
	}
}

func ServiceRouter(r *gin.RouterGroup){
	serviceGroup := r.Group("/services")
	{
		serviceGroup.GET("/", controllers.GetAllService)
		serviceGroup.GET("/:id", controllers.GetServiceById)
		serviceGroup.POST("/", controllers.CreateService)
		serviceGroup.PUT("/:id", controllers.UpdateService)
		serviceGroup.DELETE("/:id", controllers.DeleteService)
	}
}