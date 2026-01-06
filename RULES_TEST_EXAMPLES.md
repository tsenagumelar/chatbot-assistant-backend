# Test Examples untuk Response Rules & Location Rules

Script test komprehensif tersedia di: `test-rules.sh`

## Cara Menjalankan

```bash
# Pastikan server sudah running
go run main.go

# Di terminal lain, jalankan test script
./test-rules.sh
```

## Individual Test Cases

### 1. Buat SIM Baru (dengan nama)

```bash
# Step 1: User bertanya
curl -X POST "http://localhost:3000/api/v1/chat" \
  -H "Content-Type: application/json" \
  -d '{
    "session_id": "test-sim-001",
    "message": "mau bikin sim baru",
    "name": "Budi",
    "location": "Jakarta Selatan"
  }'

# Expected: AI menyapa "Budi" dan mengikuti flow dari response-rules.json
# - Pertanyaan 1: terkait dokumen
# - Response 1: harus menyebutkan nama "Budi" dan lokasi
```

### 2. Perpanjangan SIM (dengan upload dokumen)

```bash
# Step 1: Bertanya tentang perpanjangan
curl -X POST "http://localhost:3000/api/v1/chat" \
  -H "Content-Type: application/json" \
  -d '{
    "session_id": "test-sim-002",
    "message": "cara perpanjang SIM?",
    "name": "Siti",
    "location": "Tangerang Selatan"
  }'

# Step 2: Upload dokumen
curl -X POST "http://localhost:3000/api/v1/chat" \
  -H "Content-Type: multipart/form-data" \
  -F "session_id=test-sim-002" \
  -F "message=ini dokumen saya" \
  -F "name=Siti" \
  -F "documents=@test-ktp.jpg"

# Expected: 
# - Turn 1: AI bertanya dokumen
# - Turn 2: User upload, AI konfirmasi dengan ✅
# - Turn 3: AI lanjut ke pertanyaan berikutnya
```

### 3. Pajak Kendaraan Bermotor

```bash
curl -X POST "http://localhost:3000/api/v1/chat" \
  -H "Content-Type: application/json" \
  -d '{
    "session_id": "test-pajak-003",
    "message": "cara bayar pajak motor?",
    "name": "Ahmad",
    "location": "Bekasi"
  }'

# Expected:
# - AI menyapa "Ahmad"
# - Informasi dokumen: STNK asli, KTP asli, BPKB (jika ada)
# - Location rule: Arahkan ke Samsat terdekat
```

### 4. Pengesahan STNK

```bash
curl -X POST "http://localhost:3000/api/v1/chat" \
  -H "Content-Type: application/json" \
  -d '{
    "session_id": "test-stnk-004",
    "message": "STNK mau habis masa berlakunya",
    "name": "Rina",
    "location": "Bogor"
  }'

# Expected:
# - AI menyapa "Rina"
# - Dokumen: STNK asli, KTP asli, BPKB
# - Informasi biaya dan waktu proses
```

### 5. Mutasi Kendaraan

```bash
curl -X POST "http://localhost:3000/api/v1/chat" \
  -H "Content-Type: application/json" \
  -d '{
    "session_id": "test-mutasi-005",
    "message": "mau mutasi motor dari Jakarta ke Bandung",
    "name": "Deni",
    "location": "Jakarta Pusat"
  }'

# Expected:
# - AI jelaskan proses mutasi antar daerah
# - Dokumen lengkap yang diperlukan
# - Waktu proses dan biaya
```

### 6. Balik Nama Kendaraan

```bash
curl -X POST "http://localhost:3000/api/v1/chat" \
  -H "Content-Type: application/json" \
  -d '{
    "session_id": "test-balik-nama-006",
    "message": "cara balik nama mobil bekas?",
    "name": "Fitri",
    "location": "Depok"
  }'

# Expected:
# - Dokumen penjual dan pembeli
# - Surat pernyataan fiskal
# - Cek fisik kendaraan
```

### 7. Ganti Plat Nomor (5 Tahunan)

```bash
curl -X POST "http://localhost:3000/api/v1/chat" \
  -H "Content-Type: application/json" \
  -d '{
    "session_id": "test-plat-007",
    "message": "plat nomor motor mau habis 5 tahun",
    "name": "Eko",
    "location": "Tangerang"
  }'

# Expected:
# - Info tentang ganti plat 5 tahunan
# - Dokumen yang diperlukan
# - Biaya dan waktu proses
```

### 8. Upgrade SIM (C ke A)

```bash
curl -X POST "http://localhost:3000/api/v1/chat" \
  -H "Content-Type: application/json" \
  -d '{
    "session_id": "test-upgrade-008",
    "message": "punya SIM C, mau upgrade jadi SIM A",
    "name": "Hendra",
    "location": "Semarang"
  }'

# Expected:
# - Info syarat upgrade SIM
# - Test dan training yang diperlukan
# - Dokumen pendukung
```

### 9. Tanpa Nama (Fallback)

```bash
curl -X POST "http://localhost:3000/api/v1/chat" \
  -H "Content-Type: application/json" \
  -d '{
    "session_id": "test-no-name-009",
    "message": "cara perpanjang SIM?",
    "location": "Bandung"
  }'

# Expected:
# - AI menyapa dengan "Sobat Lantas" (bukan nama)
# - Flow tetap sesuai response-rules.json
```

### 10. Kombinasi: E-Tilang + Pajak

```bash
# Step 1: Cek e-tilang
curl -X POST "http://localhost:3000/api/v1/chat" \
  -H "Content-Type: application/json" \
  -d '{
    "session_id": "test-combo-010",
    "message": "cek tilang B1234XYZ",
    "name": "Rudi",
    "location": "Jakarta Barat"
  }'

# Step 2: Tanya pajak
curl -X POST "http://localhost:3000/api/v1/chat" \
  -H "Content-Type: application/json" \
  -d '{
    "session_id": "test-combo-010",
    "message": "cara bayar pajak motor?",
    "name": "Rudi",
    "location": "Jakarta Barat"
  }'

# Expected:
# - Step 1: Info e-tilang (bersih/ada pelanggaran)
# - Step 2: AI ingat nama "Rudi" dan lanjut dengan flow pajak
# - Context tetap terjaga antar pesan
```

## Validasi Response

Response yang benar harus memenuhi:

### ✅ Personalisasi
- Jika `name` diberikan, AI harus menyapa dengan nama tersebut
- Placeholder `<name>` di rules diganti dengan nama user
- Jika tanpa nama, fallback ke "Sobat Lantas"

### ✅ Context Replacement
- `<konteks>` diganti dengan informasi lokasi user
- Contoh: "Saya lihat Anda sedang di Jakarta Selatan"

### ✅ Formatting
- TIDAK ada character `#` untuk heading
- TIDAK ada `*` atau `**` untuk bold/italic
- Gunakan emoji dan huruf kapital untuk penekanan
- List menggunakan angka (1. 2. 3.) bukan asterisk

### ✅ Flow Conversation
- Ikuti urutan pertanyaan/response dari response-rules.json
- Jangan skip turn
- Sesuaikan dengan turn yang sedang berjalan

### ✅ Location Rules
- Berikan arahan sesuai location-rules.json
- Satpas SIM untuk urusan SIM
- Samsat untuk urusan kendaraan
- Online untuk layanan digital

## Monitoring Response

Perhatikan hal-hal berikut saat testing:

1. **Nama User**: Apakah AI menyapa dengan nama yang benar?
2. **Context**: Apakah `<name>` dan `<konteks>` sudah diganti?
3. **Flow**: Apakah mengikuti pertanyaan/response sesuai rules?
4. **Format**: Apakah tidak ada Markdown formatting (#, *, **)?
5. **Location**: Apakah arahan lokasi sesuai dengan rules?
6. **Document Upload**: Apakah AI konfirmasi dengan ✅?
7. **Multi-turn**: Apakah AI ingat context dari pesan sebelumnya?

## Debug Tips

Jika response tidak sesuai:

1. Cek log server untuk melihat rule yang dimatch
2. Pastikan `jenisPelayanan` di rules match dengan flow title
3. Verify format JSON di response-rules.json dan location-rules.json
4. Cek apakah RulesService berhasil load file JSON
5. Lihat system prompt yang digenerate (add logging jika perlu)

## File Dummy untuk Test Upload

Buat dummy files untuk testing upload:

```bash
# Create dummy image files
echo "Dummy KTP" > test-ktp.jpg
echo "Dummy SIM Lama" > test-sim-lama.jpg
echo "Dummy STNK" > test-stnk.jpg
echo "Dummy BPKB" > test-bpkb.jpg
```
