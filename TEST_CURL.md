# Test Persona Sobat Lantas - CURL Examples

## 1. Chat Pertama (Harus ada sapaan "Halo Sobat Lantas!")

```bash
curl -X POST http://localhost:3000/api/v1/chat \
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
  }'
```

**Expected Response:**
- Harus dimulai dengan: "Halo Sobat Lantas!"
- Isi respons ramah dan peduli keselamatan
- Return `session_id` untuk chat selanjutnya

---

## 2. Chat Kedua (TIDAK boleh ada "Halo Sobat Lantas!")

**Ganti `YOUR_SESSION_ID` dengan session_id dari response pertama**

```bash
curl -X POST http://localhost:3000/api/v1/chat \
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
  }'
```

**Expected Response:**
- TIDAK ada sapaan "Halo Sobat Lantas!"
- Langsung jawab pertanyaan dengan ramah
- Tetap menggunakan persona yang peduli

---

## 3. Chat Lanjutan dengan Konteks Berbeda

```bash
curl -X POST http://localhost:3000/api/v1/chat \
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
  }'
```

---

## 4. Session Baru (Harus ada sapaan lagi)

```bash
curl -X POST http://localhost:3000/api/v1/chat \
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
  }'
```

**Expected Response:**
- Harus ada sapaan "Halo Sobat Lantas!" lagi (karena session baru)

---

## Test dengan Pretty Print (menggunakan jq)

```bash
curl -s -X POST http://localhost:3000/api/v1/chat \
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

---

## Contoh Response Format

```json
{
  "success": true,
  "response": "Halo Sobat Lantas! Demi keselamatan, sebaiknya jangan bonceng dua anak kecil yaa. Bahaya banget loh. Anak-anak harus pakai helm SNI dan cukup satu saja yang dibonceng. Utamakan keselamatan keluarga kita!",
  "session_id": "550e8400-e29b-41d4-a716-446655440000"
}
```

---

## Menjalankan Automated Test Script

```bash
# Berikan permission execute
chmod +x test-persona.sh

# Jalankan test
./test-persona.sh
```

Script akan otomatis:
1. Test chat pertama (ada sapaan)
2. Test chat kedua dengan session sama (tanpa sapaan)
3. Test chat ketiga (tanpa sapaan)
4. Test session baru (ada sapaan lagi)
