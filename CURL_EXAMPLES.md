# Sample cURL untuk Chat dengan Upload Dokumen

## 1. Chat Biasa (Tanpa Dokumen)

```bash
curl -X POST http://localhost:3000/api/chat \
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
  }'
```

## 2. Upload 1 Dokumen

```bash
curl -X POST http://localhost:3000/api/chat \
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
        "file_name": "ktp.jpg",
        "file_type": "image/jpeg",
        "file_url": "https://example.com/uploads/ktp.jpg",
        "description": "KTP untuk pengajuan SIM"
      }
    ]
  }'
```

## 3. Upload Multiple Dokumen Sekaligus

```bash
curl -X POST http://localhost:3000/api/chat \
  -H "Content-Type: application/json" \
  -d '{
    "message": "Ini semua dokumen saya untuk perpanjang SIM",
    "context": {
      "location": "Jakarta Selatan"
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
        "description": "Surat Keterangan Sehat"
      }
    ]
  }'
```

## 4. Upload dengan Session ID (Lanjutan Percakapan)

```bash
curl -X POST http://localhost:3000/api/chat \
  -H "Content-Type: application/json" \
  -d '{
    "message": "Ini dokumen tambahannya",
    "session_id": "your-session-id-here",
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
  }'
```

## 5. Upload Berbagai Jenis File (JPG, PNG, PDF)

```bash
curl -X POST http://localhost:3000/api/chat \
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
  }'
```

## Request Body Structure

### ChatRequest
```json
{
  "message": "string (required)",
  "session_id": "string (optional)",
  "context": {
    "location": "string",
    "latitude": 0.0,
    "longitude": 0.0,
    "speed": 0.0,
    "traffic": "string"
  },
  "documents": [
    {
      "file_name": "string (required)",
      "file_type": "string (required)",
      "file_url": "string (required)",
      "description": "string (optional)"
    }
  ]
}
```

### Document Object
- **file_name**: Nama file (contoh: "ktp.jpg", "sim_lama.pdf")
- **file_type**: MIME type (contoh: "image/jpeg", "image/png", "application/pdf")
- **file_url**: URL tempat file disimpan (bisa cloud storage atau server upload)
- **description**: Deskripsi dokumen (opsional, membantu AI memahami konteks)

## Response Format

### Success Response
```json
{
  "success": true,
  "response": "Terima kasih Sobat Lantas! ✅\n\nDokumen yang Anda upload sudah kami terima...",
  "session_id": "abc-123-def-456",
  "sim_flow_info": {
    "active": true,
    "current_node": "collect_ktp",
    "node_type": "collect",
    "node_text": "Silakan upload foto KTP Anda",
    "choices": []
  }
}
```

### Response Fields
- **success**: Boolean, status request
- **response**: String, AI response text
- **session_id**: String, ID untuk melanjutkan percakapan
- **e_tilang_info**: Object (optional), info tilang jika ada
- **pelayanan_info**: Object (optional), info pelayanan jika ditanyakan
- **sim_flow_info**: Object (optional), info flow SIM jika aktif
- **error**: String (jika error), pesan error

## Testing

### Menggunakan Test Script
```bash
# Test semua scenario upload dokumen
./test-upload-document.sh

# Test full SIM flow dengan upload dokumen
./test-sim-flow.sh
```

### Manual Testing dengan curl
```bash
# 1. Start conversation
SESSION_ID=$(curl -s -X POST http://localhost:3000/api/chat \
  -H "Content-Type: application/json" \
  -d '{"message":"Mau perpanjang SIM"}' | jq -r '.session_id')

# 2. Upload document
curl -X POST http://localhost:3000/api/chat \
  -H "Content-Type: application/json" \
  -d "{
    \"message\": \"Ini KTP saya\",
    \"session_id\": \"$SESSION_ID\",
    \"documents\": [{
      \"file_name\": \"ktp.jpg\",
      \"file_type\": \"image/jpeg\",
      \"file_url\": \"https://example.com/ktp.jpg\",
      \"description\": \"KTP\"
    }]
  }" | jq '.'
```

## Tips

1. **File URL**: Dalam implementasi real, file harus di-upload dulu ke storage (S3, GCS, atau server), baru URL-nya dikirim ke API ini
2. **Multiple Files**: Bisa upload multiple files sekaligus dalam 1 request
3. **Session**: Gunakan session_id yang sama untuk melanjutkan percakapan
4. **Context**: Context membantu AI memberikan response yang lebih relevan
5. **Description**: Berikan description yang jelas untuk setiap dokumen agar AI memahami konteks
6. **File Types**: Support image (jpg, png) dan document (pdf)

## Contoh Integration dengan Frontend

### JavaScript/React
```javascript
// Upload file ke server/storage
const uploadFile = async (file) => {
  const formData = new FormData();
  formData.append('file', file);
  
  const response = await fetch('https://your-storage/upload', {
    method: 'POST',
    body: formData
  });
  
  const data = await response.json();
  return data.url; // Return URL file yang sudah di-upload
};

// Kirim chat dengan dokumen
const sendChatWithDocument = async (message, files) => {
  // Upload semua file dulu
  const uploadedDocs = await Promise.all(
    files.map(async (file) => ({
      file_name: file.name,
      file_type: file.type,
      file_url: await uploadFile(file),
      description: file.description || ''
    }))
  );
  
  // Kirim ke chat API
  const response = await fetch('http://localhost:3000/api/chat', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({
      message,
      session_id: sessionId,
      documents: uploadedDocs,
      context: {
        location: currentLocation
      }
    })
  });
  
  return await response.json();
};
```

### Python
```python
import requests

def chat_with_document(message, documents, session_id=None):
    url = "http://localhost:3000/api/chat"
    
    payload = {
        "message": message,
        "documents": documents,
        "context": {
            "location": "Jakarta"
        }
    }
    
    if session_id:
        payload["session_id"] = session_id
    
    response = requests.post(url, json=payload)
    return response.json()

# Usage
documents = [
    {
        "file_name": "ktp.jpg",
        "file_type": "image/jpeg",
        "file_url": "https://example.com/ktp.jpg",
        "description": "KTP"
    }
]

result = chat_with_document("Ini KTP saya", documents)
print(result["response"])
```
