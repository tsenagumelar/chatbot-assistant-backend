#!/bin/bash

echo "🧪 Test Session Management - Backend History Storage"
echo "======================================================"
echo ""

# Colors
GREEN='\033[0;32m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

API_URL="http://localhost:8080/api/v1"

echo -e "${BLUE}Test 1: Kirim pesan pertama (tanpa session_id - auto create)${NC}"
echo "================================================================"
RESPONSE1=$(curl -s -X POST $API_URL/chat \
  -H "Content-Type: application/json" \
  -d '{
    "message": "Halo, nama saya Taufan",
    "session_id": "",
    "context": {
      "location": "Jakarta Selatan",
      "latitude": -6.2608,
      "longitude": 106.7819,
      "speed": 0,
      "traffic": "smooth"
    }
  }')

echo "$RESPONSE1" | jq '.'
SESSION_ID=$(echo "$RESPONSE1" | jq -r '.session_id')

echo ""
echo -e "${GREEN}✅ Session ID: $SESSION_ID${NC}"
echo ""
read -p "Tekan Enter untuk lanjut ke test 2..."
echo ""

echo -e "${BLUE}Test 2: Tanya nama (dengan session_id - AI harus ingat)${NC}"
echo "=========================================================="
RESPONSE2=$(curl -s -X POST $API_URL/chat \
  -H "Content-Type: application/json" \
  -d "{
    \"message\": \"Siapa nama saya?\",
    \"session_id\": \"$SESSION_ID\",
    \"context\": {
      \"location\": \"Jakarta Selatan\",
      \"latitude\": -6.2608,
      \"longitude\": 106.7819,
      \"speed\": 0,
      \"traffic\": \"smooth\"
    }
  }")

echo "$RESPONSE2" | jq '.'
echo ""
read -p "Tekan Enter untuk lanjut ke test 3..."
echo ""

echo -e "${BLUE}Test 3: Get Session Info${NC}"
echo "============================"
curl -s $API_URL/session/$SESSION_ID | jq '.'
echo ""
read -p "Tekan Enter untuk lanjut ke test 4..."
echo ""

echo -e "${BLUE}Test 4: Skenario Samsat - Pertama${NC}"
echo "======================================"
RESPONSE3=$(curl -s -X POST $API_URL/chat \
  -H "Content-Type: application/json" \
  -d "{
    \"message\": \"Saya mau ke Kantor Samsat Tangsel\",
    \"session_id\": \"$SESSION_ID\",
    \"context\": {
      \"location\": \"Jakarta Selatan\",
      \"latitude\": -6.2608,
      \"longitude\": 106.7819,
      \"speed\": 0,
      \"traffic\": \"smooth\"
    }
  }")

echo "$RESPONSE3" | jq '.'
echo ""
read -p "Tekan Enter untuk lanjut ke test 5..."
echo ""

echo -e "${BLUE}Test 5: Skenario Samsat - Follow-up (AI harus ingat tujuan)${NC}"
echo "=============================================================="
RESPONSE4=$(curl -s -X POST $API_URL/chat \
  -H "Content-Type: application/json" \
  -d "{
    \"message\": \"Berapa jaraknya dari lokasi saya?\",
    \"session_id\": \"$SESSION_ID\",
    \"context\": {
      \"location\": \"Jakarta Selatan\",
      \"latitude\": -6.2608,
      \"longitude\": 106.7819,
      \"speed\": 0,
      \"traffic\": \"smooth\"
    }
  }")

echo "$RESPONSE4" | jq '.'
echo ""
read -p "Tekan Enter untuk lanjut ke test 6..."
echo ""

echo -e "${BLUE}Test 6: Get Session Info (setelah banyak pesan)${NC}"
echo "================================================="
curl -s $API_URL/session/$SESSION_ID | jq '.'
echo ""
read -p "Tekan Enter untuk lanjut ke test 7..."
echo ""

echo -e "${BLUE}Test 7: Clear Session History${NC}"
echo "================================"
curl -s -X POST $API_URL/session/$SESSION_ID/clear | jq '.'
echo ""
read -p "Tekan Enter untuk lanjut ke test 8..."
echo ""

echo -e "${BLUE}Test 8: Tanya nama lagi (setelah clear - AI tidak ingat)${NC}"
echo "=========================================================="
RESPONSE5=$(curl -s -X POST $API_URL/chat \
  -H "Content-Type: application/json" \
  -d "{
    \"message\": \"Siapa nama saya?\",
    \"session_id\": \"$SESSION_ID\",
    \"context\": {
      \"location\": \"Jakarta Selatan\",
      \"latitude\": -6.2608,
      \"longitude\": 106.7819,
      \"speed\": 0,
      \"traffic\": \"smooth\"
    }
  }")

echo "$RESPONSE5" | jq '.'
echo ""
echo -e "${GREEN}✅ Test selesai!${NC}"
echo ""
echo "📋 Ringkasan:"
echo "- Session ID: $SESSION_ID"
echo "- Test 2: AI harus ingat nama 'Taufan' ✅"
echo "- Test 5: AI harus ingat tujuan 'Samsat Tangsel' ✅"
echo "- Test 8: AI tidak ingat nama setelah clear ✅"
