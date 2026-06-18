# Pengaduan Masyarakat - Full Stack App

## Tech Stack
- **Backend**: Go + Gin + GORM + JWT
- **Frontend**: React.js + Tailwind CSS + Vite
- **Database**: MySQL

## Struktur Folder
```
pengaduan/
├── backend/
│   ├── main.go
│   ├── go.mod
│   ├── .env
│   ├── config/
│   │   └── database.go
│   ├── models/
│   │   └── models.go
│   ├── middleware/
│   │   └── auth.go
│   ├── controllers/
│   │   ├── auth.go
│   │   ├── pengaduan.go
│   │   ├── tanggapan.go
│   │   └── wilayah.go
│   └── routes/
│       └── routes.go
└── frontend/
    ├── src/
    │   ├── App.jsx
    │   ├── main.jsx
    │   ├── index.css
    │   ├── context/AuthContext.jsx
    │   ├── services/api.js
    │   ├── components/
    │   │   ├── AdminLayout.jsx
    │   │   └── UserLayout.jsx
    │   └── pages/
    │       ├── Login.jsx
    │       ├── Register.jsx
    │       ├── admin/
    │       │   ├── Dashboard.jsx
    │       │   ├── Pengaduan.jsx
    │       │   └── Masyarakat.jsx
    │       └── user/
    │           └── Pengaduan.jsx
    ├── package.json
    ├── vite.config.js
    ├── tailwind.config.js
    ├── postcss.config.js
    └── index.html
```

---

## Cara Menjalankan

### Prasyarat
- Go 1.21+
- Node.js 18+
- MySQL / XAMPP

---

### 1. Setup Database
Buka phpMyAdmin atau MySQL, import file SQL:
```sql
CREATE DATABASE pengaduan_masyarakat3;
-- lalu import file pengaduan_masyarakat3.sql
```

---

### 2. Jalankan Backend (Go)

```bash
cd backend

# Install dependencies
go mod tidy

# Buat folder untuk upload foto
mkdir -p assets/pengaduan

# Edit .env sesuaikan password MySQL Anda
# DB_PASSWORD=your_password

# Jalankan server
go run main.go
```
Server berjalan di: http://localhost:8080

---

### 3. Jalankan Frontend (React)

```bash
cd frontend

# Install dependencies
npm install

# Jalankan development server
npm run dev
```
Aplikasi berjalan di: http://localhost:5173

---

## Akun Default (dari database)

| Role | Username | Password |
|------|----------|----------|
| Admin | admin2 | petugas |
| Petugas | petugas | petugas |
| Masyarakat | user | user1234 |

> Catatan: Password di database sudah di-hash dengan bcrypt.
> Jika login gagal, reset password via SQL:
> ```sql
> UPDATE petugas SET password = '$2y$10$/iGk2Q1bwRNqPnfQEJDt0.IcbEN7Kth6391V6bF73l7mPYBav.huG' WHERE username = 'admin';
> ```

---

## API Endpoints

### Auth
| Method | Endpoint | Deskripsi |
|--------|----------|-----------|
| POST | /api/auth/login | Login |
| POST | /api/auth/register | Registrasi masyarakat |
| GET | /api/auth/me | Info user login |

### Pengaduan
| Method | Endpoint | Deskripsi |
|--------|----------|-----------|
| GET | /api/pengaduan | List pengaduan |
| GET | /api/pengaduan/:id | Detail pengaduan |
| POST | /api/pengaduan | Buat pengaduan (multipart/form-data) |
| PUT | /api/pengaduan/:id | Update pengaduan |
| DELETE | /api/pengaduan/:id | Hapus pengaduan |

### Admin Only
| Method | Endpoint | Deskripsi |
|--------|----------|-----------|
| GET | /api/dashboard | Statistik dashboard |
| POST | /api/tanggapan | Kirim tanggapan |
| GET | /api/masyarakat | List masyarakat |

### Wilayah (Public)
| Method | Endpoint | Deskripsi |
|--------|----------|-----------|
| GET | /api/provinces | Daftar provinsi |
| GET | /api/regencies?province_id= | Daftar kabupaten |
| GET | /api/districts?regency_id= | Daftar kecamatan |
| GET | /api/villages?district_id= | Daftar desa |
