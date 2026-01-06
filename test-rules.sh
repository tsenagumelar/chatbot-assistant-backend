#!/bin/bash

# Test Script for Response Rules and Location Rules
# Test berbagai jenis pelayanan dengan conversation flow dari response-rules.json

BASE_URL="http://localhost:8080/api/v1"

echo "════════════════════════════════════════════════════════════════"
echo "🧪 TEST RULES-BASED CONVERSATION FLOW"
echo "════════════════════════════════════════════════════════════════"
echo ""

# Colors
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# ============================================================================
# TEST 1: BUAT SIM BARU (with name parameter)
# ============================================================================
echo -e "${BLUE}════════════════════════════════════════════════════════════════${NC}"
echo -e "${GREEN}TEST 1: Buat SIM Baru (dengan nama: Budi)${NC}"
echo -e "${BLUE}════════════════════════════════════════════════════════════════${NC}"
echo ""

echo -e "${YELLOW}Step 1: User bertanya tentang pembuatan SIM${NC}"
curl -X POST "$BASE_URL/chat" \
  -H "Content-Type: application/json" \
  -d '{
    "session_id": "test-sim-baru-001",
    "message": "mau bikin sim baru",
    "name": "Budi",
    "location": "Jakarta Selatan",
    "latitude": -6.2608,
    "longitude": 106.7817,
    "speed": 0,
    "traffic": "smooth"
  }' | jq '.'

echo ""
echo -e "${YELLOW}Tunggu 3 detik...${NC}"
sleep 3

echo -e "${YELLOW}Step 2: User mau dibantu${NC}"
curl -X POST "$BASE_URL/chat" \
  -H "Content-Type: application/json" \
  -d '{
    "session_id": "test-sim-baru-001",
    "message": "iya mau dibantu",
    "name": "Budi",
    "location": "Jakarta Selatan",
    "latitude": -6.2608,
    "longitude": 106.7817
  }' | jq '.'

echo ""
echo -e "${YELLOW}Tunggu 3 detik...${NC}"
sleep 3

# ============================================================================
# TEST 2: PERPANJANGAN SIM (with name and document upload)
# ============================================================================
echo ""
echo -e "${BLUE}════════════════════════════════════════════════════════════════${NC}"
echo -e "${GREEN}TEST 2: Perpanjangan SIM (dengan nama: Siti, upload dokumen)${NC}"
echo -e "${BLUE}════════════════════════════════════════════════════════════════${NC}"
echo ""

echo -e "${YELLOW}Step 1: User bertanya tentang perpanjangan SIM${NC}"
curl -X POST "$BASE_URL/chat" \
  -H "Content-Type: application/json" \
  -d '{
    "session_id": "test-perpanjang-sim-002",
    "message": "bagaimana cara perpanjang SIM?",
    "name": "Siti",
    "location": "Tangerang Selatan",
    "latitude": -6.2933,
    "longitude": 106.6894
  }' | jq '.'

echo ""
echo -e "${YELLOW}Tunggu 3 detik...${NC}"
sleep 3

echo -e "${YELLOW}Step 2: User upload dokumen${NC}"
curl -X POST "$BASE_URL/chat" \
  -H "Content-Type: multipart/form-data" \
  -F "session_id=test-perpanjang-sim-002" \
  -F "message=ini dokumen saya" \
  -F "name=Siti" \
  -F "location=Tangerang Selatan" \
  -F "latitude=-6.2933" \
  -F "longitude=106.6894" \
  -F "documents=@test-ktp.jpg" \
  -F "documents=@test-sim-lama.jpg" | jq '.'

echo ""
echo -e "${YELLOW}Tunggu 3 detik...${NC}"
sleep 3

# ============================================================================
# TEST 3: PAJAK KENDARAAN BERMOTOR
# ============================================================================
echo ""
echo -e "${BLUE}════════════════════════════════════════════════════════════════${NC}"
echo -e "${GREEN}TEST 3: Pajak Kendaraan Bermotor (dengan nama: Ahmad)${NC}"
echo -e "${BLUE}════════════════════════════════════════════════════════════════${NC}"
echo ""

echo -e "${YELLOW}Step 1: User bertanya tentang pajak kendaraan${NC}"
curl -X POST "$BASE_URL/chat" \
  -H "Content-Type: application/json" \
  -d '{
    "session_id": "test-pajak-003",
    "message": "cara bayar pajak motor gimana ya?",
    "name": "Ahmad",
    "location": "Bekasi",
    "latitude": -6.2383,
    "longitude": 106.9756
  }' | jq '.'

echo ""
echo -e "${YELLOW}Tunggu 3 detik...${NC}"
sleep 3

echo -e "${YELLOW}Step 2: User minta bantuan${NC}"
curl -X POST "$BASE_URL/chat" \
  -H "Content-Type: application/json" \
  -d '{
    "session_id": "test-pajak-003",
    "message": "iya tolong bantu saya",
    "name": "Ahmad",
    "location": "Bekasi"
  }' | jq '.'

echo ""
echo -e "${YELLOW}Tunggu 3 detik...${NC}"
sleep 3

# ============================================================================
# TEST 4: PENGESAHAN STNK
# ============================================================================
echo ""
echo -e "${BLUE}════════════════════════════════════════════════════════════════${NC}"
echo -e "${GREEN}TEST 4: Pengesahan STNK (dengan nama: Rina)${NC}"
echo -e "${BLUE}════════════════════════════════════════════════════════════════${NC}"
echo ""

echo -e "${YELLOW}Step 1: User bertanya tentang pengesahan STNK${NC}"
curl -X POST "$BASE_URL/chat" \
  -H "Content-Type: application/json" \
  -d '{
    "session_id": "test-stnk-004",
    "message": "STNK saya mau habis masa berlakunya, harus diapakan?",
    "name": "Rina",
    "location": "Bogor",
    "latitude": -6.5971,
    "longitude": 106.8060
  }' | jq '.'

echo ""
echo -e "${YELLOW}Tunggu 3 detik...${NC}"
sleep 3

# ============================================================================
# TEST 5: MUTASI KENDARAAN
# ============================================================================
echo ""
echo -e "${BLUE}════════════════════════════════════════════════════════════════${NC}"
echo -e "${GREEN}TEST 5: Mutasi Kendaraan (dengan nama: Deni)${NC}"
echo -e "${BLUE}════════════════════════════════════════════════════════════════${NC}"
echo ""

echo -e "${YELLOW}User bertanya tentang mutasi kendaraan${NC}"
curl -X POST "$BASE_URL/chat" \
  -H "Content-Type: application/json" \
  -d '{
    "session_id": "test-mutasi-005",
    "message": "saya mau mutasi motor saya dari Jakarta ke Bandung",
    "name": "Deni",
    "location": "Jakarta Pusat",
    "latitude": -6.1754,
    "longitude": 106.8272
  }' | jq '.'

echo ""
echo -e "${YELLOW}Tunggu 3 detik...${NC}"
sleep 3

# ============================================================================
# TEST 6: BALIK NAMA KENDARAAN
# ============================================================================
echo ""
echo -e "${BLUE}════════════════════════════════════════════════════════════════${NC}"
echo -e "${GREEN}TEST 6: Balik Nama Kendaraan (dengan nama: Fitri)${NC}"
echo -e "${BLUE}════════════════════════════════════════════════════════════════${NC}"
echo ""

echo -e "${YELLOW}User bertanya tentang balik nama${NC}"
curl -X POST "$BASE_URL/chat" \
  -H "Content-Type: application/json" \
  -d '{
    "session_id": "test-balik-nama-006",
    "message": "cara balik nama mobil bekas gimana?",
    "name": "Fitri",
    "location": "Depok",
    "latitude": -6.4025,
    "longitude": 106.7942
  }' | jq '.'

echo ""
echo -e "${YELLOW}Tunggu 3 detik...${NC}"
sleep 3

# ============================================================================
# TEST 7: GANTI PLAT NOMOR
# ============================================================================
echo ""
echo -e "${BLUE}════════════════════════════════════════════════════════════════${NC}"
echo -e "${GREEN}TEST 7: Ganti Plat Nomor (dengan nama: Eko)${NC}"
echo -e "${BLUE}════════════════════════════════════════════════════════════════${NC}"
echo ""

echo -e "${YELLOW}User bertanya tentang ganti plat nomor${NC}"
curl -X POST "$BASE_URL/chat" \
  -H "Content-Type: application/json" \
  -d '{
    "session_id": "test-plat-007",
    "message": "plat nomor motor saya mau habis masa berlakunya 5 tahun",
    "name": "Eko",
    "location": "Tangerang",
    "latitude": -6.1783,
    "longitude": 106.6319
  }' | jq '.'

echo ""
echo -e "${YELLOW}Tunggu 3 detik...${NC}"
sleep 3

# ============================================================================
# TEST 8: CEK E-TILANG + PAJAK (Kombinasi)
# ============================================================================
echo ""
echo -e "${BLUE}════════════════════════════════════════════════════════════════${NC}"
echo -e "${GREEN}TEST 8: Kombinasi - Cek E-Tilang kemudian Bayar Pajak${NC}"
echo -e "${BLUE}════════════════════════════════════════════════════════════════${NC}"
echo ""

echo -e "${YELLOW}Step 1: User cek e-tilang${NC}"
curl -X POST "$BASE_URL/chat" \
  -H "Content-Type: application/json" \
  -d '{
    "session_id": "test-combo-008",
    "message": "cek tilang B1234XYZ",
    "name": "Rudi",
    "location": "Jakarta Barat"
  }' | jq '.'

echo ""
echo -e "${YELLOW}Tunggu 3 detik...${NC}"
sleep 3

echo -e "${YELLOW}Step 2: User mau bayar pajak${NC}"
curl -X POST "$BASE_URL/chat" \
  -H "Content-Type: application/json" \
  -d '{
    "session_id": "test-combo-008",
    "message": "bagaimana cara bayar pajak tahunan motor?",
    "name": "Rudi",
    "location": "Jakarta Barat"
  }' | jq '.'

echo ""
echo -e "${YELLOW}Tunggu 3 detik...${NC}"
sleep 3

# ============================================================================
# TEST 9: TANPA NAMA (fallback ke "Sobat Lantas")
# ============================================================================
echo ""
echo -e "${BLUE}════════════════════════════════════════════════════════════════${NC}"
echo -e "${GREEN}TEST 9: Tanpa Nama - Fallback ke 'Sobat Lantas'${NC}"
echo -e "${BLUE}════════════════════════════════════════════════════════════════${NC}"
echo ""

echo -e "${YELLOW}User bertanya tanpa menyebutkan nama${NC}"
curl -X POST "$BASE_URL/chat" \
  -H "Content-Type: application/json" \
  -d '{
    "session_id": "test-no-name-009",
    "message": "cara perpanjang SIM gimana?",
    "location": "Bandung"
  }' | jq '.'

echo ""
echo -e "${YELLOW}Tunggu 3 detik...${NC}"
sleep 3

# ============================================================================
# TEST 10: UPGRADE SIM
# ============================================================================
echo ""
echo -e "${BLUE}════════════════════════════════════════════════════════════════${NC}"
echo -e "${GREEN}TEST 10: Upgrade SIM (dengan nama: Hendra)${NC}"
echo -e "${BLUE}════════════════════════════════════════════════════════════════${NC}"
echo ""

echo -e "${YELLOW}User bertanya tentang upgrade SIM${NC}"
curl -X POST "$BASE_URL/chat" \
  -H "Content-Type: application/json" \
  -d '{
    "session_id": "test-upgrade-010",
    "message": "saya punya SIM C, mau upgrade jadi SIM A",
    "name": "Hendra",
    "location": "Semarang",
    "latitude": -6.9667,
    "longitude": 110.4167
  }' | jq '.'

echo ""
echo -e "${YELLOW}Tunggu 3 detik...${NC}"
sleep 3

# ============================================================================
# TEST 11: CEK E-TILANG - ADA 2 PELANGGARAN (Data Dummy)
# ============================================================================
echo ""
echo -e "${BLUE}════════════════════════════════════════════════════════════════${NC}"
echo -e "${GREEN}TEST 11: Cek E-Tilang - Ada 2 Pelanggaran (B1234SV)${NC}"
echo -e "${BLUE}════════════════════════════════════════════════════════════════${NC}"
echo ""

echo -e "${YELLOW}User cek tilang dengan plat B1234SV (data dummy: ada 2 pelanggaran)${NC}"
curl -X POST "$BASE_URL/chat" \
  -H "Content-Type: application/json" \
  -d '{
    "session_id": "test-etilang-011",
    "message": "cek tilang B1234SV",
    "name": "Budi Santoso",
    "location": "Jakarta Pusat"
  }' | jq '.'

echo ""
echo -e "${YELLOW}Tunggu 3 detik...${NC}"
sleep 3

# ============================================================================
# TEST 12: CEK E-TILANG - BERSIH (Data Dummy)
# ============================================================================
echo ""
echo -e "${BLUE}════════════════════════════════════════════════════════════════${NC}"
echo -e "${GREEN}TEST 12: Cek E-Tilang - Bersih/Tidak Ada Pelanggaran (B9999ZZ)${NC}"
echo -e "${BLUE}════════════════════════════════════════════════════════════════${NC}"
echo ""

echo -e "${YELLOW}User cek tilang dengan plat B9999ZZ (data dummy: bersih)${NC}"
curl -X POST "$BASE_URL/chat" \
  -H "Content-Type: application/json" \
  -d '{
    "session_id": "test-etilang-012",
    "message": "tolong cek e-tilang untuk B9999ZZ",
    "name": "Ahmad Fauzi",
    "location": "Jakarta Selatan"
  }' | jq '.'

echo ""
echo -e "${YELLOW}Tunggu 3 detik...${NC}"
sleep 3

# ============================================================================
# TEST 13: CEK E-TILANG - 1 PELANGGARAN PARKIR (Data Dummy)
# ============================================================================
echo ""
echo -e "${BLUE}════════════════════════════════════════════════════════════════${NC}"
echo -e "${GREEN}TEST 13: Cek E-Tilang - Pelanggaran Parkir (B5678XY)${NC}"
echo -e "${BLUE}════════════════════════════════════════════════════════════════${NC}"
echo ""

echo -e "${YELLOW}User cek tilang dengan plat B5678XY (data dummy: parkir ilegal)${NC}"
curl -X POST "$BASE_URL/chat" \
  -H "Content-Type: application/json" \
  -d '{
    "session_id": "test-etilang-013",
    "message": "cek tilang mobil B5678XY dong",
    "name": "Siti Rahayu",
    "location": "Jakarta Pusat"
  }' | jq '.'

echo ""
echo -e "${YELLOW}Tunggu 3 detik...${NC}"
sleep 3

# ============================================================================
# TEST 14: CEK E-TILANG - SUDAH DIBAYAR (Data Dummy)
# ============================================================================
echo ""
echo -e "${BLUE}════════════════════════════════════════════════════════════════${NC}"
echo -e "${GREEN}TEST 14: Cek E-Tilang - Pelanggaran Sudah Dibayar (D1111AA)${NC}"
echo -e "${BLUE}════════════════════════════════════════════════════════════════${NC}"
echo ""

echo -e "${YELLOW}User cek tilang dengan plat D1111AA (data dummy: sudah dibayar)${NC}"
curl -X POST "$BASE_URL/chat" \
  -H "Content-Type: application/json" \
  -d '{
    "session_id": "test-etilang-014",
    "message": "cek e-tilang D1111AA",
    "name": "Rina Kartika",
    "location": "Bandung"
  }' | jq '.'

echo ""
echo -e "${YELLOW}Tunggu 3 detik...${NC}"
sleep 3

# ============================================================================
# TEST 15: CEK E-TILANG - PELANGGARAN HP (Data Dummy)
# ============================================================================
echo ""
echo -e "${BLUE}════════════════════════════════════════════════════════════════${NC}"
echo -e "${GREEN}TEST 15: Cek E-Tilang - Menggunakan HP (E7777BB)${NC}"
echo -e "${BLUE}════════════════════════════════════════════════════════════════${NC}"
echo ""

echo -e "${YELLOW}User cek tilang dengan plat E7777BB (data dummy: pakai HP)${NC}"
curl -X POST "$BASE_URL/chat" \
  -H "Content-Type: application/json" \
  -d '{
    "session_id": "test-etilang-015",
    "message": "cek tilang motor E7777BB",
    "name": "Dedi Gunawan",
    "location": "Bandung"
  }' | jq '.'

echo ""
echo -e "${YELLOW}Tunggu 3 detik...${NC}"
sleep 3

# ============================================================================
# TEST 16: CEK E-TILANG - PLAT TIDAK ADA DI DATABASE
# ============================================================================
echo ""
echo -e "${BLUE}════════════════════════════════════════════════════════════════${NC}"
echo -e "${GREEN}TEST 16: Cek E-Tilang - Plat Tidak Terdaftar (B8888XX)${NC}"
echo -e "${BLUE}════════════════════════════════════════════════════════════════${NC}"
echo ""

echo -e "${YELLOW}User cek tilang dengan plat yang tidak ada di database${NC}"
curl -X POST "$BASE_URL/chat" \
  -H "Content-Type: application/json" \
  -d '{
    "session_id": "test-etilang-016",
    "message": "cek tilang B8888XX",
    "name": "Tono",
    "location": "Jakarta"
  }' | jq '.'

echo ""
echo ""
echo -e "${GREEN}════════════════════════════════════════════════════════════════${NC}"
echo -e "${GREEN}✅ SEMUA TEST SELESAI!${NC}"
echo -e "${GREEN}════════════════════════════════════════════════════════════════${NC}"
echo ""
echo "📝 CATATAN:"
echo "   - Setiap test menggunakan session_id yang berbeda"
echo "   - Parameter 'name' digunakan untuk personalisasi"
echo "   - Response harus mengikuti flow dari response-rules.json"
echo "   - Location rules menentukan routing (Satpas/Samsat/Online)"
echo "   - E-Tilang menggunakan data dummy yang sudah ada di etilang.go"
echo ""
echo "🔍 CIRI-CIRI SUKSES:"
echo "   - AI menyapa dengan nama yang disebutkan"
echo "   - Flow percakapan sesuai dengan pertanyaan/response di rules"
echo "   - <name> diganti dengan nama user yang sesuai"
echo "   - <konteks> diganti dengan informasi lokasi"
echo "   - Tidak ada character # atau * dalam response"
echo "   - E-Tilang menampilkan data sesuai dummy (bersih/ada pelanggaran)"
echo ""
echo "🎫 DATA DUMMY E-TILANG:"
echo "   B1234SV - 2 pelanggaran (lampu merah + helm) - Belum dibayar"
echo "   B5678XY - 1 pelanggaran (parkir ilegal) - Belum dibayar"
echo "   B9999ZZ - Bersih, tidak ada pelanggaran"
echo "   D1111AA - 1 pelanggaran (kecepatan) - Sudah dibayar"
echo "   E7777BB - 1 pelanggaran (HP saat berkendara) - Belum dibayar"
echo "   Lainnya - Otomatis bersih (tidak ada di database)"
echo ""
