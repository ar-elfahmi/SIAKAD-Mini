<p align="center">
  <img src="sumber%20yang%20digunakan.png" alt="Sumber yang digunakan" />
</p>

> **Penggunaan AI:**
> Saya menggunakan Chat GPT untuk membantu saya dalam mengerjakan modul ini namun kode yang di generate oleh chat gpt saya TULIS ULANG MANUAL hingga tahap 4 sehingga saya terlibat aktif dalam proses mengetik, brainstorming, mempertanyakan dan troubleshooting secara langsung. Tujuan saya adalah untuk mengetahui garis besar aplikasi dijalankan, dibuat dan digunakan. sehingga ketika saya sudah mengerti apa yang saya tulis maka kemudian saya dapat mendelegasikan nya kepada AI dengan mengikuti kode/cara pikir/cara penulisan yang telah saya lakukan sebelumnya AI DIGUNAKAN UNTUK MENGURANGI REDUNDASI kode yang sudah bisa saya tulis sendiri agar menghemat waktu saya dalam menulis kode. Tujuannya Adalah untuk efisien waktu dan tidak mengorbankan pemahaman dan pembelajaran yang saya lakukan.

---

# SIAKAD-Mini

Sistem Informasi Akademik sederhana untuk mengelola data mahasiswa, mata kuliah, dan KRS (Kartu Rencana Studi).
Cocok untuk tugas kuliah, demo, atau belajar membangun REST API dengan **Go + PostgreSQL**.

---

## Daftar Isi

1. [Apa itu SIAKAD-Mini?](#apa-itu-siakad-mini)
2. [Fitur](#fitur)
3. [Yang Perlu Disiapkan](#yang-perlu-disiapkan)
4. [Cara Install](#cara-install)
5. [Cara Setup Database](#cara-setup-database)
6. [Cara Menjalankan Server](#cara-menjalankan-server)
7. [Akun Default](#akun-default)
8. [Cara Pakai API](#cara-pakai-api)
9. [Daftar Endpoint](#daftar-endpoint)
10. [Dokumentasi API (Postman)](#dokumentasi-api-postman)
11. [Struktur Project](#struktur-project)
12. [Troubleshooting](#troubleshooting)

---

## Apa itu SIAKAD-Mini?

Bayangkan sebuah aplikasi kampus yang dipakai untuk:

- **Login** sebagai admin atau mahasiswa
- **Admin** bisa melihat, menambah, mengubah, dan menghapus data mahasiswa
- **Mahasiswa** bisa melihat daftar mata kuliah, lalu **mengambil (enroll)** mata kuliah untuk semester tertentu
- Sistem otomatis mengecek **kuota** mata kuliah dan **batas SKS** berdasarkan IPK mahasiswa

Versi "Mini" ini punya data dan fitur secukupnya untuk demo atau tugas akhir. Tidak ada antarmuka web (hanya API), jadi interaksi via tools seperti Postman, HTTPie, atau `curl`.

---

## Fitur

| Role     | Yang Bisa Dilakukan                                                                |
| -------- | ---------------------------------------------------------------------------------- |
| `admin`  | Lihat semua mahasiswa, tambah/ubah/hapus mahasiswa, lihat semua course             |
| `mahasiswa` | Lihat profil sendiri, lihat daftar course, ambil (enroll) course, batal enroll   |

Detail tiap endpoint lihat bagian [Daftar Endpoint](#daftar-endpoint).

---

## Yang Perlu Disiapkan

Pastikan sudah terinstall di komputer:

1. **Go** versi 1.22 atau lebih baru
   → Download di <https://go.dev/dl/>
2. **PostgreSQL** versi 14 atau lebih baru
   → Bisa install sendiri atau pakai Laragon (yang sudah termasuk PostgreSQL)
3. **Postman / HTTPie / Bruno / Insomnia** (opsional, untuk coba API)
   → Atau cukup pakai `curl` di terminal

Cek versi:

```bash
go version
psql --version
```

---

## Cara Install

1. **Clone repository**

   ```bash
   git clone https://github.com/ar-elfahmi/SIAKAD-Mini.git
   cd SIAKAD-Mini
   ```

2. **Download dependency Go**

   ```bash
   go mod download
   ```

3. **Buat file `.env`** di root project (sudah ada contoh default, cek isinya):

   ```env
   DATABASE_URL=postgres://postgres:@localhost:5432/siakad_mini
   JWT_SECRET=siakad-mini-secret-2026
   ```

   Sesuaikan `DATABASE_URL`:
   - Ganti `postgres` (username) jika user PostgreSQL Anda berbeda
   - Tambahkan password setelah `:` jika PostgreSQL Anda pakai password, contoh:
     `postgres://postgres:password123@localhost:5432/siakad_mini`

---

## Cara Setup Database

1. **Buat database** bernama `siakad_mini` di PostgreSQL:

   ```bash
   psql -U postgres
   ```

   Di dalam prompt psql:

   ```sql
   CREATE DATABASE siakad_mini;
   \q
   ```

2. **Jalankan semua migration + seeder** dari folder `migrations/`:

   ```bash
   psql -U postgres -d siakad_mini -f migrations/001_create_courses.sql
   psql -U postgres -d siakad_mini -f migrations/002_create_users.sql
   psql -U postgres -d siakad_mini -f migrations/003_create_students.sql
   psql -U postgres -d siakad_mini -f migrations/004_create_enrollments.sql
   psql -U postgres -d siakad_mini -f migrations/005_seed_users.sql
   psql -U postgres -d siakad_mini -f migrations/006_seed_students.sql
   psql -U postgres -d siakad_mini -f migrations/007_seed_courses.sql
   psql -U postgres -d siakad_mini -f migrations/008_seed_enrollments.sql
   ```

   Atau kalau mau cepat, gunakan `migrations/reset.sql` (hati-hati: menghapus semua data dulu, lalu jalankan ulang 001–008).

---

## Cara Menjalankan Server

```bash
go run main.go
```

Kalau berhasil, akan muncul:

```
   ____   __    _  ___    _  _
  / __/  / _ \ / |/ _ \  / |/ /
 _\ \   / ___ / || \_, / / || /
/___/  /_//_//_/|_/___/_/_/|_/

SIAKAD-Mini listening on :3000
```

Server jalan di `http://localhost:3000`.

Cek cepat:

```bash
curl http://localhost:3000/
```

Harus balas:

```
yey my first hello world
```

---

## Akun Default

Dari seeder (`migrations/005_seed_users.sql`):

| Email                          | Password       | Role       |
| ------------------------------ | -------------- | ---------- |
| `admin@siakad.local`           | `Password123!` | `admin`    |
| `mahasiswa01@siakad.local`     | `Password123!` | `mahasiswa`|
| `mahasiswa02@siakad.local`     | `Password123!` | `mahasiswa`|
| ... sampai `mahasiswa20@siakad.local` | `Password123!` | `mahasiswa` |

> **Catatan keamanan:** password `Password123!` hanya untuk demo. Untuk production, ganti dengan password yang kuat dan simpan hash-nya (bukan plaintext).

---

## Cara Pakai API

### 1. Login (ambil token)

```bash
curl -X POST http://localhost:3000/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"admin@siakad.local","password":"Password123!"}'
```

Response:

```json
{
  "access_token": "eyJhbGciOi...",
  "token_type": "Bearer",
  "expires_in": 3600,
  "user": {
    "id": 1,
    "email": "admin@siakad.local",
    "role": "admin"
  }
}
```

### 2. Pakai token untuk akses endpoint lain

Simpan `access_token`, lalu kirim di header `Authorization`:

```bash
curl http://localhost:3000/api/v1/courses \
  -H "Authorization: Bearer eyJhbGciOi..."
```

### 3. Endpoint mahasiswa butuh role `admin`

```bash
# Ambil semua mahasiswa (hanya admin)
curl http://localhost:3000/api/v1/students \
  -H "Authorization: Bearer <token-admin>"

# Tambah mahasiswa baru
curl -X POST http://localhost:3000/api/v1/students \
  -H "Authorization: Bearer <token-admin>" \
  -H "Content-Type: application/json" \
  -d '{
    "nim": "187221000099",
    "nama": "Mahasiswa Baru",
    "email": "baru@siakad.local",
    "prodi": "Teknik Informatika",
    "angkatan": 2024,
    "ipk_terakhir": 3.0
  }'
```

### 4. Enroll mata kuliah (mahasiswa)

```bash
curl -X POST http://localhost:3000/api/v1/enrollments \
  -H "Authorization: Bearer <token-mahasiswa>" \
  -H "Content-Type: application/json" \
  -d '{
    "course_id": 1,
    "tahun_akademik": "2024-2025-1"
  }'
```

---

## Daftar Endpoint

Semua endpoint di bawah prefix `/api/v1` dan butuh JSON. Kecuali `POST /auth/login` dan `GET /`, semua butuh header `Authorization: Bearer <token>`.

| Method | Path                       | Auth   | Role       | Keterangan                              |
| ------ | -------------------------- | ------ | ---------- | --------------------------------------- |
| GET    | `/`                        | -      | -          | Health check                            |
| POST   | `/api/v1/auth/login`       | -      | -          | Login, balikin token                    |
| GET    | `/api/v1/auth/me`          | Bearer | any        | Info user dari token                    |
| GET    | `/api/v1/courses`          | Bearer | any        | List course (filter: semester, search, available) |
| GET    | `/api/v1/students`         | Bearer | admin      | List mahasiswa (paginated, filter, sort)|
| GET    | `/api/v1/students/:id`     | Bearer | any        | Detail 1 mahasiswa (mahasiswa hanya boleh lihat diri sendiri) |
| POST   | `/api/v1/students`         | Bearer | admin      | Tambah mahasiswa baru                   |
| PUT    | `/api/v1/students/:id`     | Bearer | admin      | Update data mahasiswa                   |
| DELETE | `/api/v1/students/:id`     | Bearer | admin      | Hapus mahasiswa (soft delete)           |
| POST   | `/api/v1/enrollments`      | Bearer | mahasiswa  | Ambil mata kuliah (enroll)              |
| DELETE | `/api/v1/enrollments/:id`  | Bearer | mahasiswa  | Batal ambil mata kuliah                 |

Detail lengkap setiap endpoint (parameter, body, response) ada di [`endpoints.json`](./endpoints.json) — bisa di-import ke Postman.

---

## Dokumentasi API (Postman)

File [`endpoints.json`](./endpoints.json) adalah **Postman Collection v2.1.0** yang berisi semua endpoint di project ini.

Cara pakai:

1. Buka Postman (atau HTTPie Desktop, Insomnia, Bruno).
2. **Import** file `endpoints.json`.
3. Collection **SIAKAD-Mini** akan muncul dengan folder:
   - **Health** — `GET /`
   - **Auth** — `POST /login`, `GET /me`
   - **Courses** — `GET /`
   - **Students** — list, filter, search, sort, pagination, detail, create, update, delete
   - **Enrollments** — create, delete
4. Edit variable `base_url` di collection kalau server jalan di host/port lain (default `http://localhost:3000`).
5. Setelah login, copy `access_token` dari response dan simpan di variable `token` (atau langsung paste di header Authorization).

> File `endpoints.json` **wajib** di-update setiap kali ada endpoint baru/berubah/dihapus. Lihat `AGENTS.md` di root project untuk aturannya.

---

## Struktur Project

```
SIAKAD-Mini/
├── main.go                       # Entry point + semua route Fiber
├── go.mod / go.sum               # Dependency Go
├── .env                          # Konfigurasi (DATABASE_URL, JWT_SECRET)
├── endpoints.json                # Postman Collection
├── README.md                     # File ini
├── AGENTS.md                     # Aturan maintenance untuk AI agent
├── migrations/                   # Skema DB + data awal
│   ├── 001_create_courses.sql
│   ├── 002_create_users.sql
│   ├── 003_create_students.sql
│   ├── 004_create_enrollments.sql
│   ├── 005_seed_users.sql
│   ├── 006_seed_students.sql
│   ├── 007_seed_courses.sql
│   ├── 008_seed_enrollments.sql
│   └── reset.sql                 # Drop semua + re-seed (dev only)
└── app/
    ├── model/                    # Struct request/response
    │   ├── auth.go
    │   ├── course.go
    │   ├── enrollment.go
    │   └── student.go
    └── repository/               # Akses DB (pgx raw SQL)
        ├── auth_repository.go
        ├── course_repository.go
        ├── enrollment_repository.go
        └── student_repository.go
```

---

## Troubleshooting

**Q: `connection refused` saat konek PostgreSQL**
- Pastikan service PostgreSQL jalan (`pg_ctl status` atau via Laragon).
- Cek `DATABASE_URL` di `.env` — username, password, port, nama DB harus sesuai.

**Q: Login selalu gagal padahal email/password benar**
- Pastikan sudah menjalankan `migrations/005_seed_users.sql` (dan hash passwordnya cocok).
- Cek apakah hash di database masih bcrypt — lihat file seed.

**Q: `401 unauthorized` padahal baru login**
- Header harus `Authorization: Bearer <token>` (dengan spasi setelah `Bearer`).
- Token expired → login ulang.

**Q: `403 akses ditolak` saat lihat detail mahasiswa**
- Mahasiswa hanya boleh lihat profil sendiri. Untuk lihat mahasiswa lain, login sebagai admin.

**Q: `409 conflict` saat enroll**
- Mata kuliah sudah pernah diambil. Batal dulu dengan `DELETE /api/v1/enrollments/:id`, lalu enroll ulang.

**Q: `422 kuota mata kuliah sudah penuh`**
- Course sudah penuh kapasitasnya. Pilih course lain.

**Q: `422 batas SKS terlampaui`**
- Total SKS yang sudah diambil + SKS course baru > batas IPK mahasiswa.
  Aturan batas: IPK ≥ 3.0 → 24 SKS, ≥ 2.5 → 21, ≥ 2.0 → 18, < 2.0 → 15.

---

## Lisensi

Project ini untuk kebutuhan tugas/pembelajaran. Bebas dipakai dan dimodifikasi.
