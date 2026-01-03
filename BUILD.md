# Docker Build & Push Guide

## Cara Build dan Push ke Docker Hub

### 1. Login ke Docker Hub
```bash
docker login
```
Masukkan username dan password Docker Hub Anda.

### 2. Build Docker Image

**Untuk server Linux AMD64 (paling umum):**
```bash
docker build --platform linux/amd64 -t taufansena/chatbot-backend:latest .
```

**Atau tanpa platform flag (sudah di-set di Dockerfile):**
```bash
docker build -t taufansena/chatbot-backend:latest .
```

### 3. Push ke Docker Hub
```bash
docker push taufansena/chatbot-backend:latest
```

### 4. (Optional) Build dengan Tag Versi Spesifik
```bash
# Build dengan tag versi
docker build -t taufansena/chatbot-backend:v1.0.0 -t taufansena/chatbot-backend:latest .

# Push semua tags
docker push taufansena/chatbot-backend:v1.0.0
docker push taufansena/chatbot-backend:latest
```

---

## Environment Variables yang Perlu di-Set di Portainer

Saat deploy di Portainer, set environment variables berikut:

```env
PORT=8080
OPENAI_API_KEY=your_openai_api_key_here
ORS_API_KEY=your_openroute_service_api_key_here
```

### Contoh Konfigurasi di Portainer:

1. **Port Mapping**: `8080:8080`
2. **Environment Variables**:
   - `PORT` = `8080`
   - `OPENAI_API_KEY` = `sk-proj-xxxxx...`
   - `ORS_API_KEY` = `5b3ce3597851110001cf6248xxxxx...`

---

## Test Lokal Sebelum Push

### 1. Build image lokal
```bash
docker build -t chatbot-backend:test .
```

### 2. Jalankan container dengan env variables
```bash
docker run -d \
  -p 8080:8080 \
  -e PORT=8080 \
  -e OPENAI_API_KEY=your_key \
  -e ORS_API_KEY=your_key \
  --name chatbot-test \
  chatbot-backend:test
```

### 3. Test API
```bash
# Health check
curl http://localhost:8080/health

# Root endpoint
curl http://localhost:8080/
```

### 4. Lihat logs
```bash
docker logs -f chatbot-test
```

### 5. Stop dan hapus test container
```bash
docker stop chatbot-test
docker rm chatbot-test
```

---

## Quick Commands

```bash
# Build & Push (satu kali jalan) - untuk Linux AMD64
docker build --platform linux/amd64 -t taufansena/chatbot-backend:latest . && \
docker push taufansena/chatbot-backend:latest

# Pull di server (setelah push)
docker pull taufansena/chatbot-backend:latest
```

---

## Troubleshooting

### Image terlalu besar?
Dockerfile sudah menggunakan multi-stage build dengan Alpine Linux untuk hasil yang minimal (~20-30MB).

### Build gagal?
```bash
# Clear cache dan build ulang
docker build --no-cache -t taufansena/chatbot-backend:latest .
```

### Lupa login Docker Hub?
```bash
docker login
# Username: taufansena
# Password: [your Docker Hub password/token]
```

### Error "exec format error" di server?
Ini terjadi karena arsitektur binary tidak cocok dengan server. **SOLUSI:**

```bash
# Rebuild dengan platform yang benar (Linux AMD64)
docker build --platform linux/amd64 -t taufansena/chatbot-backend:latest .
docker push taufansena/chatbot-backend:latest

# Di Portainer: Pull image ulang dan recreate container
```

Dockerfile sudah di-update untuk build GOARCH=amd64, jadi rebuild dan push ulang!
