<div align="center">

# 🛡️ JevBeez

**Automated Evidence-Grounded Fact-Checking & Misinformation Detection Engine**

*Melawan disinformasi dan hoax digital menggunakan penalaran probabilistik terkalibrasi System One (TypeSafe Jev) dan penelusuran bukti empiris berbasis web.*

[![Go Version](https://img.shields.io/badge/Go-1.23-00ADD8?style=flat-square&logo=go)](https://golang.org)
[![Nuxt](https://img.shields.io/badge/Nuxt-4.x-00DC82?style=flat-square&logo=nuxtdotjs)](https://nuxt.com)
[![Vue 3](https://img.shields.io/badge/Vue-3.5-4FC08D?style=flat-square&logo=vuedotjs)](https://vuejs.org)
[![Tailwind CSS](https://img.shields.io/badge/Tailwind_CSS-3.4-38B2AC?style=flat-square&logo=tailwind-css)](https://tailwindcss.com)
[![Redis](https://img.shields.io/badge/Redis-7-DC382D?style=flat-square&logo=redis)](https://redis.io)
[![Docker](https://img.shields.io/badge/Docker-Enabled-2496ED?style=flat-square&logo=docker)](https://www.docker.com)
[![License](https://img.shields.io/badge/License-MIT-blue.flat-square)](LICENSE)

</div>

---

## 📸 Antarmuka Aplikasi (User Interface)

### 1. Halaman Utama (Fact-Check Workbench)
Desain editorial modern dengan tipografi berwibawa, bebas dari distraksi visual generik (*anti-slop*), dan berfokus pada kecepatan input:

![JevBeez Workbench](docs/screenshots/hero.png)

### 2. Dossier Hasil Analisis & Verifikasi Bukti
Menampilkan putusan keabsahan fakta (*Verdict*), skor keyakinan (*Confidence Score*), rincian metrik probabilitas terkalibrasi (*Calibrated Probabilities*), serta rujukan dokumen primer yang ditemukan di internet:

![JevBeez Verdict Dossier](docs/screenshots/verdict-analysis.png)

---

## 🌟 Sekilas Tentang JevBeez

Penyebaran disinformasi dan hoax di grup percakapan keluarga maupun media sosial sering kali sulit dibendung karena lambatnya proses verifikasi manual. Di sisi lain, model bahasa besar (LLM) konvensional kerap mengalami **halusinasi** (*hallucination*) dan menyajikan data yang terdengar meyakinkan meski tanpa dasar fakta empiris.

**JevBeez** hadir sebagai solusi pemeriksa fakta otomatis yang dibangun dengan arsitektur hibrida:
1. **Web Scraper & Social Media Content Extractor**: Mengambil konten klaim dari artikel berita maupun postingan media sosial (Twitter/X, TikTok, YouTube, Instagram, Threads, Facebook).
2. **Multi-Source Evidence Grounding**: Melakukan pencarian bukti tervalidasi dari situs berita arus utama dan portal klarifikasi fakta terpercaya secara langsung (*live search* via Tavily & Firecrawl).
3. **System One Decision Engine (TypeSafe Jev)**: Menganalisis korelasi klaim dengan bukti melalui pemodelan probabilitas terkalibrasi matematis, bukan sekadar respons teks generatif bebas.

---

## ✨ Fitur Utama

- 🔍 **Ekstraksi Konten Multi-Platform**: Mendukung ekstraksi teks dari berbagai portal berita web dan tautan media sosial melalui integrasi scraping aman dan protokol OEMBED.
- ⚖️ **Probabilistic Decision Modeling (TypeSafe Jev)**: Menghasilkan 3 metrik kuantitatif independen:
  - `supported_by_evidence`: Derajat probabilitas bahwa klaim didukung oleh fakta web.
  - `contradicted_by_evidence`: Derajat probabilitas bahwa klaim dibantah oleh bukti kredibel.
  - `sensational_framing`: Derajat indikasi bahasa klikbait, provokatif, atau judul bombastis.
- 🎯 **Tiga Klasifikasi Putusan Faktual**:
  - `VALID` *(Terbukti Fakta Faktual)*: Didukung secara konsisten oleh sumber terpercaya.
  - `HOAX` *(Terindikasi Hoax / Bantahan)*: Klaim terbantahkan oleh laporan valid atau terkonfirmasi palsu.
  - `MISLEADING` *(Menyesatkan / Belum Terverifikasi)*: Klaim ambigu, keluar konteks, atau bukti belum memadai.
- 📱 **Anti-Slop Editorial Design**: Antarmuka bersih, tenang, bertipografi tegas, dan ramah pengguna lintas generasi tanpa elemen visual berlebih.
- 🚀 **1-Click WhatsApp Debunk Sharing**: Tombol bagikan langsung ke WhatsApp untuk memudahkan pengguna mengklarifikasi hoax di grup pesan instan hanya dalam satu ketukan.
- ⚡ **Ultra-Fast Redis Caching**: Hasil evaluasi disimpan dalam cache dengan TTL optimal untuk memberikan respon instan bagi URL yang sering diperiksa.
- 🛡️ **Enterprise Security & SSRF Protection**: Dilengkapi dengan pencegahan SSRF (blokir private/loopback IP), sanitasi URL berlapis, rate limiting, dan Content Security Policy (CSP) ketat.

---

## 🏛️ Arsitektur Sistem

```
                                  +-----------------------+
                                  |   Pengguna / Klien    |
                                  |   (Browser / Mobile)  |
                                  +-----------+-----------+
                                              |
                                              | HTTPS
                                              v
                              +-------------------------------+
                              |    Nuxt 4 Frontend (SSR)      |
                              |   (Vue 3, Tailwind, Vite)     |
                              +---------------+---------------+
                                              |
                                              | Internal Network / Proxy
                                              v
                              +-------------------------------+
                              |      Go Fiber API Service     |
                              |   (Core Engine, Port 8080)    |
                              +---------------+---------------+
                                  |           |           |
            +---------------------+           |           +---------------------+
            | Redis Cache                     |                                 |
            v (Key-Value TTL)                 v (Live Grounding)                v (Probabilistic Evaluation)
   +-----------------+            +------------------------+            +-----------------------+
   |   Redis Store   |            | Search & Scraping Core |            |   TypeSafe Jev AI     |
   |   (Port 6379)   |            |  - Tavily Search API   |            | System One Evaluator  |
   +-----------------+            |  - Firecrawl Engine    |            +-----------------------+
                                  |  - SSRF Safe Scraper   |
                                  +------------------------+
```

---

## 🛠️ Tech Stack

| Layer | Komponen / Library | Kegunaan |
|---|---|---|
| **Frontend** | Nuxt 4 (Vue 3, TypeScript) | Server-Side Rendering (SSR), performa tinggi & SEO-ready |
| **Styling** | Tailwind CSS v3 | Desain sistem bertipe editorial, modular, dan responsif |
| **Backend Core** | Go 1.23, Fiber v2 | RESTful API mikrografis, routing berkecepatan tinggi & hemat memori |
| **Decision Model** | TypeSafe Jev AI | Evaluasi keputusan probabilistik non-halusinasi (*System One*) |
| **Search & Evidence** | Tavily Search & Firecrawl | Mesin pencarian bukti rujukan web dan ekstraksi dokumen primer |
| **Caching Layer** | Redis 7 Alpine | Penyimpanan hasil fact-check berbasis key-value dengan proteksi sandi |
| **Containerization** | Docker & Docker Compose | Standardisasi deployment multi-container untuk kemudahan operasional |

---

## 🚀 Panduan Memulai (Getting Started)

### Prasyarat
- [Docker](https://www.docker.com/) & Docker Compose (v2.x atau lebih baru)
- Atau jika ingin menjalankan tanpa Docker:
  - [Node.js](https://nodejs.org/) (>= 22.x) & [Bun](https://bun.sh/)
  - [Go](https://golang.org/) (>= 1.23)
  - [Redis Server](https://redis.io/) (>= 7.x)

### 1. Kloning Repositori
```bash
git clone https://github.com/Scr34m1nG/JevBeez.git
cd JevBeez
```

### 2. Konfigurasi Environment Variable
Salin berkas template konfigurasi dan masukkan API Key yang sesuai:
```bash
cp .env.example .env
```

Isi variabel pada berkas `.env`:
```env
PORT=8080
CORS_ORIGIN=*
REDIS_ADDR=redis:6379
REDIS_PASSWORD=

# Kunci API Layanan Eksternal
TYPESAFE_API_KEY=your_typesafe_api_key_here
TAVILY_API_KEY=your_tavily_api_key_here
FIRECRAWL_API_KEY=your_firecrawl_api_key_here

API_INTERNAL_URL=http://api:8080
```

> **Catatan Keamanan:** Jangan pernah melakukan commit berkas `.env` asli atau mempublikasikan kunci API pribadi ke git repository publik.

---

### 3. Menjalankan Menggunakan Docker Compose (Direkomendasikan)

Cukup jalankan satu perintah berikut:
```bash
docker compose up --build -d
```

Layanan akan aktif pada port:
- **Web Interface:** `http://localhost:3000`
- **Backend API:** `http://localhost:8080`
- **Redis Cache:** `127.0.0.1:6379`

Untuk menghentikan layanan:
```bash
docker compose down
```

---

### 4. Menjalankan Manual untuk Pengembangan Lokal (Local Dev)

**Langkah A: Menjalankan Redis**
```bash
docker run -d --name jevbeez-redis -p 6379:6379 redis:7-alpine
```

**Langkah B: Menjalankan Backend API (Go)**
```bash
cd apps/api
go run cmd/server/main.go
```

**Langkah C: Menjalankan Frontend (Nuxt)**
```bash
cd apps/web
bun install
bun run dev
```
Buka browser pada tautan `http://localhost:3000`.

---

## 📡 Dokumentasi Endpoint API

### 1. `GET /health`
Memeriksa status ketersediaan layanan API.

**Response:**
```json
{
  "service": "jevbeez-api",
  "status": "ok"
}
```

---

### 2. `POST /api/v1/analyze`
Menganalisis dan memverifikasi keabsahan fakta dari sebuah tautan artikel atau media sosial.

**Request Headers:**
- `Content-Type: application/json`

**Request Body:**
```json
{
  "url": "https://example.com/berita-tertentu"
}
```

**Contoh Response Sukses (`200 OK`):**
```json
{
  "verdict": "VALID",
  "confidence": 88,
  "judgments": {
    "supported_by_evidence": 0.88,
    "contradicted_by_evidence": 0.05,
    "sensational_framing": 0.08
  },
  "article_title": "Judul Artikel Terverifikasi",
  "article_url": "https://example.com/berita-tertentu",
  "source_domain": "example.com",
  "summary": "Informasi dalam artikel ini terverifikasi dan didukung oleh bukti laporan kredibel.",
  "red_flags": [],
  "evidences": [
    {
      "source": "antaranews.com",
      "title": "Klarifikasi dan Laporan Terkait Fakta",
      "snippet": "Kutipan laporan verifikasi pendukung...",
      "url": "https://www.antaranews.com/berita/..."
    }
  ],
  "is_cached": false,
  "is_social": false,
  "platform": ""
}
```

**Status Kode HTTP:**
- `200 OK`: Evaluasi berhasil dilakukan.
- `400 Bad Request`: Format JSON tidak valid atau URL tidak memenuhi kriteria keamanan.
- `422 Unprocessable Entity`: Konten pada URL tidak dapat diakses atau teks artikel terlalu singkat.
- `429 Too Many Requests`: Melebihi batas kuota permintaan (rate limit).
- `500 Internal Server Error`: Kendala internal pada server.

---

## 🔒 Keamanan & Kebijakan Privasi (Security & Privacy)

Sistem JevBeez dirancang dengan mengedepankan prinsip **Zero-Risk** dan **Privacy-by-Design**:
1. **Mitigasi SSRF (Server-Side Request Forgery)**:
   - Validasi ketat terhadap skema protokol (hanya `http` dan `https` yang diizinkan).
   - Resolusi DNS sebelum *fetch* untuk mencegah pemanggilan alamat IP privat (RFC 1918), loopback (`127.0.0.0/8`), link-local (`169.254.0.0/16`), atau IP khusus internal.
2. **Perlindungan Data Pribadi**:
   - Sistem **tidak menyimpan rekaman identitas pribadi** pengguna.
   - Analisis berfokus murni pada teks berita dan bukti publik di internet.
3. **Hardened HTTP Headers**:
   - Header keamanan lengkap: `Content-Security-Policy`, `Strict-Transport-Security` (HSTS), `X-Content-Type-Options: nosniff`, `X-Frame-Options: DENY`, dan `Referrer-Policy: strict-origin-when-cross-origin`.
4. **Zero Secret Exposure**:
   - Seluruh kredensial dan kunci API terisolasi secara eksklusif melalui *environment variable* aman dan tidak pernah masuk ke dalam bundle frontend klien.

---

## 🧪 Menjalankan Pengujian Otomatis

**Menjalankan Test Suite Backend (Go):**
```bash
cd apps/api
go test -v ./...
```

**Memverifikasi Kompilasi Frontend (Nuxt):**
```bash
cd apps/web
bun run build
```

---

## 🤝 Kontribusi & GitFlow

Proyek ini menerapkan standar **GitFlow** dan **Conventional Commits**:
- `main`: Branch stabil khusus rilis produksi.
- `develop`: Branch integrasi untuk siklus pengembangan harian.
- `feature/<nama>`: Fitur baru dibuat dari `develop` dan di-merge kembali setelah pengujian selesai.
- `fix/<nama>`: Perbaikan bug berbasis branch `develop`.

Setiap kontribusi pull request sangat terbuka! Pastikan seluruh pengujian unit (`go test` dan `bun run build`) berhasil lolos sebelum mengajukan merge.

---

## 📄 Lisensi

Didistribusikan di bawah lisensi **MIT License**. Lihat berkas [LICENSE](LICENSE) untuk informasi lebih lanjut.

---

<div align="center">
  <sub>Dibangun untuk ekosistem internet yang lebih sehat, bersih dari disinformasi, dan berbasis bukti empiris.</sub>
</div>
