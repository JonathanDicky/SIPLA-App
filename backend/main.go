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

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

type Kategori struct {
	ID   int    `json:"id"`
	Nama string `json:"nama"`
}

type Aspirasi struct {
	ID         int    `json:"id"`
	IdUser     int    `json:"id_user"`
	IdKategori int    `json:"id_kategori"`
	Deskripsi  string `json:"deskripsi"`
	Foto       string `json:"foto"`
	Status     string `json:"status"`
}

func main() {
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

	// ==================== ENDPOINT LOGIN & REGISTER (KODE ASLI LU) ====================
	app.Post("/register", func(c *fiber.Ctx) error {
		var user User
		if err := c.BodyParser(&user); err != nil {
			return c.Status(400).SendString(err.Error())
		}
		user.Role = "masyarakat"
		_, err = db.Exec("INSERT INTO user (username, password, role) VALUES (?, ?, ?)", user.Username, user.Password, user.Role)
		if err != nil {
			return c.Status(500).SendString(err.Error())
		}
		return c.JSON(user)
	})

	app.Post("/login", func(c *fiber.Ctx) error {
		var user User
		if err := c.BodyParser(&user); err != nil {
			return c.Status(400).SendString(err.Error())
		}
		err = db.QueryRow("SELECT id, role FROM user WHERE username = ? AND password = ?", user.Username, user.Password).Scan(&user.ID, &user.Role)
		if err != nil {
			return c.Status(401).SendString("Username atau password salah")
		}
		return c.JSON(user)
	})

	app.Get("/masyarakat", func(c *fiber.Ctx) error {
		rows, err := db.Query("SELECT id, username, role FROM user WHERE role = 'masyarakat'")
		if err != nil {
			return c.Status(500).SendString(err.Error())
		}
		defer rows.Close()

		var users []User
		for rows.Next() {
			var u User
			if err := rows.Scan(&u.ID, &u.Username, &u.Role); err != nil {
				return c.Status(500).SendString(err.Error())
			}
			users = append(users, u)
		}
		return c.JSON(users)
	})

	app.Get("/kategori", func(c *fiber.Ctx) error {
		rows, err := db.Query("SELECT id, nama_kategori FROM kategori")
		if err != nil {
			return c.Status(500).SendString(err.Error())
		}
		defer rows.Close()

		var kategoris []Kategori
		for rows.Next() {
			var k Kategori
			if err := rows.Scan(&k.ID, &k.Nama); err != nil {
				return c.Status(500).SendString(err.Error())
			}
			kategoris = append(kategoris, k)
		}
		return c.JSON(kategoris)
	})

	// ==================== ENDPOINT ASPIRASI (SUDAH DI-FIX JALUR GAMBARNYA) ====================
	app.Get("/aspirasi", func(c *fiber.Ctx) error {
		rows, err := db.Query("SELECT id, id_user, id_kategori, deskripsi, foto, status FROM aspirasi")
		if err != nil {
			return c.Status(500).SendString(err.Error())
		}
		defer rows.Close()

		var aspirasis []Aspirasi

		// Mengambil base URL dinamis dari Env Vercel, fallback ke localhost jika di laptop
		baseURL := os.Getenv("APP_URL")
		if baseURL == "" {
			baseURL = "http://localhost:8080"
		}

		for rows.Next() {
			var a Aspirasi
			if err := rows.Scan(&a.ID, &a.IdUser, &a.IdKategori, &a.Deskripsi, &a.Foto, &a.Status); err != nil {
				return c.Status(500).SendString(err.Error())
			}
			if a.Foto != "" {
				a.Foto = fmt.Sprintf("%s/assets/pengaduan/%s", baseURL, a.Foto)
			}
			aspirasis = append(aspirasis, a)
		}
		return c.JSON(aspirasis)
	})

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

		baseURL := os.Getenv("APP_URL")
		if baseURL == "" {
			baseURL = "http://localhost:8080"
		}

		return c.JSON(fiber.Map{
			"message": "Aspirasi berhasil ditambahkan",
			"foto":    fmt.Sprintf("%s/assets/pengaduan/%s", baseURL, filename),
		})
	})

	app.Put("/aspirasi/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")
		status := c.FormValue("status")
		_, err = db.Exec("UPDATE aspirasi SET status = ? WHERE id = ?", status, id)
		if err != nil {
			return c.Status(500).SendString(err.Error())
		}
		return c.JSON(fiber.Map{"message": "Aspirasi berhasil diupdate"})
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Fatal(app.Listen(":" + port))
}
