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
