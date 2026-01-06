# 🚓 Chatbot Assistant - Testing Documentation

Dokumentasi lengkap untuk testing AI Chatbot dengan persona "Sobat Lantas" dan fitur E-Tilang.

---

## 📋 Daftar Isi

1. [Testing Persona Sobat Lantas](#-testing-persona-sobat-lantas)
2. [Testing E-Tilang Feature](#-testing-e-tilang-feature)
3. [Data Dummy E-Tilang](#-data-dummy-e-tilang)
4. [Automated Test Script](#-automated-test-script)
5. [Response Format](#-response-format)
6. [API Endpoints](#-api-endpoints)
7. [Troubleshooting](#-troubleshooting)

---

## 🎭 Testing Persona Sobat Lantas

### 1. Chat Pertama (Harus ada sapaan "Halo Sobat Lantas!")

```bash
curl -X POST http://localhost:8080/api/v1/chat \
  -H "Content-Type: application/json" \
  -d '{
    "message": "Boleh gak bonceng dua anak kecil di motor?",
    "context": {
      "location": "Jl. Sudirman, Jakarta",
      "speed": 35.5,
      "traffic": "lancar",
      "latitude": -6.2088,
      "longitude": 106.8456
    }
  }' | jq '.'
```

**Expected Response:**
- ✅ Harus dimulai dengan: **"Halo Sobat Lantas!"**
- ✅ Isi respons ramah, santai, dan peduli keselamatan
- ✅ Return `session_id` untuk chat selanjutnya
- ✅ Menggunakan kata-kata seperti "yaa", "loh", "nih" dengan natural

**Contoh:**
```
"Halo Sobat Lantas! Demi keselamatan, sebaiknya jangan bonceng dua anak kecil yaa. 
Bahaya banget loh. Anak-anak harus pakai helm SNI dan cukup satu saja yang dibonceng. 
Utamakan keselamatan keluarga kita!"
```

---

### 2. Chat Kedua (TIDAK boleh ada "Halo Sobat Lantas!")

**Ganti `YOUR_SESSION_ID` dengan session_id dari response pertama**

```bash
curl -X POST http://localhost:8080/api/v1/chat \
  -H "Content-Type: application/json" \
  -d '{
    "message": "Kalau satu anak kecil gimana?",
    "context": {
      "location": "Jl. Sudirman, Jakarta",
      "speed": 35.5,
      "traffic": "lancar",
      "latitude": -6.2088,
      "longitude": 106.8456
    },
    "session_id": "YOUR_SESSION_ID"
  }' | jq '.'
```

**Expected Response:**
- ❌ **TIDAK** ada sapaan "Halo Sobat Lantas!"
- ✅ Langsung jawab pertanyaan dengan ramah
- ✅ Tetap menggunakan persona yang peduli keselamatan

---

### 3. Chat Lanjutan dengan Konteks Berbeda

```bash
curl -X POST http://localhost:8080/api/v1/chat \
  -H "Content-Type: application/json" \
  -d '{
    "message": "Berapa batas kecepatan di tol?",
    "context": {
      "location": "Tol Jagorawi",
      "speed": 85.0,
      "traffic": "lancar",
      "latitude": -6.2345,
      "longitude": 106.8765
    },
    "session_id": "YOUR_SESSION_ID"
  }' | jq '.'
```

---

### 4. Session Baru (Harus ada sapaan lagi)

```bash
curl -X POST http://localhost:8080/api/v1/chat \
  -H "Content-Type: application/json" \
  -d '{
    "message": "Apa aturan pakai helm untuk anak?",
    "context": {
      "location": "Tangerang Selatan",
      "speed": 40.0,
      "traffic": "padat",
      "latitude": -6.3024,
      "longitude": 106.6519
    }
  }' | jq '.'
```

**Expected Response:**
- ✅ Harus ada sapaan **"Halo Sobat Lantas!"** lagi (karena session baru)

---

## 🎫 Testing E-Tilang Feature

### Cek E-Tilang dengan Pelanggaran (B1234SV)

```bash
curl -X POST http://localhost:8080/api/v1/chat \
  -H "Content-Type: application/json" \
  -d '{
    "message": "Tolong bantu cek tilang dengan no polisi B1234SV",
    "context": {
      "location": "Jakarta",
      "speed": 30.0,
      "traffic": "lancar",
      "latitude": -6.2088,
      "longitude": 106.8456
    }
  }' | jq '.'
```

**Expected Response:**
- ✅ AI menjawab dengan persona Sobat Lantas
- ✅ Menampilkan informasi kendaraan (nomor polisi, rangka, pemilik, jenis)
- ✅ Detail pelanggaran yang ada
- ✅ Total denda yang harus dibayar
- ✅ Saran untuk segera melunasi

**Contoh Response:**
```
"Halo Sobat Lantas! Untuk kendaraan dengan nomor polisi B 1234 SV atas nama Budi Santoso, 
ada 2 pelanggaran yang tercatat nih:

1. Melanggar lampu merah di Jl. Sudirman - Jakarta Pusat (15 Des 2025) - Rp 500.000
2. Tidak menggunakan helm SNI di Jl. Gatot Subroto (20 Des 2025) - Rp 250.000

Total denda: Rp 750.000

Sebaiknya segera dilunasi yaa biar gak kena denda tambahan. Jaga keselamatan berkendara!"
```

---

### Cek E-Tilang Bersih (B9999ZZ)

```bash
curl -X POST http://localhost:8080/api/v1/chat \
  -H "Content-Type: application/json" \
  -d '{
    "message": "Cek tilang motor saya nomor polisi B 9999 ZZ",
    "context": {
      "location": "Bandung",
      "speed": 40.0,
      "traffic": "lancar",
      "latitude": -6.9175,
      "longitude": 107.6191
    }
  }' | jq '.'
```

**Expected Response:**
- ✅ AI memberitahu tidak ada tilang tercatat
- ✅ Memberi apresiasi karena kendaraan bersih

**Contoh Response:**
```
"Alhamdulillah, untuk kendaraan dengan nomor polisi B 9999 ZZ tidak ada tilang yang 
tercatat. Keren! Tetap patuhi peraturan lalu lintas yaa dan jaga keselamatan berkendara!"
```

---

### Cek E-Tilang dengan Pelanggaran Parkir (B5678XY)

```bash
curl -X POST http://localhost:8080/api/v1/chat \
  -H "Content-Type: application/json" \
  -d '{
    "message": "Cek denda tilang mobil B5678XY dong",
    "context": {
      "location": "Jakarta",
      "speed": 0.0,
      "traffic": "padat",
      "latitude": -6.2088,
      "longitude": 106.8456
    }
  }' | jq '.'
```

---

### Cek E-Tilang Plat Tidak Terdaftar

```bash
curl -X POST http://localhost:8080/api/v1/chat \
  -H "Content-Type: application/json" \
  -d '{
    "message": "Bantu cek pelanggaran nomor polisi F1111XX",
    "context": {
      "location": "Bogor",
      "speed": 35.0,
      "traffic": "lancar",
      "latitude": -6.5971,
      "longitude": 106.8060
    }
  }' | jq '.'
```

**Expected Response:**
- ✅ Menampilkan data sebagai tidak ada pelanggaran (default untuk plat tidak terdaftar)

---

## 📊 Data Dummy E-Tilang

Backend memiliki data dummy untuk testing fitur e-tilang:

### 1. B1234SV - Ada 2 Pelanggaran ⚠️

- **Pemilik:** Budi Santoso
- **Jenis:** Motor Honda Beat
- **Rangka:** MH1RP6701FK123456
- **Pelanggaran:**
  1. Melanggar lampu merah (15 Des 2025) - Rp 500.000 - Belum dibayar ❌
  2. Tidak menggunakan helm SNI (20 Des 2025) - Rp 250.000 - Belum dibayar ❌
- **Total Denda:** Rp 750.000

### 2. B5678XY - Ada 1 Pelanggaran ⚠️

- **Pemilik:** Siti Rahayu
- **Jenis:** Mobil Toyota Avanza
- **Rangka:** MHKA42BA7JK098765
- **Pelanggaran:**
  1. Parkir di tempat terlarang (02 Jan 2026) - Rp 300.000 - Belum dibayar ❌
- **Total Denda:** Rp 300.000

### 3. B9999ZZ - Bersih ✅

- **Pemilik:** Ahmad Fauzi
- **Jenis:** Motor Yamaha NMAX
- **Rangka:** MH1JC5101FK234567
- **Pelanggaran:** Tidak ada
- **Total Denda:** Rp 0

### 4. D1111AA - Ada 1 Pelanggaran (Sudah Dibayar) ✅

- **Pemilik:** Rina Kartika
- **Jenis:** Mobil Honda CR-V
- **Rangka:** MHRGN81235K876543
- **Pelanggaran:**
  1. Melebihi batas kecepatan 120 km/jam (28 Des 2025) - Rp 500.000 - Sudah dibayar ✅
- **Total Denda:** Rp 500.000

### 5. E7777BB - Ada 1 Pelanggaran ⚠️

- **Pemilik:** Dedi Gunawan
- **Jenis:** Motor Kawasaki Ninja
- **Rangka:** MH1JFJ110FK345678
- **Pelanggaran:**
  1. Menggunakan handphone saat berkendara (05 Jan 2026) - Rp 750.000 - Belum dibayar ❌
- **Total Denda:** Rp 750.000

### Nomor Polisi Lainnya

- Semua nomor polisi yang **tidak terdaftar** akan dikembalikan sebagai **bersih/tidak ada pelanggaran**

---

## 🤖 Automated Test Script

Script bash otomatis untuk testing semua fitur sekaligus.

### Persiapan

```bash
# Berikan permission execute (jika belum)
chmod +x test-persona.sh
```

### Menjalankan Test

```bash
./test-persona.sh
```

### Apa yang Ditest?

Script akan otomatis menjalankan **8 test** berurutan:

1. ✅ Chat pertama - harus ada sapaan "Halo Sobat Lantas!"
2. ✅ Chat kedua dengan session sama - tanpa sapaan
3. ✅ Chat ketiga dengan pertanyaan berbeda - tanpa sapaan
4. ✅ Session baru - harus ada sapaan lagi
5. 🎫 E-Tilang dengan pelanggaran (B1234SV)
6. 🎫 E-Tilang bersih (B9999ZZ)
7. 🎫 E-Tilang pelanggaran parkir (B5678XY)
8. 🎫 E-Tilang plat tidak terdaftar (F1111XX)

---

## 🔍 Keyword Detection untuk E-Tilang

AI dapat mendeteksi permintaan e-tilang dengan berbagai keyword:

- ✅ "cek tilang nomor polisi..."
- ✅ "tolong bantu cek tilang dengan no polisi..."
- ✅ "ada pelanggaran tidak untuk kendaraan..."
- ✅ "cek denda tilang mobil..."
- ✅ "bantu cek pelanggaran nomor polisi..."
- ✅ "cek e-tilang..."

**Format Nomor Polisi:**

- Bisa dengan spasi: `B 1234 SV`
- Bisa tanpa spasi: `B1234SV`
- Case insensitive: `b1234sv` atau `B1234SV`

---

## 📝 Response Format

### Success Response

```json
{
  "success": true,
  "response": "Halo Sobat Lantas! Demi keselamatan, sebaiknya jangan bonceng dua anak kecil yaa...",
  "session_id": "550e8400-e29b-41d4-a716-446655440000"
}
```

### Error Response

```json
{
  "success": false,
  "error": "Message is required"
}
```

---

## 🎯 Karakteristik Persona "Sobat Lantas"

AI akan menjawab dengan karakteristik:

- ✅ **Ramah dan bersahabat** - seperti teman yang peduli
- ✅ **Peduli keselamatan** - fokus pada keselamatan pengguna dan keluarga
- ✅ **Bahasa santai tapi informatif** - menggunakan "yaa", "loh", "nih", "dong"
- ✅ **Tegas tapi ramah** - saat memberikan peringatan tetap bersahabat
- ✅ **Menggunakan emoji** - untuk visual clarity
- ✅ **Praktis dan mudah dipahami** - memberikan saran yang actionable

---

## 📞 API Endpoints

### POST /api/v1/chat

Endpoint utama untuk chat dengan AI.

**Request Body:**

```json
{
  "message": "string (required)",
  "context": {
    "location": "string",
    "speed": "number",
    "traffic": "string",
    "latitude": "number",
    "longitude": "number"
  },
  "session_id": "string (optional)"
}
```

**Response:**

```json
{
  "success": "boolean",
  "response": "string",
  "session_id": "string",
  "error": "string (optional)"
}
```

**Example:**

```bash
curl -X POST http://localhost:8080/api/v1/chat \
  -H "Content-Type: application/json" \
  -d '{
    "message": "Halo, boleh nanya?",
    "context": {
      "location": "Jakarta",
      "speed": 30.0,
      "traffic": "lancar"
    }
  }'
```

---

## 🚀 Quick Start Test

### Test Cepat Persona

```bash
curl -s -X POST http://localhost:8080/api/v1/chat \
  -H "Content-Type: application/json" \
  -d '{
    "message": "Halo, boleh nanya?",
    "context": {
      "location": "Jakarta",
      "speed": 30.0,
      "traffic": "lancar",
      "latitude": -6.2088,
      "longitude": 106.8456
    }
  }' | jq '.response'
```

### Test Cepat E-Tilang

```bash
curl -s -X POST http://localhost:8080/api/v1/chat \
  -H "Content-Type: application/json" \
  -d '{
    "message": "Cek tilang B1234SV",
    "context": {
      "location": "Jakarta",
      "speed": 30.0,
      "traffic": "lancar"
    }
  }' | jq '.response'
```

---

## ⚙️ Requirements

- Server running di `http://localhost:8080` (atau sesuaikan port)
- `curl` installed
- `jq` installed (untuk pretty print JSON) - optional

### Install jq (jika belum ada)

**macOS:**

```bash
brew install jq
```

**Linux:**

```bash
sudo apt-get install jq
```

**Windows (via Chocolatey):**

```bash
choco install jq
```

---

## 🐛 Troubleshooting

### Server tidak running

```bash
# Jalankan server terlebih dahulu
go run main.go
```

### Port berbeda

Edit script `test-persona.sh` dan ubah `localhost:8080` sesuai port Anda di file `.env`:

```bash
PORT=8080
```

### jq not found

Install jq atau hilangkan `| jq '.'` dari command curl:

```bash
# Tanpa jq
curl -X POST http://localhost:8080/api/v1/chat \
  -H "Content-Type: application/json" \
  -d '{"message": "test"}'
```

### Connection refused

Pastikan server sudah running dan port sudah benar:

```bash
# Cek apakah server running
lsof -i :8080

# Atau cek dengan curl ke health endpoint
curl http://localhost:8080/health
```

---

## 💡 Tips Testing

1. **Simpan Session ID** - untuk testing chat berkelanjutan dan melihat konsistensi persona
2. **Variasikan keyword** - untuk melihat kemampuan AI mendeteksi intent e-tilang
3. **Test dengan berbagai nomor polisi** - untuk melihat respons berbeda (ada pelanggaran vs bersih)
4. **Perhatikan persona** - harus konsisten ramah dan peduli keselamatan di semua respons
5. **Cek sapaan** - hanya muncul di pesan pertama per session, tidak boleh ada di chat lanjutan
6. **Test edge cases** - nomor polisi dengan format berbeda, typo, dll

---

## 📚 Additional Resources

- [README.md](README.md) - Project overview dan fitur
- [BUILD.md](BUILD.md) - Build dan deployment instructions
- [PORTAINER_DEPLOY.md](PORTAINER_DEPLOY.md) - Deployment dengan Portainer
- [test-persona.sh](test-persona.sh) - Automated test script

---

**Happy Testing! 🚓🚦**

_Last updated: January 6, 2026_
