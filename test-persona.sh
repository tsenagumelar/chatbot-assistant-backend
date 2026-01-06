#!/bin/bash

# Test Persona "Sobat Lantas" - Chatbot Assistant
# Pastikan server sudah running di localhost:8080

echo "═══════════════════════════════════════════════════"
echo "🧪 TESTING PERSONA SOBAT LANTAS"
echo "═══════════════════════════════════════════════════"
echo ""

# 1. TEST CHAT PERTAMA (harus ada sapaan "Halo Sobat Lantas!")
echo "📝 Test 1: Chat Pertama (harus ada sapaan 'Halo Sobat Lantas!')"
echo "───────────────────────────────────────────────────"
RESPONSE1=$(curl -s -X POST http://localhost:8080/api/v1/chat \
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
  }')

echo "$RESPONSE1" | jq '.'
SESSION_ID=$(echo "$RESPONSE1" | jq -r '.session_id')
echo ""
echo "Session ID: $SESSION_ID"
echo ""
sleep 2

# 2. TEST CHAT KEDUA (TIDAK boleh ada sapaan "Halo Sobat Lantas!")
echo "───────────────────────────────────────────────────"
echo "📝 Test 2: Chat Kedua dengan session yang sama (TIDAK boleh ada 'Halo Sobat Lantas!')"
echo "───────────────────────────────────────────────────"
curl -s -X POST http://localhost:8080/api/v1/chat \
  -H "Content-Type: application/json" \
  -d "{
    \"message\": \"Kalau satu anak kecil gimana?\",
    \"context\": {
      \"location\": \"Jl. Sudirman, Jakarta\",
      \"speed\": 35.5,
      \"traffic\": \"lancar\",
      \"latitude\": -6.2088,
      \"longitude\": 106.8456
    },
    \"session_id\": \"$SESSION_ID\"
  }" | jq '.'

echo ""
sleep 2

# 3. TEST CHAT KETIGA dengan konteks berbeda
echo "───────────────────────────────────────────────────"
echo "📝 Test 3: Chat Ketiga - pertanyaan berbeda"
echo "───────────────────────────────────────────────────"
curl -s -X POST http://localhost:8080/api/v1/chat \
  -H "Content-Type: application/json" \
  -d "{
    \"message\": \"Berapa batas kecepatan di tol?\",
    \"context\": {
      \"location\": \"Tol Jagorawi\",
      \"speed\": 85.0,
      \"traffic\": \"lancar\",
      \"latitude\": -6.2345,
      \"longitude\": 106.8765
    },
    \"session_id\": \"$SESSION_ID\"
  }" | jq '.'

echo ""
sleep 2

# 4. TEST SESSION BARU (harus ada sapaan lagi karena session baru)
echo "───────────────────────────────────────────────────"
echo "📝 Test 4: Session BARU (harus ada sapaan 'Halo Sobat Lantas!' lagi)"
echo "───────────────────────────────────────────────────"
curl -s -X POST http://localhost:8080/api/v1/chat \
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

echo ""
echo "═══════════════════════════════════════════════════"
echo "✅ TESTING SELESAI"
echo "═══════════════════════════════════════════════════"

# 5. TEST E-TILANG - Ada Pelanggaran
echo ""
echo "═══════════════════════════════════════════════════"
echo "🎫 TESTING E-TILANG FEATURE"
echo "═══════════════════════════════════════════════════"
echo ""
echo "📝 Test 5: Cek E-Tilang dengan pelanggaran (B1234SV)"
echo "───────────────────────────────────────────────────"
curl -s -X POST http://localhost:8080/api/v1/chat \
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

echo ""
sleep 2

# 6. TEST E-TILANG - Tidak Ada Pelanggaran
echo "───────────────────────────────────────────────────"
echo "📝 Test 6: Cek E-Tilang bersih (B9999ZZ)"
echo "───────────────────────────────────────────────────"
curl -s -X POST http://localhost:8080/api/v1/chat \
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

echo ""
sleep 2

# 7. TEST E-TILANG - Plat Lain dengan Pelanggaran
echo "───────────────────────────────────────────────────"
echo "📝 Test 7: Cek E-Tilang dengan pelanggaran parkir (B5678XY)"
echo "───────────────────────────────────────────────────"
curl -s -X POST http://localhost:8080/api/v1/chat \
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

echo ""
sleep 2

# 8. TEST E-TILANG - Plat tidak terdaftar
echo "───────────────────────────────────────────────────"
echo "📝 Test 8: Cek E-Tilang plat tidak terdaftar (F1111XX)"
echo "───────────────────────────────────────────────────"
curl -s -X POST http://localhost:8080/api/v1/chat \
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

echo ""
echo "═══════════════════════════════════════════════════"
echo "✅ ALL TESTS COMPLETED"
echo "═══════════════════════════════════════════════════"
