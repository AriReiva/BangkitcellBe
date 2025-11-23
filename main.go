package main

import (
	"BangkitcellBe/config"
	"BangkitcellBe/model"
	"BangkitcellBe/routes"
	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/cors"
	"time"
)

func main() {
	config.InitDB()

	config.DB.AutoMigrate(&model.Brand{})
	config.DB.AutoMigrate(&model.Device{})

	r := gin.Default()

	// === CORS FIX ===
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173", "http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
	// =================

	routes.BrandsRouter(r)
	routes.DeviceRouter(r)
	routes.TransactionRouter(r)
	routes.UserRouter(r)
	routes.AuthRouter(r)
	routes.StatsRouter(r)

	r.Run(":8000")
}
