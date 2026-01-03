# Deploy ke Portainer - Panduan Lengkap

## Metode 1: Deploy via Portainer Stacks (RECOMMENDED)

### Langkah-langkah:

1. **Login ke Portainer** (biasanya http://your-server-ip:9000)

2. **Pilih Environment** (Local/Docker)

3. **Masuk ke Menu "Stacks"**

4. **Klik "Add Stack"**

5. **Isi Form:**
   - **Name**: `police-assistant-backend`
   - **Build method**: Pilih `Web editor`
   - **Web editor**: Copy paste isi file `portainer-stack.yml`

6. **Set Environment Variables** (scroll ke bawah):

   Klik "Add environment variable" untuk setiap variable:

   | Name | Value |
   |------|-------|
   | `OPENAI_API_KEY` | `sk-proj-xxxxx...` |
   | `ORS_API_KEY` | `5b3ce3597851110001cf6248xxxxx...` |

7. **Deploy Stack** - Klik "Deploy the stack"

8. **Verify** - Cek logs dan akses:
   - Logs: Klik nama stack → klik container → logs
   - Test: `http://your-server-ip:8080/health`

---

## Metode 2: Deploy via Container (Manual)

### Langkah-langkah:

1. **Masuk ke Menu "Containers"**

2. **Klik "Add container"**

3. **Isi Form:**

   **General Settings:**
   - **Name**: `police-assistant-backend`
   - **Image**: `taufansena/chatbot-backend:latest`

   **Network ports configuration:**
   - **+publish a new network port**
   - Host: `8080` → Container: `8080`

   **Env (Environment variables):**
   - Klik "+ add environment variable"
   - `PORT` = `8080`
   - `OPENAI_API_KEY` = `sk-proj-xxxxx...`
   - `ORS_API_KEY` = `5b3ce3597851110001cf6248xxxxx...`

   **Restart policy:**
   - Pilih: `Unless stopped`

4. **Deploy container** - Klik "Deploy the container"

5. **Verify** - Lihat logs dan test endpoint

---

## Update Image (Pull Latest Version)

### Via Stacks:
1. Masuk ke Stack → Klik stack name
2. Scroll ke container
3. Klik "Pull latest image" icon (↓)
4. Klik "Recreate" container

### Via Container:
1. Masuk ke Containers
2. Stop container
3. Pull latest image:
   - Klik container → "Duplicate/Edit"
   - Checklist "Always pull the image"
   - Recreate container

---

## Monitoring & Logs

### Lihat Logs:
1. Masuk ke **Containers**
2. Klik nama container `police-assistant-backend`
3. Klik tab **"Logs"**
4. Enable "Auto-refresh logs"

### Health Check:
- Portainer akan otomatis monitor via health check endpoint
- Status "healthy" akan muncul di container list

### Test Endpoint:
```bash
# Health check
curl http://your-server-ip:8080/health

# Root endpoint
curl http://your-server-ip:8080/

# Chat test (gunakan Postman atau curl)
curl -X POST http://your-server-ip:8080/api/v1/chat \
  -H "Content-Type: application/json" \
  -d '{
    "message": "Hello",
    "session_id": "test-123"
  }'
```

---

## Troubleshooting

### Container tidak start?
- Cek logs untuk error message
- Pastikan port 8080 tidak dipakai service lain
- Verify environment variables sudah di-set dengan benar

### Image tidak bisa pull?
- Pastikan image sudah di-push ke Docker Hub
- Cek nama image: `taufansena/chatbot-backend:latest`
- Pastikan image bersifat public di Docker Hub

### API Key tidak bekerja?
- Cek environment variables di container settings
- Pastikan tidak ada spasi atau karakter aneh
- Test API key langsung di OpenAI dashboard

### Port tidak accessible?
- Cek firewall server
- Pastikan port mapping benar: `8080:8080`
- Test dari server: `curl localhost:8080/health`

---

## Auto-Update dengan Watchtower (Optional)

Untuk auto-update image saat ada versi baru:

```yaml
# Tambahkan service ini di stack
watchtower:
  image: containrrr/watchtower
  volumes:
    - /var/run/docker.sock:/var/run/docker.sock
  command: --interval 300 --cleanup
  restart: unless-stopped
```

Watchtower akan cek update setiap 5 menit dan auto-restart container jika ada image baru.

---

## Environment Variables Reference

| Variable | Required | Description | Example |
|----------|----------|-------------|---------|
| `PORT` | Yes | Port aplikasi berjalan | `8080` |
| `OPENAI_API_KEY` | Yes | OpenAI API Key | `sk-proj-xxxxx...` |
| `ORS_API_KEY` | Yes | OpenRouteService API Key | `5b3ce3597851110001cf6248xxxxx...` |

---

## Quick Commands (SSH ke Server)

```bash
# Pull latest image
docker pull taufansena/chatbot-backend:latest

# Recreate container dengan image baru
docker stop police-assistant-backend
docker rm police-assistant-backend
docker run -d \
  --name police-assistant-backend \
  -p 8080:8080 \
  -e PORT=8080 \
  -e OPENAI_API_KEY=your_key \
  -e ORS_API_KEY=your_key \
  --restart unless-stopped \
  taufansena/chatbot-backend:latest

# Lihat logs realtime
docker logs -f police-assistant-backend
```
