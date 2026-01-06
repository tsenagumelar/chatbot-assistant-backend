# Police Assistant Backend - Dokumentasi Lengkap

## Daftar Isi

1. [Tentang Project](#tentang-project)
2. [Chat History & Context Management](#chat-history--context-management)
3. [Session Management](#session-management)
4. [API Documentation](#api-documentation)
5. [Frontend Implementation](#frontend-implementation)
6. [Testing Guide](#testing-guide)
7. [Troubleshooting](#troubleshooting)

---

## Tentang Project

Backend untuk Police Assistant - asisten virtual berbasis AI untuk membantu pengguna dengan informasi terkait layanan kepolisian, navigasi, dan layanan publik.

### Fitur Utama

- **Chat dengan AI**: Menggunakan OpenAI GPT-4o untuk memberikan response yang kontekstual
- **Chat History**: AI mengingat konteks percakapan sebelumnya
- **Session Management**: Backend otomatis mengelola history chat per session
- **Location Context**: Mendukung informasi lokasi, kecepatan, dan kondisi lalu lintas
- **Auto Cleanup**: Session otomatis dibersihkan setelah 24 jam tidak aktif
- **E-Tilang Check**: Cek pelanggaran e-tilang berdasarkan nomor polisi
- **Pelayanan Info**: Informasi pelayanan polisi dan dokumen yang diperlukan
- **Document Upload**: Dukungan upload dokumen untuk berbagai pelayanan
- **🆕 SIM Flow**: Alur percakapan terstruktur untuk perpanjangan/pembuatan SIM (lihat [SIM_FLOW.md](SIM_FLOW.md))

---

## Chat History & Context Management

### Masalah yang Diselesaikan

Sebelumnya, AI tidak mengingat konteks percakapan:

```
User: "Nama saya Taufan"
AI: "Halo Taufan!"

User: "Siapa nama saya?"
AI: "Saya tidak memiliki informasi nama Anda" ❌
```

Sekarang dengan chat history:

```
User: "Nama saya Taufan"
AI: "Halo Taufan!"

User: "Siapa nama saya?"
AI: "Nama Anda Taufan!" ✅
```

### Cara Kerja

Backend sekarang menyimpan history chat menggunakan **session management**. Frontend hanya perlu:

1. Kirim pesan dengan atau tanpa `session_id`
2. Backend otomatis create session jika belum ada
3. Backend return `session_id` di response
4. Frontend simpan dan kirim `session_id` di request berikutnya

### Keuntungan Session Management

1. **Lebih Simple**: Frontend tidak perlu track history sendiri
2. **Konsisten**: History disimpan terpusat di backend
3. **Aman**: History otomatis dibersihkan setelah 24 jam tidak aktif
4. **Efficient**: Otomatis batasi max 30 pesan per session

---

## Session Management

### Flow Sederhana

```
1. Frontend kirim pesan TANPA session_id
   → Backend otomatis buat session baru
   → Return response + session_id

2. Frontend simpan session_id, kirim pesan berikutnya DENGAN session_id
   → Backend ambil history dari session
   → AI ingat konteks
   → Return response + session_id yang sama

3. Frontend terus pakai session_id yang sama untuk percakapan
   → AI selalu ingat konteks
```

### Storage Info

- **Type**: In-memory (sync.Map)
- **Persistence**: Session hilang saat server restart
- **Cleanup**: Auto delete session tidak aktif > 24 jam
- **Max Messages**: Otomatis limit 30 pesan per session
- **Future**: Bisa upgrade ke Redis/Database jika perlu persistence

---

## API Documentation

### Base URL

```
http://localhost:8080/api/v1
```

### 1. Chat Endpoint (Main)

**Endpoint**: `POST /api/v1/chat`

**Request Body**:

```json
{
  "message": "Halo, nama saya Taufan",
  "session_id": "",  // Kosong untuk session baru, atau kirim session_id yang ada
  "context": {
    "location": "Jakarta Selatan",
    "latitude": -6.2608,
    "longitude": 106.7819,
    "speed": 0,
    "traffic": "smooth"
  }
}
```

**Response**:

```json
{
  "success": true,
  "response": "Halo Taufan! Senang berkenalan dengan Anda...",
  "session_id": "550e8400-e29b-41d4-a716-446655440000"
}
```

**Catatan**:
- Jika `session_id` kosong, backend otomatis create session baru
- Simpan `session_id` dari response untuk request berikutnya
- Backend otomatis manage chat history berdasarkan session

### 2. Create Session (Optional)

**Endpoint**: `POST /api/v1/session`

**Response**:

```json
{
  "success": true,
  "session_id": "550e8400-e29b-41d4-a716-446655440000",
  "message": "Session created successfully"
}
```

**Catatan**: Endpoint ini optional karena chat endpoint sudah otomatis create session jika belum ada.

### 3. Clear Session History

**Endpoint**: `POST /api/v1/session/{session_id}/clear`

**Response**:

```json
{
  "success": true,
  "session_id": "550e8400-e29b-41d4-a716-446655440000",
  "message": "Session history cleared"
}
```

**Use Case**: Untuk reset chat history tanpa create session baru.

### 4. Delete Session

**Endpoint**: `DELETE /api/v1/session/{session_id}`

**Response**:

```json
{
  "success": true,
  "message": "Session deleted successfully"
}
```

**Use Case**: Untuk menghapus session sepenuhnya.

### 5. Get Session Info (Debug)

**Endpoint**: `GET /api/v1/session/{session_id}`

**Response**:

```json
{
  "success": true,
  "session_id": "550e8400-e29b-41d4-a716-446655440000",
  "message_count": 4,
  "created_at": "2025-12-16T17:00:00Z",
  "updated_at": "2025-12-16T17:05:00Z"
}
```

**Use Case**: Untuk debugging, cek jumlah pesan dalam session.

---

## Frontend Implementation

### Implementasi Sederhana (Vanilla JS)

```javascript
// Simpan session ID di state/memory
let currentSessionId = null;

async function sendMessage(message, context) {
  try {
    const response = await fetch("http://localhost:8080/api/v1/chat", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({
        message: message,
        session_id: currentSessionId || "",  // Kosong jika belum ada
        context: context,
      }),
    });

    const data = await response.json();

    if (data.success) {
      // ⭐ SIMPAN session_id dari response
      currentSessionId = data.session_id;

      return data.response;
    }
  } catch (error) {
    console.error("Error:", error);
  }
}

// Reset chat = hapus session_id lokal atau call clear API
function resetChat() {
  // Option 1: Set null untuk create session baru
  currentSessionId = null;

  // Option 2: Clear history tapi keep session
  // fetch(`http://localhost:8080/api/v1/session/${currentSessionId}/clear`, {method: 'POST'})
}
```

### Implementasi dengan React

```javascript
import { useState } from 'react';

function ChatApp() {
  const [sessionId, setSessionId] = useState(null);
  const [messages, setMessages] = useState([]);

  const sendMessage = async (message, context) => {
    try {
      const response = await fetch("http://localhost:8080/api/v1/chat", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({
          message: message,
          session_id: sessionId || "",
          context: context,
        }),
      });

      const data = await response.json();

      if (data.success) {
        // Simpan session ID jika belum ada
        if (!sessionId) {
          setSessionId(data.session_id);
        }

        // Update UI
        setMessages([
          ...messages,
          { role: "user", content: message },
          { role: "assistant", content: data.response },
        ]);

        return data.response;
      }
    } catch (error) {
      console.error("Error:", error);
    }
  };

  const resetChat = () => {
    setSessionId(null);
    setMessages([]);
  };

  return (
    <div>
      {/* UI components */}
    </div>
  );
}
```

### Best Practices

1. **Simpan Session ID**: Selalu simpan `session_id` dari response pertama
2. **Reset Chat**: Beri opsi user untuk mulai percakapan baru dengan set `session_id = null`
3. **Error Handling**: Handle error jika session expired atau tidak valid
4. **LocalStorage (Optional)**: Simpan `session_id` di localStorage untuk persist saat refresh
5. **Context Update**: Selalu kirim context terbaru (lokasi, speed, traffic) di setiap request

---

## Testing Guide

### Test 1: Basic Chat Flow

```bash
# Pesan pertama - TANPA session_id
curl -X POST http://localhost:8080/api/v1/chat \
  -H "Content-Type: application/json" \
  -d '{
    "message": "Halo, nama saya Taufan",
    "session_id": "",
    "context": {
      "location": "Jakarta",
      "latitude": -6.2,
      "longitude": 106.8,
      "speed": 0,
      "traffic": "smooth"
    }
  }'

# Response akan include session_id:
# {"success":true,"response":"...","session_id":"550e8400-..."}
```

### Test 2: Context Memory

```bash
# Pesan kedua - DENGAN session_id dari response pertama
curl -X POST http://localhost:8080/api/v1/chat \
  -H "Content-Type: application/json" \
  -d '{
    "message": "Siapa nama saya?",
    "session_id": "PASTE_SESSION_ID_DARI_RESPONSE_PERTAMA",
    "context": {
      "location": "Jakarta",
      "latitude": -6.2,
      "longitude": 106.8,
      "speed": 0,
      "traffic": "smooth"
    }
  }'

# AI akan jawab "Taufan" ✅
```

### Test 3: Follow-up Questions

Skenario untuk test apakah AI ingat konteks:

```
1. "Saya mau ke Kantor Samsat Tangsel"
   → AI akan kasih info lokasi Samsat Tangsel

2. "Berapa jaraknya dari lokasi saya?"
   → AI harus ingat tujuan (Samsat Tangsel) dari pesan sebelumnya ✅

3. "Ada rute alternatif?"
   → AI masih ingat tujuan ✅
```

### Test 4: Clear History

```bash
curl -X POST http://localhost:8080/api/v1/session/YOUR_SESSION_ID/clear
```

Setelah clear, AI tidak akan ingat percakapan sebelumnya.

### Test 5: Get Session Info

```bash
curl http://localhost:8080/api/v1/session/YOUR_SESSION_ID
```

### Test 6: Delete Session

```bash
curl -X DELETE http://localhost:8080/api/v1/session/YOUR_SESSION_ID
```

### Test dengan HTML Page

Lihat file test yang tersedia:
- `test_frontend.html` - Test page dengan UI (jika ada)
- `test_session_frontend.html` - Test page untuk session management (jika ada)

---

## Troubleshooting

### 1. AI Masih Lupa Konteks

**Kemungkinan Penyebab**:
- Frontend tidak mengirim `session_id`
- Session sudah expired (> 24 jam tidak aktif)
- `session_id` salah atau tidak valid

**Solusi**:
```javascript
// Pastikan session_id disimpan dan dikirim
console.log("Session ID:", currentSessionId);  // Harus ada value

// Cek response dari backend
const data = await response.json();
console.log("Session ID dari backend:", data.session_id);
```

### 2. Error "Session Not Found"

**Penyebab**: Session sudah dihapus (expired atau di-delete)

**Solusi**:
```javascript
// Reset session_id untuk create session baru
currentSessionId = null;
```

### 3. Backend Tidak Respond

**Checklist**:
- [ ] Backend sudah running? (`go run main.go`)
- [ ] Port 8080 tidak dipakai aplikasi lain?
- [ ] CORS sudah dikonfigurasi dengan benar?
- [ ] API key Claude sudah diset?

### 4. Session Hilang Setelah Server Restart

**Ini Normal**: Session disimpan di memory (in-memory storage).

**Solusi**:
- Untuk production: upgrade ke Redis atau Database
- Untuk development: session akan auto-create ulang saat user chat

### 5. Cek Backend Logs

Backend akan log informasi penting:

```
📚 Including X messages from history    // History diterima dengan benar
🆕 Creating new session for user        // Session baru dibuat
♻️  Using existing session: xxx         // Menggunakan session yang ada
🧹 Cleaned up X expired sessions        // Auto cleanup berjalan
```

---

## Key Points

### Yang WAJIB di Frontend:

1. ✅ Simpan `session_id` dari response
2. ✅ Kirim `session_id` di request berikutnya
3. ✅ Set `session_id: ""` untuk session baru
4. ✅ Kirim `context` yang updated di setiap request

### Yang OTOMATIS di Backend:

1. ✅ Create session jika belum ada
2. ✅ Simpan chat history per session
3. ✅ Limit max 30 pesan per session
4. ✅ Clean up session expired (> 24 jam)

### Perbandingan: Before vs After

**Before (Manual History - Kompleks)**:
```javascript
// ❌ Frontend harus manage history
let chatHistory = [];

fetch("/api/chat", {
  body: JSON.stringify({
    message: msg,
    history: chatHistory,  // Frontend track semua pesan
  }),
});

// Frontend harus push manual
chatHistory.push({role: "user", content: msg});
chatHistory.push({role: "assistant", content: resp});
```

**After (Session-based - Simple)**:
```javascript
// ✅ Backend yang manage, frontend hanya kirim ID
let sessionId = null;

fetch("/api/v1/chat", {
  body: JSON.stringify({
    message: msg,
    session_id: sessionId,  // Hanya kirim ID
  }),
});

// Backend otomatis simpan, frontend hanya simpan ID
sessionId = response.session_id;
```

---

## Status & Version

- **Status**: ✅ Production Ready
- **Version**: 2.0 (Session-based Management)
- **Last Update**: 16 Desember 2025
- **Backend**: Ready
- **Frontend**: Perlu update untuk support session management

---

## Contact & Support

Untuk pertanyaan atau issue, silakan hubungi tim development.

**Dokumentasi ini menggabungkan**:
- README_CONTEXT_FIX.md - Informasi tentang perbaikan context issue
- CHAT_HISTORY.md - Dokumentasi chat history implementation
- FIX_CONTEXT_ISSUE.md - Panduan fix untuk frontend
- SESSION_MANAGEMENT.md - Dokumentasi session management system
