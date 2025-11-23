package routes

import (
	"BangkitcellBe/controllers"
	"BangkitcellBe/middleware"

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
	deviceGroup := r.Group("/devices")
	{
		deviceGroup.GET("/", controllers.GetAllDevice)
		deviceGroup.GET("/:id", controllers.GetDeviceById)
		deviceGroup.POST("", controllers.CreateDevice)
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

func UserRouter(r *gin.Engine) {
	transactionGroup := r.Group("/users")
	{
		transactionGroup.GET("/", controllers.GetAllUser)
		transactionGroup.GET("/:id", controllers.GetUserById)
		transactionGroup.POST("/", controllers.CreateUser)
		transactionGroup.PATCH("/:id/payment", controllers.UpdateUser)
		transactionGroup.DELETE("/:id", controllers.DeleteUser)
	}
}


func AuthRouter(r *gin.Engine) {
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

func StatsRouter(r *gin.Engine) {
    statsGroup := r.Group("/stats")
    {
        statsGroup.GET("/", controllers.StatsIndex)          // GET /stats
        // statsGroup.GET("/report", controllers.StatsReport)   // GET /stats/report
        // statsGroup.GET("/omset", controllers.GetOmset)       // GET /stats/omset
    }
}