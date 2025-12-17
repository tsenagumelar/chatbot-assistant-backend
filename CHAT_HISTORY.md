# Chat History Implementation Guide

## 🎯 Overview

Backend sekarang mendukung chat history untuk mempertahankan konteks percakapan. AI akan mengingat apa yang sudah dibicarakan sebelumnya.

**PENTING:** Frontend WAJIB mengirimkan field `history` di setiap request agar AI dapat mengingat konteks. Tanpa history, setiap pesan dianggap sebagai percakapan baru.

## 🧪 Testing

Gunakan salah satu cara berikut untuk test:

1. **Test dengan HTML**: Buka file `test_frontend.html` di browser
2. **Test dengan Script**: Jalankan `./test_chat_history.sh` (pastikan backend running)
3. **Test Manual dengan curl**: Lihat contoh di bawah

## API Request Format

### Request Body

```json
{
  "message": "berapa jaraknya dari lokasi saya?",
  "context": {
    "location": "Jakarta Selatan",
    "latitude": -6.2608,
    "longitude": 106.7819,
    "speed": 45.5,
    "traffic": "moderate"
  },
  "history": [
    {
      "role": "user",
      "content": "Saya mau ke Kantor Samsat Tangsel"
    },
    {
      "role": "assistant",
      "content": "Baik, Kantor Samsat Tangsel berlokasi di Jl. Pajajaran No.100, Pamulang..."
    }
  ]
}
```

## Frontend Implementation

### Menyimpan Chat History di Frontend

```javascript
// Inisialisasi array untuk menyimpan history
let chatHistory = [];

// Fungsi untuk mengirim pesan
async function sendMessage(message, context) {
  try {
    const response = await fetch("http://localhost:8080/api/chat", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({
        message: message,
        context: context,
        history: chatHistory, // Kirim history
      }),
    });

    const data = await response.json();

    if (data.success) {
      // Tambahkan pesan user ke history
      chatHistory.push({
        role: "user",
        content: message,
      });

      // Tambahkan respons assistant ke history
      chatHistory.push({
        role: "assistant",
        content: data.response,
      });

      // Optional: Batasi history (misal max 10 pesan terakhir)
      if (chatHistory.length > 10) {
        chatHistory = chatHistory.slice(-10);
      }

      return data.response;
    }
  } catch (error) {
    console.error("Error:", error);
  }
}

// Fungsi untuk reset history (saat mulai percakapan baru)
function resetChatHistory() {
  chatHistory = [];
}
```

## Contoh Skenario

### Skenario 1: Pertanyaan Follow-up

```
User: "Saya mau ke Kantor Samsat Tangsel"
AI: "Baik, Kantor Samsat Tangsel berlokasi di..."

User: "Berapa jaraknya dari lokasi saya?"
AI: "Jarak dari lokasi Anda di Jakarta Selatan ke Kantor Samsat Tangsel sekitar 15 km..."
     ☝️ AI mengingat tujuan dari pesan sebelumnya
```

### Skenario 2: Lanjutan Percakapan

```
User: "Ada rute alternatif ke Samsat Tangsel?"
AI: "Tentu! Untuk perjalanan ke Samsat Tangsel yang tadi kita bicarakan, ada 2 rute alternatif..."
     ☝️ AI mengingat konteks lokasi tujuan
```

## Best Practices

1. **Batasi Jumlah History**: Simpan maksimal 10-20 pesan terakhir untuk menghindari token yang terlalu banyak
2. **Reset Saat Diperlukan**: Beri opsi untuk user memulai percakapan baru
3. **Simpan di Local Storage**: Untuk persistensi saat refresh halaman (optional)
4. **Update Context**: Selalu kirim context terbaru (lokasi, speed, traffic) di setiap request

## Notes

- System prompt sudah dioptimasi untuk memahami dan menggunakan konteks dari history
- AI diinstruksikan untuk TIDAK meminta informasi yang sudah disebutkan sebelumnya
- Log akan menampilkan jumlah pesan history yang dikirim: `📚 Including X messages from history`
