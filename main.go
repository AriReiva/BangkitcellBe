package main

import (
	"BangkitcellBe/config"
	"BangkitcellBe/model"
	"BangkitcellBe/routes"
	"github.com/gin-gonic/gin"
)

func main() {
	config.InitDB()

	config.DB.AutoMigrate(&model.Brand{})
	config.DB.AutoMigrate(&model.Device{})
	config.DB.AutoMigrate(&model.Service{})
	config.DB.AutoMigrate(&model.DeviceServiceVariant{})

	r := gin.Default()

	apiGroup := r.Group("/api")
	{
		routes.BrandsRouter(apiGroup)
		routes.DeviceRouter(apiGroup)
		routes.TransactionRouter(apiGroup)
		routes.DeviceServiceVariantRouter(apiGroup)
		routes.ServiceRouter(apiGroup)
		routes.UserRouter(apiGroup)
		routes.AuthRouter(apiGroup)
		routes.StatsRouter(apiGroup)
	}

	
	r.Run(":8000")
}