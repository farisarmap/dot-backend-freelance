# DOT Backend Freelance Online Test

**Project Name** adalah aplikasi backend yang dibangun menggunakan Go, Echo, dan GORM. Aplikasi ini menyediakan API untuk mengelola pengguna dan pesanan dengan fitur validasi, caching, dan pengujian end-to-end (E2E) yang komprehensif.

## Daftar Isi

- [DOT Backend Freelance Online Test](#dot-backend-freelance-online-test)
  - [Daftar Isi](#daftar-isi)
  - [Arsitektur dan Pola Desain](#arsitektur-dan-pola-desain)
    - [Layered Architecture](#layered-architecture)
    - [Service Layer](#service-layer)
    - [Handler Layer](#handler-layer)
    - [Adapters](#adapters)
    - [Validasi](#validasi)
    - [Caching](#caching)
    - [Strategi Pengujian](#strategi-pengujian)
      - [End-to-End (E2E) Testing](#end-to-end-e2e-testing)
  - [Struktur Proyek](#struktur-proyek)
  - [API Endpoints](#api-endpoints)
    - [Deploy App](#deploy-app)
      - [Konfigurasi Environment Variables](#konfigurasi-environment-variables)
      - [Menjalankan Layanan dengan Docker Compose](#menjalankan-layanan-dengan-docker-compose)
  - [Menjalankan Aplikasi di Local](#menjalankan-aplikasi-di-local)
    - [Langkah-Langkah](#langkah-langkah)
      - [1. Menjalankan Layanan Pendukung dengan Docker](#1-menjalankan-layanan-pendukung-dengan-docker)
      - [Jalankan Layanan](#jalankan-layanan)
      - [2. Menjalankan Migrasi Database](#2-menjalankan-migrasi-database)
      - [3. Menjalankan Aplikasi](#3-menjalankan-aplikasi)
      - [4. Test End to End](#4-test-end-to-end)

## Arsitektur dan Pola Desain

Proyek ini diorganisir dengan mempertimbangkan prinsip-prinsip **Layered Architecture** dan berbagai pola desain yang memastikan **skalabilitas**, **maintainabilitas**, dan **testabilitas**. Berikut adalah alasan mengapa Saya memilih pola-pola ini:

### Layered Architecture

Saya menggunakan **Layered Architecture** untuk memisahkan tanggung jawab ke dalam lapisan-lapisan yang berbeda. Ini membuat kode lebih terstruktur dan memudahkan pengembangan serta pemeliharaan.

- **Entities**: Menyimpan struktur data dan logika bisnis inti.
- **Services**: Mengatur logika bisnis dan berinteraksi dengan repositories.
- **Handlers**: Menangani permintaan HTTP dan mengembalikan respons.
- **Adapters**: Menghubungkan aplikasi dengan sistem eksternal seperti caching dan database.

### Service Layer

Lapisan **Service** bertanggung jawab untuk mengelola logika bisnis dan mengoordinasikan operasi antara berbagai repositories. Ini memastikan bahwa logika bisnis terpusat dan mudah diatur.

- **Pengelolaan Logika Bisnis**: Menangani aturan-aturan bisnis yang kompleks.
- **Interaksi Repositories**: Mengatur komunikasi antara berbagai repositories.
- **Reusabilitas**: Fungsi-fungsi bisnis dapat digunakan kembali di berbagai handler.

### Handler Layer

Lapisan **Handler** bertanggung jawab untuk menangani permintaan HTTP, memproses data yang diterima, dan mengembalikan respons yang sesuai kepada klien.

- **Binding dan Validasi**: Mengikat data dari permintaan ke struct dan memvalidasinya.
- **Pengelolaan Respons**: Mengembalikan data atau pesan error yang sesuai.
- **Error Handling**: Menangani error dan memastikan respons yang konsisten.

### Adapters

Saya menggunakan **Adapters** untuk menghubungkan aplikasi dengan sistem eksternal seperti caching dan database. Ini memungkinkan penggantian atau modifikasi sistem eksternal tanpa mempengaruhi lapisan bisnis.

- **CacheManager**: Mengelola operasi caching menggunakan Redis atau sistem caching lain.
- **Database Repositories**: Mengelola operasi database menggunakan GORM.

### Validasi

Saya menggunakan **go-playground/validator** untuk memastikan bahwa input yang diterima oleh API memenuhi kriteria yang ditetapkan.

### Caching

Implementasi **caching** membantu meningkatkan performa aplikasi dengan menyimpan data yang sering diakses dalam memori.

- **CacheManager**:
  - Mengelola operasi **GET**, **SET**, dan **DELETE** untuk data yang di-cache.
  - Mengurangi beban pada database dengan menyimpan hasil query yang sering digunakan.
  - Meningkatkan waktu respons API dengan mengakses data dari cache yang lebih cepat.

### Strategi Pengujian

Saya menerapkan strategi pengujian yang komprehensif untuk memastikan kualitas dan stabilitas aplikasi.

#### End-to-End (E2E) Testing

E2E testing dilakukan untuk memastikan bahwa seluruh alur aplikasi bekerja dengan baik dari awal hingga akhir.

## Struktur Proyek

Berikut adalah struktur direktori utama dalam proyek ini:

```graphql
.
├── cmd
│   └── main.go                # Entry point aplikasi
├── internal
│   ├── adapter
│   │   ├── redis.go   # Implementasi CacheManager
│   │   ├── order.go# Repository untuk Order
│   │   └── user.go # Repository untuk User
│   ├── api
│   │   ├── user.go           # Struct untuk User request
│   │   ├── order.go          # Struct untuk Order request  
│   ├── entity
│   │   ├── order.go            # Entity Order
│   │   └── user.go             # Entity User
│   ├── handler
│   │   ├── order.go    # Handler untuk Order endpoints
│   │   ├── user.go     # Handler untuk User endpoints
│   └── service
│       ├── order.go    # Service untuk Order
│       └── user.go     # Service untuk User
├── pkg
│   └── response.go             # Struct dan fungsi untuk response API
└── test
    ├── user_e2e_test.go            # E2E Test untuk User
    └── order_e2e_test.go           # E2E Test untuk Order

```

## API Endpoints
Proyek ini menyediakan berbagai endpoint API untuk mengelola pengguna dan pesanan. Berikut adalah dokumentasi lengkap mengenai endpoint yang tersedia:

- User Endpoints
- Membuat User
- Endpoint: /users

- Method: POST

- Deskripsi: Membuat pengguna baru.

- Request Body:

```json
{
    "name": "John Doe",
    "email": "john.doe@example.com"
}
```

- Respons:

- 201 Created

```json
{
    "status": "success",
    "message": "User created",
    "data": {
        "id": 1,
        "name": "John Doe",
        "email": "john.doe@example.com",
        "created_at": "2024-12-23T12:00:00Z",
        "updated_at": "2024-12-23T12:00:00Z"
    }
}
```

- 400 Bad Request

```json
{
    "status": "error",
    "message": "Field 'Name' failed on the 'required' tag",
    "data": ""
}
```

- Menghapus User
- Endpoint: /users/{id}
- Method: DELETE
- Deskripsi: Menghapus pengguna berdasarkan ID.
- Path Parameters:
- id (integer): ID pengguna yang akan dihapus.
- Respons:
- 200 OK

```json
{
    "status": "success",
    "message": "User deleted",
    "data": null
}
```

- 404 Not Found

```json
{
    "status": "error",
    "message": "User not found",
    "data": ""
}
```

- Order Endpoints
- Membuat Order
- Endpoint: /orders

- Method: POST

- Deskripsi: Membuat order baru.

- Request Body:

```json
{
    "order_name": "Order ABC",
    "user_id": 1
}
```

- Respons:

- 201 Created

```json
{
    "status": "success",
    "message": "Order created",
    "data": {
        "id": 1,
        "order_name": "Order ABC",
        "user_id": 1,
        "created_at": "2024-12-23T12:00:00Z",
        "updated_at": "2024-12-23T12:00:00Z"
    }
}
```

- 400 Bad Request

```json
{
    "status": "error",
    "message": "Field 'OrderName' failed on the 'required' tag",
    "data": ""
}
```

- Mengambil Semua Orders
- Endpoint: /orders
- Method: GET
- Deskripsi: Mengambil semua order yang tersedia.
- Respons:
- 200 OK

```json
{
    "status": "success",
    "message": "Orders retrieved",
    "data": [
        {
            "id": 1,
            "order_name": "Order ABC",
            "user_id": 1,
            "created_at": "2024-12-23T12:00:00Z",
            "updated_at": "2024-12-23T12:00:00Z"
        },
        {
            "id": 2,
            "order_name": "Order XYZ",
            "user_id": 2,
            "created_at": "2024-12-24T08:30:00Z",
            "updated_at": "2024-12-24T08:30:00Z"
        }
    ]
}
```

- Mengambil Order Berdasarkan ID
- Endpoint: /orders/{id}
- Method: GET
- Deskripsi: Mengambil order berdasarkan ID.
- Path Parameters:
- id (integer): ID order yang akan diambil.
- Respons:

- 200 OK

```json
{
    "status": "success",
    "message": "Order retrieved",
    "data": {
        "id": 1,
        "order_name": "Order ABC",
        "user_id": 1,
        "created_at": "2024-12-23T12:00:00Z",
        "updated_at": "2024-12-23T12:00:00Z"
    }
}
```

- 404 Not Found

```json
{
    "status": "error",
    "message": "Order not found",
    "data": ""
}
```

- Memperbarui Order
- Endpoint: /orders/{id}

- Method: PUT

- Deskripsi: Memperbarui order secara keseluruhan.

- Path Parameters:

- id (integer): ID order yang akan diperbarui.
- Request Body:

```json
{
    "order_name": "Order DEF",
    "user_id": 1
}
```

- Respons:

- 200 OK

```json
{
    "status": "success",
    "message": "Order updated",
    "data": {
        "id": 1,
        "order_name": "Order DEF",
        "user_id": 1,
        "created_at": "2024-12-23T12:00:00Z",
        "updated_at": "2024-12-23T14:00:00Z"
    }
}
```

- 400 Bad Request

```json
{
    "status": "error",
    "message": "Field 'OrderName' failed on the 'min' tag",
    "data": ""
}
```

- 404 Not Found

```json
{
    "status": "error",
    "message": "Order not found",
    "data": ""
}
```

- Memperbarui Order Secara Parsial
- Endpoint: /orders/{id}

- Method: PATCH

- Deskripsi: Memperbarui sebagian data order.

- Path Parameters:

- id (integer): ID order yang akan diperbarui.
- Request Body:

```json
{
    "order_name": "Order GHI"
}
```

- Respons:

- 200 OK

```json
{
    "status": "success",
    "message": "Order updated",
    "data": {
        "id": 1,
        "order_name": "Order GHI",
        "user_id": 1,
        "created_at": "2024-12-23T12:00:00Z",
        "updated_at": "2024-12-23T15:00:00Z"
    }
}
```

- 400 Bad Request

```json
{
    "status": "error",
    "message": "Field 'OrderName' failed on the 'min' tag",
    "data": ""
}
```

- 404 Not Found

```json
{
    "status": "error",
    "message": "Order not found",
    "data": ""
}
```

- Menghapus Order
- Endpoint: /orders/{id}
- Method: DELETE
- Deskripsi: Menghapus order berdasarkan ID.
- Path Parameters:
- id (integer): ID order yang akan dihapus.
- Respons:
- 200 OK

```json
{
    "status": "success",
    "message": "Order deleted",
    "data": null
}
```

- 404 Not Found

```json
{
    "status": "error",
    "message": "Order not found",
    "data": ""
}
```

### Deploy App

#### Konfigurasi Environment Variables

1. **Buat file `config.json`** berdasarkan contoh `config.json.example` jika ada.
2. **Sesuaikan variabel lingkungan** sesuai kebutuhan.

#### Menjalankan Layanan dengan Docker Compose

```bash
docker compose -f "docker-compose.yml" up -d --build 
```

Mengakses Aplikasi
Aplikasi akan berjalan di http://localhost:8080 atau sesuai dengan konfigurasi yang Anda tetapkan dalam docker-compose.yml.

Menghentikan Layanan

```bash
docker-compose down
```

Perintah ini akan menghentikan dan menghapus kontainer yang dijalankan oleh Docker Compose.

## Menjalankan Aplikasi di Local

Panduan ini menjelaskan langkah-langkah untuk menjalankan aplikasi secara lokal menggunakan Docker untuk layanan pendukung, serta menjalankan migrasi database dan aplikasi utama.

### Langkah-Langkah

#### 1. Menjalankan Layanan Pendukung dengan Docker

Aplikasi ini membutuhkan layanan **PostgreSQL** dan **Redis**. Pastikan Anda telah menginstal **Docker** di sistem Anda.

#### Jalankan Layanan

Gunakan perintah berikut untuk menjalankan PostgreSQL dan Redis:

```bash
docker compose  -f "docker-compose.yml" up -d --build postgres redis 
```

postgres: Layanan database PostgreSQL.
redis: Layanan cache Redis.

#### 2. Menjalankan Migrasi Database

Setelah layanan database berjalan, lakukan migrasi database dengan perintah berikut:

```bash
make migrate-up
```

make migrate-up akan menjalankan migrasi schema database menggunakan file migrasi yang telah disiapkan.


#### 3. Menjalankan Aplikasi

Jalankan aplikasi utama menggunakan perintah berikut:

```bash
make run
```

#### 4. Test End to End

```bash
go test ./test/
```
