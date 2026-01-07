#!/bin/bash

# Test Script for Korlantas Knowledge Base
# Test pertanyaan tentang pejabat Korlantas Polri

BASE_URL="http://localhost:8080/api/v1"

echo "════════════════════════════════════════════════════════════════"
echo "🚓 TEST KORLANTAS KNOWLEDGE BASE"
echo "════════════════════════════════════════════════════════════════"
echo ""

# Colors
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# ============================================================================
# TEST 1: KAKORLANTAS
# ============================================================================
echo -e "${BLUE}════════════════════════════════════════════════════════════════${NC}"
echo -e "${GREEN}TEST 1: Siapa Kakorlantas Polri saat ini?${NC}"
echo -e "${BLUE}════════════════════════════════════════════════════════════════${NC}"
echo ""

curl -X POST "$BASE_URL/chat" \
  -H "Content-Type: application/json" \
  -d '{
    "session_id": "test-korlantas-001",
    "message": "siapa kakorlantas sekarang?",
    "name": "Taufan",
    "location": "Jakarta"
  }' | jq '.'

echo ""
echo -e "${YELLOW}Tunggu 3 detik...${NC}"
sleep 3

# ============================================================================
# TEST 2: DIRLANTAS POLDA JABAR
# ============================================================================
echo ""
echo -e "${BLUE}════════════════════════════════════════════════════════════════${NC}"
echo -e "${GREEN}TEST 2: Dirlantas Polda Jabar${NC}"
echo -e "${BLUE}════════════════════════════════════════════════════════════════${NC}"
echo ""

curl -X POST "$BASE_URL/chat" \
  -H "Content-Type: application/json" \
  -d '{
    "session_id": "test-korlantas-002",
    "message": "siapa dirlantas polda jawa barat?",
    "name": "Budi",
    "location": "Bandung"
  }' | jq '.'

echo ""
echo -e "${YELLOW}Tunggu 3 detik...${NC}"
sleep 3

# ============================================================================
# TEST 3: DIRLANTAS POLDA METRO JAYA
# ============================================================================
echo ""
echo -e "${BLUE}════════════════════════════════════════════════════════════════${NC}"
echo -e "${GREEN}TEST 3: Dirlantas Polda Metro Jaya${NC}"
echo -e "${BLUE}════════════════════════════════════════════════════════════════${NC}"
echo ""

curl -X POST "$BASE_URL/chat" \
  -H "Content-Type: application/json" \
  -d '{
    "session_id": "test-korlantas-003",
    "message": "siapa kepala lalu lintas polda metro jaya?",
    "name": "Siti",
    "location": "Jakarta Selatan"
  }' | jq '.'

echo ""
echo -e "${YELLOW}Tunggu 3 detik...${NC}"
sleep 3

# ============================================================================
# TEST 4: DIRLANTAS POLDA BANTEN
# ============================================================================
echo ""
echo -e "${BLUE}════════════════════════════════════════════════════════════════${NC}"
echo -e "${GREEN}TEST 4: Dirlantas Polda Banten${NC}"
echo -e "${BLUE}════════════════════════════════════════════════════════════════${NC}"
echo ""

curl -X POST "$BASE_URL/chat" \
  -H "Content-Type: application/json" \
  -d '{
    "session_id": "test-korlantas-004",
    "message": "dirlantas banten siapa ya?",
    "name": "Ahmad",
    "location": "Tangerang"
  }' | jq '.'

echo ""
echo -e "${YELLOW}Tunggu 3 detik...${NC}"
sleep 3

# ============================================================================
# TEST 5: DIRLANTAS POLDA JATENG
# ============================================================================
echo ""
echo -e "${BLUE}════════════════════════════════════════════════════════════════${NC}"
echo -e "${GREEN}TEST 5: Dirlantas Polda Jawa Tengah${NC}"
echo -e "${BLUE}════════════════════════════════════════════════════════════════${NC}"
echo ""

curl -X POST "$BASE_URL/chat" \
  -H "Content-Type: application/json" \
  -d '{
    "session_id": "test-korlantas-005",
    "message": "siapa nama dirlantas polda jateng?",
    "name": "Rina",
    "location": "Semarang"
  }' | jq '.'

echo ""
echo -e "${YELLOW}Tunggu 3 detik...${NC}"
sleep 3

# ============================================================================
# TEST 6: DIRLANTAS POLDA JATIM
# ============================================================================
echo ""
echo -e "${BLUE}════════════════════════════════════════════════════════════════${NC}"
echo -e "${GREEN}TEST 6: Dirlantas Polda Jawa Timur${NC}"
echo -e "${BLUE}════════════════════════════════════════════════════════════════${NC}"
echo ""

curl -X POST "$BASE_URL/chat" \
  -H "Content-Type: application/json" \
  -d '{
    "session_id": "test-korlantas-006",
    "message": "dirlantas jawa timur siapa?",
    "name": "Deni",
    "location": "Surabaya"
  }' | jq '.'

echo ""
echo -e "${YELLOW}Tunggu 3 detik...${NC}"
sleep 3

# ============================================================================
# TEST 7: DIRLANTAS POLDA SUMUT
# ============================================================================
echo ""
echo -e "${BLUE}════════════════════════════════════════════════════════════════${NC}"
echo -e "${GREEN}TEST 7: Dirlantas Polda Sumatera Utara${NC}"
echo -e "${BLUE}════════════════════════════════════════════════════════════════${NC}"
echo ""

curl -X POST "$BASE_URL/chat" \
  -H "Content-Type: application/json" \
  -d '{
    "session_id": "test-korlantas-007",
    "message": "siapa dirlantas polda sumut?",
    "name": "Hendra",
    "location": "Medan"
  }' | jq '.'

echo ""
echo -e "${YELLOW}Tunggu 3 detik...${NC}"
sleep 3

# ============================================================================
# TEST 8: DIRLANTAS POLDA BALI
# ============================================================================
echo ""
echo -e "${BLUE}════════════════════════════════════════════════════════════════${NC}"
echo -e "${GREEN}TEST 8: Dirlantas Polda Bali${NC}"
echo -e "${BLUE}════════════════════════════════════════════════════════════════${NC}"
echo ""

curl -X POST "$BASE_URL/chat" \
  -H "Content-Type: application/json" \
  -d '{
    "session_id": "test-korlantas-008",
    "message": "kepala lalu lintas bali siapa?",
    "name": "Fitri",
    "location": "Denpasar"
  }' | jq '.'

echo ""
echo -e "${YELLOW}Tunggu 3 detik...${NC}"
sleep 3

# ============================================================================
# TEST 9: KASUBDIT BPKB
# ============================================================================
echo ""
echo -e "${BLUE}════════════════════════════════════════════════════════════════${NC}"
echo -e "${GREEN}TEST 9: Kasubdit BPKB${NC}"
echo -e "${BLUE}════════════════════════════════════════════════════════════════${NC}"
echo ""

curl -X POST "$BASE_URL/chat" \
  -H "Content-Type: application/json" \
  -d '{
    "session_id": "test-korlantas-009",
    "message": "siapa kepala subdit BPKB?",
    "name": "Eko",
    "location": "Jakarta"
  }' | jq '.'

echo ""
echo -e "${YELLOW}Tunggu 3 detik...${NC}"
sleep 3

# ============================================================================
# TEST 10: KASUBDIT STNK
# ============================================================================
echo ""
echo -e "${BLUE}════════════════════════════════════════════════════════════════${NC}"
echo -e "${GREEN}TEST 10: Kasubdit STNK${NC}"
echo -e "${BLUE}════════════════════════════════════════════════════════════════${NC}"
echo ""

curl -X POST "$BASE_URL/chat" \
  -H "Content-Type: application/json" \
  -d '{
    "session_id": "test-korlantas-010",
    "message": "kasubdit stnk siapa namanya?",
    "name": "Rudi",
    "location": "Jakarta"
  }' | jq '.'

echo ""
echo -e "${YELLOW}Tunggu 3 detik...${NC}"
sleep 3

# ============================================================================
# TEST 11: DIRLANTAS POLDA SULSEL
# ============================================================================
echo ""
echo -e "${BLUE}════════════════════════════════════════════════════════════════${NC}"
echo -e "${GREEN}TEST 11: Dirlantas Polda Sulawesi Selatan${NC}"
echo -e "${BLUE}════════════════════════════════════════════════════════════════${NC}"
echo ""

curl -X POST "$BASE_URL/chat" \
  -H "Content-Type: application/json" \
  -d '{
    "session_id": "test-korlantas-011",
    "message": "dirlantas sulawesi selatan siapa?",
    "name": "Wahyu",
    "location": "Makassar"
  }' | jq '.'

echo ""
echo -e "${YELLOW}Tunggu 3 detik...${NC}"
sleep 3

# ============================================================================
# TEST 12: DIRLANTAS POLDA PAPUA
# ============================================================================
echo ""
echo -e "${BLUE}════════════════════════════════════════════════════════════════${NC}"
echo -e "${GREEN}TEST 12: Dirlantas Polda Papua${NC}"
echo -e "${BLUE}════════════════════════════════════════════════════════════════${NC}"
echo ""

curl -X POST "$BASE_URL/chat" \
  -H "Content-Type: application/json" \
  -d '{
    "session_id": "test-korlantas-012",
    "message": "siapa dirlantas papua?",
    "name": "Tono",
    "location": "Jayapura"
  }' | jq '.'

echo ""
echo -e "${YELLOW}Tunggu 3 detik...${NC}"
sleep 3

# ============================================================================
# TEST 13: MULTIPLE QUESTIONS - CONTEXT AWARENESS
# ============================================================================
echo ""
echo -e "${BLUE}════════════════════════════════════════════════════════════════${NC}"
echo -e "${GREEN}TEST 13: Context Awareness - Multiple Questions${NC}"
echo -e "${BLUE}════════════════════════════════════════════════════════════════${NC}"
echo ""

echo -e "${YELLOW}Step 1: Tanya Kakorlantas${NC}"
curl -X POST "$BASE_URL/chat" \
  -H "Content-Type: application/json" \
  -d '{
    "session_id": "test-korlantas-013",
    "message": "siapa kakorlantas?",
    "name": "Lisa",
    "location": "Jakarta"
  }' | jq '.'

echo ""
echo -e "${YELLOW}Tunggu 3 detik...${NC}"
sleep 3

echo -e "${YELLOW}Step 2: Tanya Dirlantas (harusnya ingat context)${NC}"
curl -X POST "$BASE_URL/chat" \
  -H "Content-Type: application/json" \
  -d '{
    "session_id": "test-korlantas-013",
    "message": "kalau dirlantas polda jabar siapa?",
    "name": "Lisa",
    "location": "Jakarta"
  }' | jq '.'

echo ""
echo -e "${YELLOW}Tunggu 3 detik...${NC}"
sleep 3

# ============================================================================
# TEST 14: DIRLANTAS POLDA KALTIM
# ============================================================================
echo ""
echo -e "${BLUE}════════════════════════════════════════════════════════════════${NC}"
echo -e "${GREEN}TEST 14: Dirlantas Polda Kalimantan Timur${NC}"
echo -e "${BLUE}════════════════════════════════════════════════════════════════${NC}"
echo ""

curl -X POST "$BASE_URL/chat" \
  -H "Content-Type: application/json" \
  -d '{
    "session_id": "test-korlantas-014",
    "message": "siapa dirlantas kaltim?",
    "name": "Agus",
    "location": "Balikpapan"
  }' | jq '.'

echo ""
echo -e "${YELLOW}Tunggu 3 detik...${NC}"
sleep 3

# ============================================================================
# TEST 15: KABAGOPS
# ============================================================================
echo ""
echo -e "${BLUE}════════════════════════════════════════════════════════════════${NC}"
echo -e "${GREEN}TEST 15: Kepala Bagian Operasional Korlantas${NC}"
echo -e "${BLUE}════════════════════════════════════════════════════════════════${NC}"
echo ""

curl -X POST "$BASE_URL/chat" \
  -H "Content-Type: application/json" \
  -d '{
    "session_id": "test-korlantas-015",
    "message": "siapa kabagops korlantas?",
    "name": "Dewi",
    "location": "Jakarta"
  }' | jq '.'

echo ""
echo ""
echo -e "${GREEN}════════════════════════════════════════════════════════════════${NC}"
echo -e "${GREEN}✅ SEMUA TEST KORLANTAS KNOWLEDGE SELESAI!${NC}"
echo -e "${GREEN}════════════════════════════════════════════════════════════════${NC}"
echo ""
echo "📝 VALIDASI RESPONSE:"
echo "   ✅ Nama pejabat lengkap dengan gelar"
echo "   ✅ Tidak ada character # atau * untuk formatting"
echo "   ✅ Ada disclaimer untuk cek website resmi"
echo "   ✅ Menyapa dengan nama user yang diberikan"
echo "   ✅ Bahasa santai tapi informatif"
echo ""
echo "🎯 EXPECTED ANSWERS:"
echo "   Kakorlantas: Irjen. Pol. Drs. Agus Suryonugroho, S.H., M.Hum."
echo "   Dirlantas Jabar: Kombes. Pol. Dodi Darjanto, S.I.K., M.H."
echo "   Dirlantas Metro Jaya: Kombes. Pol. Komarudin, S.I.K., M.M."
echo "   Dirlantas Banten: Kombes. Pol. Dr. Leganek Mawardi, S.H., S.I.K., M.Si."
echo "   Kasubdit BPKB: Kombes. Pol. Sumardji, S.H."
echo "   Kasubdit STNK: Kombes. Pol. Dedy Suhartono, S.I.K., M.M."
echo ""
