package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

type Aspirasi struct {
	ID         int    `json:"id"`
	IdUser     int    `json:"id_user"`
	IdKategori int    `json:"id_kategori"`
	Deskripsi  string `json:"deskripsi"`
	Foto       string `json:"foto"`
	Status     string `json:"status"`
}

func main() {
	// Koneksi Database
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "root:@tcp(127.0.0.1:3306)/sipla"
	}

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	app := fiber.New()

	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept",
	}))

	app.Static("/assets/pengaduan", "./assets/pengaduan")

	// Ambil Base URL secara dinamis
	getBaseURL := func() string {
		baseURL := os.Getenv("APP_URL")
		if baseURL == "" {
			return "http://localhost:8080"
		}
		return baseURL
	}

	// GET ASPIRASI
	app.Get("/aspirasi", func(c *fiber.Ctx) error {
		rows, err := db.Query("SELECT id, id_user, id_kategori, deskripsi, foto, status FROM aspirasi")
		if err != nil {
			return c.Status(500).SendString(err.Error())
		}
		defer rows.Close()

		var aspirasis []Aspirasi
		baseURL := getBaseURL()

		for rows.Next() {
			var a Aspirasi
			if err := rows.Scan(&a.ID, &a.IdUser, &a.IdKategori, &a.Deskripsi, &a.Foto, &a.Status); err != nil {
				return c.Status(500).SendString(err.Error())
			}
			// Gabungkan dengan domain dinamis
			a.Foto = fmt.Sprintf("%s/assets/pengaduan/%s", baseURL, a.Foto)
			aspirasis = append(aspirasis, a)
		}

		return c.JSON(aspirasis)
	})

	// POST ASPIRASI
	app.Post("/aspirasi", func(c *fiber.Ctx) error {
		idUser := c.FormValue("id_user")
		idKategori := c.FormValue("id_kategori")
		deskripsi := c.FormValue("deskripsi")
		status := c.FormValue("status")

		file, err := c.FormFile("foto")
		var filename string
		if err == nil {
			filename = fmt.Sprintf("%d_%s", time.Now().UnixNano(), file.Filename)
			if err := c.SaveFile(file, fmt.Sprintf("./assets/pengaduan/%s", filename)); err != nil {
				return c.Status(500).SendString("Gagal menyimpan gambar")
			}
		}

		_, err = db.Exec("INSERT INTO aspirasi (id_user, id_kategori, deskripsi, foto, status) VALUES (?, ?, ?, ?, ?)",
			idUser, idKategori, deskripsi, filename, status)
		if err != nil {
			return c.Status(500).SendString(err.Error())
		}

		baseURL := getBaseURL()
		return c.JSON(fiber.Map{
			"message": "Aspirasi berhasil ditambahkan",
			"foto":    fmt.Sprintf("%s/assets/pengaduan/%s", baseURL, filename),
		})
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Fatal(app.Listen(":" + port))
}
