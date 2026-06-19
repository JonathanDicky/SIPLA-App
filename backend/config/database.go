package config

import (
	"fmt"
	"log"
	"os"
	"strings"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() {
	var dsn string
	rawURL := os.Getenv("DATABASE_URL")

	if rawURL != "" {

		dsn = strings.TrimPrefix(rawURL, "mysql://")
		if strings.Contains(dsn, "?") {
			dsn += "&parseTime=True&loc=Local"
		} else {
			dsn += "?parseTime=True&loc=Local"
		}
		log.Println("--- KONEK KE DB PRODUCTION (RAILWAY) ---")
	} else {

		dsn = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			os.Getenv("DB_USER"),
			os.Getenv("DB_PASSWORD"),
			os.Getenv("DB_HOST"),
			os.Getenv("DB_PORT"),
			os.Getenv("DB_NAME"),
		)
		log.Println("--- KONEK KE DB LOKAL (LOCALHOST) ---")
	}

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Gagal konek database: ", err)
	}

	DB = db
	log.Println("Database sukses terhubung!")
}
