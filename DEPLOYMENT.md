# Deployment Checklist

## ⚠️ PENTING - File yang Harus Ada di Production Server

### Required JSON Files
Pastikan file-file berikut ada di directory yang sama dengan binary executable:

1. ✅ `response-rules.json` - Conversation flow rules
2. ✅ `location-rules.json` - Location routing rules  
3. ✅ `data_pelayanan.json` - Service information data
4. ✅ `perpanjangan_sim.json` - SIM renewal flow (optional, if using SIM flow)

### File Structure di Server

```
/app/
├── chatbot-assistant          # Binary executable
├── response-rules.json        # ⚠️ HARUS ADA
├── location-rules.json        # ⚠️ HARUS ADA
├── data_pelayanan.json        # ⚠️ HARUS ADA
└── perpanjangan_sim.json      # Optional
```

## Common Error: "nil pointer dereference"

### Penyebab
```
runtime error: invalid memory address or nil pointer dereference
```

Error ini terjadi jika:
1. File JSON tidak ada di server
2. Path file tidak benar
3. Permission denied untuk read file

### Solusi

#### Option 1: Copy Manual Files ke Server
```bash
# Copy all JSON files to server
scp response-rules.json user@server:/app/
scp location-rules.json user@server:/app/
scp data_pelayanan.json user@server:/app/
scp perpanjangan_sim.json user@server:/app/
```

#### Option 2: Update Dockerfile
```dockerfile
# Tambahkan di Dockerfile
COPY response-rules.json .
COPY location-rules.json .
COPY data_pelayanan.json .
COPY perpanjangan_sim.json .
```

#### Option 3: Docker Compose Volume
```yaml
# docker-compose.yml
services:
  chatbot-api:
    volumes:
      - ./response-rules.json:/app/response-rules.json:ro
      - ./location-rules.json:/app/location-rules.json:ro
      - ./data_pelayanan.json:/app/data_pelayanan.json:ro
      - ./perpanjangan_sim.json:/app/perpanjangan_sim.json:ro
```

## Verify Deployment

### 1. Check if Files Exist
```bash
# SSH ke server
ssh user@server

# Check files
ls -la /app/*.json

# Should show:
# -rw-r--r-- 1 root root  xxxx response-rules.json
# -rw-r--r-- 1 root root  xxxx location-rules.json
# -rw-r--r-- 1 root root  xxxx data_pelayanan.json
```

### 2. Check Server Logs
```bash
# Look for success messages
docker logs chatbot-api | grep "Rules"

# Should see:
# ✅ Response Rules loaded: XX rules
# ✅ Location Rules loaded: XX rules
# ✅ Rules Service initialized

# If you see WARNING:
# ⚠️ WARNING: Failed to load response-rules.json: open response-rules.json: no such file or directory
# → File tidak ada, perlu di-copy
```

### 3. Test API
```bash
curl -X POST https://chatbotapi.activa.id/api/v1/chat \
  -H "Content-Type: application/json" \
  -d '{
    "message": "mau bikin SIM A",
    "name": "Test User",
    "session_id": "test-123"
  }'

# Should return proper response without 500 error
```

## Docker Deployment Steps

### Build with JSON Files
```bash
# 1. Build Docker image
docker build -t chatbot-assistant:latest .

# 2. Verify JSON files are in image
docker run --rm chatbot-assistant:latest ls -la *.json

# 3. Run container
docker run -d \
  --name chatbot-api \
  -p 8080:8080 \
  -e OPENAI_API_KEY="your-key" \
  chatbot-assistant:latest
```

### Using Docker Compose
```bash
# Make sure docker-compose.yml has volume mounts
docker-compose up -d

# Check logs
docker-compose logs -f chatbot-api
```

## Portainer Stack Deployment

If using Portainer stack (`portainer-stack-dev.yml`):

1. **Upload JSON files ke server terlebih dahulu**
   ```bash
   scp *.json user@server:/path/to/app/
   ```

2. **Update stack configuration** to mount volumes
   ```yaml
   version: '3.8'
   services:
     chatbot-api:
       image: your-registry/chatbot-assistant:latest
       volumes:
         - /path/to/app/response-rules.json:/app/response-rules.json:ro
         - /path/to/app/location-rules.json:/app/location-rules.json:ro
         - /path/to/app/data_pelayanan.json:/app/data_pelayanan.json:ro
   ```

3. **Deploy/Update stack** via Portainer UI

## Environment Variables

Tidak ada environment variable khusus untuk JSON files. Service akan:
- Load dari current working directory
- Log WARNING jika file tidak ada
- Continue berjalan dengan fallback behavior (tanpa rules)

## Troubleshooting

### Issue: 500 error "nil pointer dereference"
**Solution:** 
- Check file JSON ada di server
- Check permission file (chmod 644)
- Check logs untuk error message detail

### Issue: Rules tidak terload tapi tidak ada error
**Solution:**
- Check JSON format valid (use `jq . response-rules.json`)
- Check field names match struct tags
- Check file encoding (UTF-8)

### Issue: File ada tapi tetap error "no such file"
**Solution:**
- Check working directory: `pwd` di dalam container
- Update path ke absolute: `/app/response-rules.json`
- Check container volumes dengan `docker inspect`

## Testing Production

```bash
# Test with real session
curl -X POST https://chatbotapi.activa.id/api/v1/chat \
  -H "Content-Type: application/json" \
  -d '{
    "message": "mau perpanjang SIM",
    "name": "Taufan Gumelar",
    "location": "Cibeureum, West Java",
    "session_id": "b70d2c64-74fb-43a3-b498-a06097f0bd04"
  }'

# Should return response with:
# - Personalized greeting with "Taufan Gumelar"
# - Flow from response-rules.json
# - No 500 error
```

## Quick Fix for Production

Jika server production sedang error dan perlu fix cepat:

```bash
# 1. SSH ke server
ssh user@production-server

# 2. Navigate ke app directory
cd /app  # atau directory dimana binary ada

# 3. Copy JSON files dari local
scp -r user@local:/path/to/json/*.json .

# 4. Verify
ls -la *.json

# 5. Restart service
docker restart chatbot-api
# atau
systemctl restart chatbot-assistant

# 6. Check logs
docker logs -f chatbot-api
```

## Monitoring

Add logging to track file loading:
```go
// Already implemented in rules.go:
log.Printf("✅ Response Rules loaded: %d rules", len(service.responseRules.Sheet1))
log.Printf("✅ Location Rules loaded: %d rules", len(service.locationRules))
```

Monitor these messages on server startup to ensure files are loaded successfully.
