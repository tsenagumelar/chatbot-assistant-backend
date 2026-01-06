#!/bin/bash

# Colors for output
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

API_URL="http://localhost:8080/api/v1/chat"

echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo -e "${BLUE}  🪪 TEST SIM FLOW - Perpanjangan SIM${NC}"
echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo ""

# Step 1: Start conversation - trigger SIM flow
echo -e "${GREEN}Step 1: Trigger SIM flow${NC}"
echo -e "${YELLOW}Request: \"Saya mau perpanjang SIM\"${NC}"
RESPONSE=$(curl -s -X POST "$API_URL" \
  -H "Content-Type: application/json" \
  -d '{
    "message": "Saya mau perpanjang SIM",
    "context": {
      "location": "Jakarta",
      "latitude": -6.200000,
      "longitude": 106.816666
    }
  }')

echo "$RESPONSE" | jq '.'
SESSION_ID=$(echo "$RESPONSE" | jq -r '.session_id')
echo -e "\n${BLUE}Session ID: $SESSION_ID${NC}\n"
sleep 2

# Step 2: Pilih jenis SIM
echo -e "${GREEN}Step 2: Pilih jenis SIM${NC}"
echo -e "${YELLOW}Request: \"SIM A\"${NC}"
RESPONSE=$(curl -s -X POST "$API_URL" \
  -H "Content-Type: application/json" \
  -d "{
    \"message\": \"SIM A\",
    \"session_id\": \"$SESSION_ID\",
    \"context\": {
      \"location\": \"Jakarta\",
      \"latitude\": -6.200000,
      \"longitude\": 106.816666
    }
  }")

echo "$RESPONSE" | jq '.'
sleep 2

# Step 3: Jawab apakah pernah punya SIM sebelumnya
echo -e "\n${GREEN}Step 3: Jawab pernah punya SIM${NC}"
echo -e "${YELLOW}Request: \"Ya, pernah\"${NC}"
RESPONSE=$(curl -s -X POST "$API_URL" \
  -H "Content-Type: application/json" \
  -d "{
    \"message\": \"Ya, pernah\",
    \"session_id\": \"$SESSION_ID\",
    \"context\": {
      \"location\": \"Jakarta\",
      \"latitude\": -6.200000,
      \"longitude\": 106.816666
    }
  }")

echo "$RESPONSE" | jq '.'
sleep 2

# Step 4: Cek status SIM saat ini
echo -e "\n${GREEN}Step 4: Status SIM${NC}"
echo -e "${YELLOW}Request: \"Masih berlaku\"${NC}"
RESPONSE=$(curl -s -X POST "$API_URL" \
  -H "Content-Type: application/json" \
  -d "{
    \"message\": \"Masih berlaku\",
    \"session_id\": \"$SESSION_ID\",
    \"context\": {
      \"location\": \"Jakarta\",
      \"latitude\": -6.200000,
      \"longitude\": 106.816666
    }
  }")

echo "$RESPONSE" | jq '.'
sleep 2

# Step 5: Konfirmasi mau dibantu
echo -e "\n${GREEN}Step 5: Konfirmasi bantuan${NC}"
echo -e "${YELLOW}Request: \"Ya, tolong bantu\"${NC}"
RESPONSE=$(curl -s -X POST "$API_URL" \
  -H "Content-Type: application/json" \
  -d "{
    \"message\": \"Ya, tolong bantu\",
    \"session_id\": \"$SESSION_ID\",
    \"context\": {
      \"location\": \"Jakarta\",
      \"latitude\": -6.200000,
      \"longitude\": 106.816666
    }
  }")

echo "$RESPONSE" | jq '.'
sleep 2

# Step 6: Upload dokumen KTP
echo -e "\n${GREEN}Step 6: Upload KTP${NC}"
echo -e "${YELLOW}Request: Upload KTP (simulated)${NC}"
RESPONSE=$(curl -s -X POST "$API_URL" \
  -H "Content-Type: application/json" \
  -d "{
    \"message\": \"Ini KTP saya\",
    \"session_id\": \"$SESSION_ID\",
    \"context\": {
      \"location\": \"Jakarta\",
      \"latitude\": -6.200000,
      \"longitude\": 106.816666
    },
    \"documents\": [
      {
        \"file_name\": \"ktp.jpg\",
        \"file_type\": \"image/jpeg\",
        \"file_url\": \"https://example.com/uploads/ktp.jpg\",
        \"description\": \"KTP\"
      }
    ]
  }")

echo "$RESPONSE" | jq '.'
sleep 2

# Step 7: Upload SIM lama
echo -e "\n${GREEN}Step 7: Upload SIM lama${NC}"
echo -e "${YELLOW}Request: Upload SIM lama (simulated)${NC}"
RESPONSE=$(curl -s -X POST "$API_URL" \
  -H "Content-Type: application/json" \
  -d "{
    \"message\": \"Ini SIM lama saya\",
    \"session_id\": \"$SESSION_ID\",
    \"context\": {
      \"location\": \"Jakarta\",
      \"latitude\": -6.200000,
      \"longitude\": 106.816666
    },
    \"documents\": [
      {
        \"file_name\": \"sim_lama.jpg\",
        \"file_type\": \"image/jpeg\",
        \"file_url\": \"https://example.com/uploads/sim_lama.jpg\",
        \"description\": \"SIM Lama\"
      }
    ]
  }")

echo "$RESPONSE" | jq '.'
sleep 2

# Step 8: Upload surat keterangan sehat
echo -e "\n${GREEN}Step 8: Upload surat sehat${NC}"
echo -e "${YELLOW}Request: Upload surat sehat (simulated)${NC}"
RESPONSE=$(curl -s -X POST "$API_URL" \
  -H "Content-Type: application/json" \
  -d "{
    \"message\": \"Ini surat sehat saya\",
    \"session_id\": \"$SESSION_ID\",
    \"context\": {
      \"location\": \"Jakarta\",
      \"latitude\": -6.200000,
      \"longitude\": 106.816666
    },
    \"documents\": [
      {
        \"file_name\": \"surat_sehat.jpg\",
        \"file_type\": \"image/jpeg\",
        \"file_url\": \"https://example.com/uploads/surat_sehat.jpg\",
        \"description\": \"Surat Keterangan Sehat\"
      }
    ]
  }")

echo "$RESPONSE" | jq '.'
sleep 2

# Step 9: Upload hasil tes psikologi
echo -e "\n${GREEN}Step 9: Upload hasil tes psikologi${NC}"
echo -e "${YELLOW}Request: Upload tes psikologi (simulated)${NC}"
RESPONSE=$(curl -s -X POST "$API_URL" \
  -H "Content-Type: application/json" \
  -d "{
    \"message\": \"Ini hasil tes psikologi saya\",
    \"session_id\": \"$SESSION_ID\",
    \"context\": {
      \"location\": \"Jakarta\",
      \"latitude\": -6.200000,
      \"longitude\": 106.816666
    },
    \"documents\": [
      {
        \"file_name\": \"tes_psikologi.jpg\",
        \"file_type\": \"image/jpeg\",
        \"file_url\": \"https://example.com/uploads/tes_psikologi.jpg\",
        \"description\": \"Hasil Tes Psikologi\"
      }
    ]
  }")

echo "$RESPONSE" | jq '.'

echo -e "\n${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo -e "${BLUE}  ✅ Test SIM Flow Complete${NC}"
echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
