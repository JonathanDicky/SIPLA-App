package controllers

import (
	"fmt"
	"net/http"
	"path/filepath"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"pengaduan/config"
	"pengaduan/models"
)

func saveBuktiFile(c *gin.Context, fieldName string) string {
	file, err := c.FormFile(fieldName)
	if err != nil {
		return ""
	}
	ext      := filepath.Ext(file.Filename)
	filename := fmt.Sprintf("assets/bukti/%d_%s%s", time.Now().UnixNano(), fieldName, ext)
	if err := c.SaveUploadedFile(file, filename); err != nil {
		return ""
	}
	return filename
}

// POST /api/tanggapan — max 3 file bukti, disimpan ke kolom terpisah
func CreateTanggapan(c *gin.Context) {
	username, _ := c.Get("username")

	var petugas models.Petugas
	if err := config.DB.Where("username = ?", username).First(&petugas).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Petugas tidak ditemukan"})
		return
	}

	idPengaduanStr := c.PostForm("id_pengaduan")
	tanggapanTeks  := c.PostForm("tanggapan")
	status         := c.PostForm("status")

	if idPengaduanStr == "" || tanggapanTeks == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id_pengaduan dan tanggapan wajib diisi"})
		return
	}

	var idP uint64
	fmt.Sscanf(idPengaduanStr, "%d", &idP)

	// Upload 3 file bukti ke kolom masing-masing
	fileBukti1 := saveBuktiFile(c, "file_bukti_1")
	fileBukti2 := saveBuktiFile(c, "file_bukti_2")
	fileBukti3 := saveBuktiFile(c, "file_bukti_3")

	now := time.Now()
	t := models.Tanggapan{
		IDPengaduan:   idP,
		TglTanggapan:  now,
		TanggapanTeks: tanggapanTeks,
		FileBukti:     fileBukti1,
		FileBukti2:    fileBukti2,
		FileBukti3:    fileBukti3,
		IDPetugas:     petugas.IDPetugas,
		CreatedAt:     &now,
		UpdatedAt:     &now,
	}

	if err := config.DB.Create(&t).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan tanggapan"})
		return
	}

	if status != "" {
		config.DB.Model(&models.Pengaduan{}).
			Where("id_pengaduan = ?", idP).
			Update("status", status)
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Tanggapan berhasil dikirim", "data": t})
}

func UpdateTanggapan(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	var t models.Tanggapan
	if err := config.DB.First(&t, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Tanggapan tidak ditemukan"})
		return
	}
	var req struct {
		Tanggapan string `json:"tanggapan"`
	}
	c.ShouldBindJSON(&req)
	t.TanggapanTeks = req.Tanggapan
	now := time.Now()
	t.UpdatedAt = &now
	config.DB.Save(&t)
	c.JSON(http.StatusOK, gin.H{"message": "Tanggapan diperbarui", "data": t})
}

func DeleteTanggapan(c *gin.Context) {
	id := c.Param("id")
	var t models.Tanggapan
	if err := config.DB.First(&t, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Tanggapan tidak ditemukan"})
		return
	}
	config.DB.Delete(&t)
	c.JSON(http.StatusOK, gin.H{"message": "Tanggapan dihapus"})
}
