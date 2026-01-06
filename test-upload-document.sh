#!/bin/bash

# Colors for output
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

API_URL="http://localhost:8080/api/v1/chat"

echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo -e "${BLUE}  📤 TEST CHAT WITH DOCUMENT UPLOAD${NC}"
echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo ""

# Test 1: Chat biasa tanpa dokumen
echo -e "${GREEN}Test 1: Chat biasa tanpa dokumen${NC}"
echo -e "${YELLOW}Request: \"Saya mau tanya tentang pembuatan SIM C\"${NC}"
curl -X POST "$API_URL" \
  -H "Content-Type: application/json" \
  -d '{
    "message": "Saya mau tanya tentang pembuatan SIM C",
    "context": {
      "location": "Jakarta Selatan",
      "latitude": -6.2608,
      "longitude": 106.7818,
      "speed": 0,
      "traffic": "normal"
    }
  }' | jq '.'

echo -e "\n${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}\n"
sleep 2

# Test 2: Upload 1 dokumen
echo -e "${GREEN}Test 2: Upload 1 dokumen (KTP)${NC}"
echo -e "${YELLOW}Request: Chat dengan upload 1 file${NC}"
curl -X POST "http://localhost:8080/api/v1/chat" \
  -H "Content-Type: application/json" \
  -d '{
    "message": "Ini KTP saya untuk pengajuan SIM",
    "context": {
      "location": "Jakarta Selatan",
      "latitude": -6.2608,
      "longitude": 106.7818
    },
    "documents": [
      {
        "file_name": "ktp_taufan.jpg",
        "file_type": "image/jpeg",
        "file_url": "https://example.com/uploads/ktp_taufan.jpg",
        "description": "KTP untuk pengajuan SIM"
      }
    ]
  }' | jq '.'

echo -e "\n${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}\n"
sleep 2

# Test 3: Upload multiple dokumen sekaligus
echo -e "${GREEN}Test 3: Upload multiple dokumen (KTP + SIM Lama + Surat Sehat)${NC}"
echo -e "${YELLOW}Request: Chat dengan upload 3 file sekaligus${NC}"
curl -X POST "$API_URL" \
  -H "Content-Type: application/json" \
  -d '{
    "message": "Ini semua dokumen yang saya punya untuk perpanjang SIM",
    "context": {
      "location": "Jakarta Selatan",
      "latitude": -6.2608,
      "longitude": 106.7818
    },
    "documents": [
      {
        "file_name": "ktp.jpg",
        "file_type": "image/jpeg",
        "file_url": "https://example.com/uploads/ktp.jpg",
        "description": "KTP Asli"
      },
      {
        "file_name": "sim_lama.jpg",
        "file_type": "image/jpeg",
        "file_url": "https://example.com/uploads/sim_lama.jpg",
        "description": "SIM A yang mau diperpanjang"
      },
      {
        "file_name": "surat_sehat.pdf",
        "file_type": "application/pdf",
        "file_url": "https://example.com/uploads/surat_sehat.pdf",
        "description": "Surat Keterangan Sehat dari dokter"
      }
    ]
  }' | jq '.'

echo -e "\n${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}\n"
sleep 2

# Test 4: Upload dengan session_id (lanjutan percakapan)
echo -e "${GREEN}Test 4: Upload dokumen dengan session_id${NC}"
echo -e "${YELLOW}Request: Upload file dalam session yang sudah ada${NC}"
echo -e "${YELLOW}Note: Ganti <SESSION_ID> dengan session_id dari response sebelumnya${NC}"
echo ""
echo 'curl -X POST "'"$API_URL"'" \'
echo '  -H "Content-Type: application/json" \'
echo '  -d '"'"'{
    "message": "Ini dokumen tambahannya",
    "session_id": "<SESSION_ID>",
    "context": {
      "location": "Jakarta Selatan"
    },
    "documents": [
      {
        "file_name": "tes_psikologi.jpg",
        "file_type": "image/jpeg",
        "file_url": "https://example.com/uploads/tes_psikologi.jpg",
        "description": "Hasil tes psikologi"
      }
    ]
  }'"'"' | jq '"'"'.'"'"

echo -e "\n${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}\n"

# Test 5: Upload berbagai jenis file
echo -e "${GREEN}Test 5: Upload berbagai jenis file (JPG, PNG, PDF)${NC}"
echo -e "${YELLOW}Request: Upload file dengan berbagai format${NC}"
curl -X POST "$API_URL" \
  -H "Content-Type: application/json" \
  -d '{
    "message": "Saya upload dokumen untuk SKCK",
    "context": {
      "location": "Jakarta Selatan"
    },
    "documents": [
      {
        "file_name": "ktp.jpg",
        "file_type": "image/jpeg",
        "file_url": "https://example.com/uploads/ktp.jpg",
        "description": "KTP"
      },
      {
        "file_name": "kk.png",
        "file_type": "image/png",
        "file_url": "https://example.com/uploads/kk.png",
        "description": "Kartu Keluarga"
      },
      {
        "file_name": "akta_lahir.pdf",
        "file_type": "application/pdf",
        "file_url": "https://example.com/uploads/akta_lahir.pdf",
        "description": "Akta Kelahiran"
      }
    ]
  }' | jq '.'

echo -e "\n${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo -e "${GREEN}✅ All tests completed!${NC}"
echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo ""
echo -e "${YELLOW}Tips:${NC}"
echo "- Response akan include: success, response, session_id"
echo "- Jika upload dokumen, context.HasUploadedDocuments akan true"
echo "- AI akan memberikan konfirmasi dokumen diterima dengan emoji ✅"
echo "- Gunakan session_id untuk melanjutkan percakapan"
echo ""
