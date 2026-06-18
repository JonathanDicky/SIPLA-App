package main

import (
	"log"
	"os"
	"pengaduan/config"
	"pengaduan/routes"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	config.InitDB()
	r := gin.Default()
	r.Use(cors.New(cors.Config{
		// AllowOrigins:     []string{"http://localhost:5173"}, // Hanya izinkan FrontEnd/ atau FE di port 5173
		AllowAllOrigins: true, // Izinkan akses dari Frontend FE semua port
		AllowMethods:    []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:    []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:   []string{"Content-Length"},
		// AllowCredentials: true, // Jika hanya port tertentu uncomment yang ini
		AllowCredentials: false, // Jika semua port pakai yang ini
	}))
	r.Static("/assets", "./assets")
	routes.SetupRoutes(r)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Server berjalan di port %s", port)
	r.Run(":" + port)
}
