# API Ticketing

Backend service menggunakan **Golang** dengan arsitektur **Clean Layered Architecture**. Menggunakan **Fiber** untuk performa tinggi, **GORM** sebagai ORM, dan **PostgreSQL** untuk database.

## 🛠 Tech Stack

* **Framework:** [Gofiber/Fiber v2](https://gofiber.io/)
* **Database:** PostgreSQL
* **ORM:** [GORM](https://gorm.io/)
* **Migration:** [Golang-Migrate](https://github.com/golang-migrate/migrate)
* **Validation:** [Go-Playground/Validator v10](https://github.com/go-playground/validator)
* **Auth:** JWT (JSON Web Token)

---

## 📂 Project Structure

Proyek ini mengikuti pola **Package by Feature** untuk skalabilitas maksimal:

* `cmd/api/`: Entry point aplikasi dan tempat perakitan (Wiring) Dependency Injection.
* `internal/delivery/http/`: Layer Handler (Fiber) untuk menangani request/response.
* `internal/usecase/`: Layer Business Logic & interface kontrak usecase.
* `internal/repository/`: Layer Database access & interface kontrak repository.
* `models/`: Entity GORM yang merepresentasikan tabel database.
* `pkg/utils/`: Utility global seperti custom validator.

---

## 🚀 Getting Started

### 1. Prasyarat

* Go 1.24+
* PostgreSQL running locally
* `golang-migrate` CLI installed:
```bash
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

```

### 2. Setup Environment

Salin file `.env.example` ke `.env` dan sesuaikan kredensial database Anda:

```bash
cp .env.example .env

```

### 3. Database Migration

Gunakan Makefile untuk menjalankan migrasi database:

```bash
# Menjalankan semua migrasi (Up)
make migrate-up

# Membuat file migrasi baru (Timestamp based)
make migrate-create name=create_example_table

# Rollback migrasi terakhir (Down)
make migrate-down

```

### 4. Running the App

```bash
# Download dependencies
make deps

# Jalankan aplikasi
make run

```

---

### 4. Testing the App

```bash
# Generate mock
make gen-mock

# Buat unit test dengan _test.go

# Jalankan test keseluruhan
make test

# Jalankan test coverage
make test-cover

```

---

## 🛠 Development Workflow

Untuk menambah fitur baru (misal: "Regencies"):

1. **Migration:** Buat tabel di `database/migrations` via `make migrate-create`.
2. **Model:** Definisikan struct di `models/regencies.go`.
3. **Repository:** Buat interface & impl di `internal/repository/regencies/`.
4. **Usecase:** Buat logic & interface di `internal/usecase/regencies/`.
5. **Handler:** Buat Fiber handler di `internal/delivery/http/regencies/`.
6. **Wiring:** Daftarkan di `cmd/api/main.go`.

---

## 🔒 Authentication

Beberapa endpoint dilindungi oleh Middleware JWT. Gunakan Header berikut setelah Login:
`Authorization: Bearer <your_token>`

---

## 📜 Makefile Commands

| Command | Description |
| --- | --- |
| `make run` | Menjalankan aplikasi secara lokal |
| `make build` | Mengompilasi binary aplikasi |
| `make migrate-up` | Menjalankan migrasi SQL ke database |
| `make migrate-down` | Membatalkan migrasi SQL terakhir |
| `make deps` | Merapikan dan mengunduh dependencies |

---
