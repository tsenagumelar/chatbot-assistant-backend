#!/bin/bash

# Test Script untuk Perubahan Terbaru
# 1. Greeting dengan nama user
# 2. Link e-tilang https://etle-pmj.id/
# 3. SIM hilang/rusak routing
# 4. SIM internasional routing
# 5. Mutasi kendaraan routing

BASE_URL="http://localhost:8080/api/v1"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo "═══════════════════════════════════════════════════"
echo "🧪 TEST PERUBAHAN TERBARU"
echo "═══════════════════════════════════════════════════"
echo ""

# ============================================================================
# TEST 1: GREETING DENGAN NAMA USER (bukan "Halo Sobat Lantas")
# ============================================================================
echo -e "${BLUE}════════════════════════════════════════════════════════════════${NC}"
echo -e "${GREEN}TEST 1: Greeting dengan Nama User${NC}"
echo -e "${BLUE}════════════════════════════════════════════════════════════════${NC}"
echo ""

echo -e "${YELLOW}Expected: Harus menyapa 'Halo Taufan!' bukan 'Halo Sobat Lantas!'${NC}"
curl -X POST "$BASE_URL/chat" \
  -H "Content-Type: application/json" \
  -d '{
    "session_id": "test-greeting-001",
    "message": "halo",
    "name": "Taufan",
    "location": "Jakarta Selatan"
  }' | jq '.'

echo ""
echo -e "${YELLOW}Tunggu 3 detik...${NC}"
sleep 3

# ============================================================================
# TEST 2: GREETING TANPA NAMA (fallback ke Sobat Lantas)
# ============================================================================
echo ""
echo -e "${BLUE}════════════════════════════════════════════════════════════════${NC}"
echo -e "${GREEN}TEST 2: Greeting Tanpa Nama (Fallback ke Sobat Lantas)${NC}"
echo -e "${BLUE}════════════════════════════════════════════════════════════════${NC}"
echo ""

echo -e "${YELLOW}Expected: Harus menyapa 'Halo Sobat Lantas!' karena nama tidak ada${NC}"
curl -X POST "$BASE_URL/chat" \
  -H "Content-Type: application/json" \
  -d '{
    "session_id": "test-greeting-002",
    "message": "halo",
    "location": "Jakarta Pusat"
  }' | jq '.'

echo ""
echo -e "${YELLOW}Tunggu 3 detik...${NC}"
sleep 3

# ============================================================================
# TEST 3: CEK E-TILANG - HARUS ADA LINK https://etle-pmj.id/
# ============================================================================
echo ""
echo -e "${BLUE}════════════════════════════════════════════════════════════════${NC}"
echo -e "${GREEN}TEST 3: Cek E-Tilang - Link Pembayaran${NC}"
echo -e "${BLUE}════════════════════════════════════════════════════════════════${NC}"
echo ""

echo -e "${YELLOW}Expected: Response harus ada link 'https://etle-pmj.id/' untuk bayar${NC}"
curl -X POST "$BASE_URL/chat" \
  -H "Content-Type: application/json" \
  -d '{
    "session_id": "test-etilang-link-001",
    "message": "cek tilang B1234SV",
    "name": "Andi",
    "location": "Jakarta"
  }' | jq '.'

echo ""
echo -e "${YELLOW}Tunggu 3 detik...${NC}"
sleep 3

# ============================================================================
# TEST 4: SIM HILANG - ROUTING KE "SIM Hilang / Rusak"
# ============================================================================
echo ""
echo -e "${BLUE}════════════════════════════════════════════════════════════════${NC}"
echo -e "${GREEN}TEST 4: SIM Hilang - Routing Correct${NC}"
echo -e "${BLUE}════════════════════════════════════════════════════════════════${NC}"
echo ""

echo -e "${YELLOW}Expected: Harus ke flow 'SIM Hilang / Rusak', BUKAN 'Buat SIM'${NC}"
curl -X POST "$BASE_URL/chat" \
  -H "Content-Type: application/json" \
  -d '{
    "session_id": "test-sim-hilang-001",
    "message": "SIM saya hilang bisa dibantu?",
    "name": "Budi",
    "location": "Bandung"
  }' | jq '.'

echo ""
echo -e "${YELLOW}Tunggu 3 detik...${NC}"
sleep 3

# ============================================================================
# TEST 5: SIM RUSAK - ROUTING KE "SIM Hilang / Rusak"
# ============================================================================
echo ""
echo -e "${BLUE}════════════════════════════════════════════════════════════════${NC}"
echo -e "${GREEN}TEST 5: SIM Rusak - Routing Correct${NC}"
echo -e "${BLUE}════════════════════════════════════════════════════════════════${NC}"
echo ""

echo -e "${YELLOW}Expected: Harus ke flow 'SIM Hilang / Rusak', BUKAN 'Buat SIM'${NC}"
curl -X POST "$BASE_URL/chat" \
  -H "Content-Type: application/json" \
  -d '{
    "session_id": "test-sim-rusak-001",
    "message": "SIM saya rusak mau ganti",
    "name": "Citra",
    "location": "Surabaya"
  }' | jq '.'

echo ""
echo -e "${YELLOW}Tunggu 3 detik...${NC}"
sleep 3

# ============================================================================
# TEST 6: SIM INTERNASIONAL - ROUTING KE "SIM Internasional"
# ============================================================================
echo ""
echo -e "${BLUE}════════════════════════════════════════════════════════════════${NC}"
echo -e "${GREEN}TEST 6: SIM Internasional - Routing Correct${NC}"
echo -e "${BLUE}════════════════════════════════════════════════════════════════${NC}"
echo ""

echo -e "${YELLOW}Expected: Harus ke flow 'SIM Internasional', BUKAN 'Buat SIM'${NC}"
curl -X POST "$BASE_URL/chat" \
  -H "Content-Type: application/json" \
  -d '{
    "session_id": "test-sim-intl-001",
    "message": "mau bikin SIM internasional",
    "name": "Deni",
    "location": "Jakarta"
  }' | jq '.'

echo ""
echo -e "${YELLOW}Tunggu 3 detik...${NC}"
sleep 3

# ============================================================================
# TEST 7: MUTASI KENDARAAN - ROUTING KE "Mutasi Kendaraan"
# ============================================================================
echo ""
echo -e "${BLUE}════════════════════════════════════════════════════════════════${NC}"
echo -e "${GREEN}TEST 7: Mutasi Kendaraan - Routing Correct${NC}"
echo -e "${BLUE}════════════════════════════════════════════════════════════════${NC}"
echo ""

echo -e "${YELLOW}Expected: Harus ke flow 'Mutasi Kendaraan', BUKAN 'Balik Nama'${NC}"
curl -X POST "$BASE_URL/chat" \
  -H "Content-Type: application/json" \
  -d '{
    "session_id": "test-mutasi-001",
    "message": "saya mau mutasi kendaraan",
    "name": "Eko",
    "location": "Semarang"
  }' | jq '.'

echo ""
echo -e "${YELLOW}Tunggu 3 detik...${NC}"
sleep 3

# ============================================================================
# TEST 8: BALIK NAMA - ROUTING KE "Balik Nama Kendaraan"
# ============================================================================
echo ""
echo -e "${BLUE}════════════════════════════════════════════════════════════════${NC}"
echo -e "${GREEN}TEST 8: Balik Nama - Routing Correct${NC}"
echo -e "${BLUE}════════════════════════════════════════════════════════════════${NC}"
echo ""

echo -e "${YELLOW}Expected: Harus ke flow 'Balik Nama Kendaraan', BUKAN 'Mutasi'${NC}"
curl -X POST "$BASE_URL/chat" \
  -H "Content-Type: application/json" \
  -d '{
    "session_id": "test-balikNama-001",
    "message": "mau balik nama motor bisa dibantu?",
    "name": "Fitri",
    "location": "Yogyakarta"
  }' | jq '.'

echo ""
echo -e "${YELLOW}Tunggu 3 detik...${NC}"
sleep 3

# ============================================================================
# TEST 9: GREETING LANJUTAN - TIDAK ADA SAPAAN LAGI
# ============================================================================
echo ""
echo -e "${BLUE}════════════════════════════════════════════════════════════════${NC}"
echo -e "${GREEN}TEST 9: Greeting Lanjutan - No Greeting${NC}"
echo -e "${BLUE}════════════════════════════════════════════════════════════════${NC}"
echo ""

echo -e "${YELLOW}Step 1: Pesan pertama (ada greeting)${NC}"
curl -X POST "$BASE_URL/chat" \
  -H "Content-Type: application/json" \
  -d '{
    "session_id": "test-no-greeting-001",
    "message": "halo",
    "name": "Gita",
    "location": "Malang"
  }' | jq '.'

echo ""
echo -e "${YELLOW}Tunggu 2 detik...${NC}"
sleep 2

echo -e "${YELLOW}Step 2: Pesan lanjutan (TIDAK ada greeting)${NC}"
curl -X POST "$BASE_URL/chat" \
  -H "Content-Type: application/json" \
  -d '{
    "session_id": "test-no-greeting-001",
    "message": "mau tanya cara perpanjang SIM",
    "name": "Gita",
    "location": "Malang"
  }' | jq '.'

echo ""
echo -e "${YELLOW}Tunggu 3 detik...${NC}"
sleep 3

# ============================================================================
# SUMMARY
# ============================================================================
echo ""
echo "═══════════════════════════════════════════════════"
echo "✅ ALL TESTS COMPLETED"
echo "═══════════════════════════════════════════════════"
echo ""
echo "📋 CHECKLIST YANG HARUS DICEK:"
echo ""
echo "✓ TEST 1-2: Greeting dengan nama user (bukan Sobat Lantas)"
echo "✓ TEST 3: Link e-tilang ada 'https://etle-pmj.id/'"
echo "✓ TEST 4-5: SIM hilang/rusak routing ke flow yang benar"
echo "✓ TEST 6: SIM internasional routing ke flow yang benar"
echo "✓ TEST 7: Mutasi kendaraan routing ke flow yang benar"
echo "✓ TEST 8: Balik nama routing ke flow yang benar"
echo "✓ TEST 9: Pesan lanjutan tidak ada greeting lagi"
echo ""
echo "═══════════════════════════════════════════════════"
