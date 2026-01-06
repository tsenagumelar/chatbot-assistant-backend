# SIM Flow - Predefined Conversation Flow

## Overview

Sistem ini mengimplementasikan alur percakapan terstruktur (predefined flow) untuk proses perpanjangan dan pembuatan SIM. Alur ini menggunakan JSON flow definition yang memandu pengguna melalui langkah-langkah yang jelas dan terstruktur.

## Arsitektur

### File-file Terkait:
1. **perpanjangan_sim.json** - Definisi flow lengkap dengan nodes, transitions, dan choices
2. **services/sim_flow.go** - Service untuk parse dan navigate flow
3. **models/types.go** - Model data untuk SIM flow (SIMFlowInfo, SIMFlowChoice)
4. **handlers/chat.go** - Handler untuk deteksi dan proses SIM flow
5. **services/openai.go** - Inject flow context ke AI prompt
6. **services/session.go** - Session management untuk menyimpan state flow

### Flow Structure (perpanjangan_sim.json):
```json
{
  "flow_id": "perpanjangan_sim",
  "name": "Alur Perpanjangan SIM",
  "version": "1.0",
  "entry_node": "welcome",
  "nodes": [
    {
      "id": "welcome",
      "type": "message",
      "text": "Selamat datang di layanan perpanjangan SIM...",
      "transitions": [...]
    },
    ...
  ]
}
```

## Cara Kerja

### 1. Flow Detection
Sistem mendeteksi intent SIM melalui keyword matching:
```go
func (s *SIMFlowService) DetectSIMIntent(message string) bool {
    keywords := []string{"sim", "perpanjang", "buat", "bikin", "pembuatan"}
    // Check if message contains any keyword
}
```

### 2. State Management
State flow disimpan di session dengan key `sim_flow_current_node`:
```go
sessionStore.SetData(req.SessionID, "sim_flow_current_node", nodeID)
currentNodeID := sessionStore.GetData(req.SessionID, "sim_flow_current_node")
```

### 3. Node Navigation
Setiap node memiliki transitions dengan choices:
```go
type FlowTransition struct {
    Condition string       `json:"condition"`
    Target    string       `json:"target"`
    Choices   []FlowChoice `json:"choices,omitempty"`
}
```

User input dicocokkan dengan choices untuk menentukan next node.

### 4. AI Prompt Injection
Flow context di-inject ke system prompt agar AI mengikuti flow:
```go
if context.SIMFlowInfo != nil && context.SIMFlowInfo.Active {
    simFlowContext := "MODE ALUR PERPANJANGAN/PEMBUATAN SIM AKTIF..."
    // Include node text and choices
}
```

## Node Types

### 1. Message Node
Hanya menampilkan informasi, tidak butuh input:
```json
{
  "id": "welcome",
  "type": "message",
  "text": "Selamat datang..."
}
```

### 2. Question Node
Meminta user memilih dari opsi yang diberikan:
```json
{
  "id": "ask_sim_type",
  "type": "question",
  "text": "Jenis SIM apa yang ingin Anda perpanjang?",
  "transitions": [
    {
      "condition": "choice_based",
      "choices": [
        {"id": "sim_a", "label": "SIM A", "target": "ask_ever_had_sim"},
        {"id": "sim_c", "label": "SIM C", "target": "ask_ever_had_sim"}
      ]
    }
  ]
}
```

### 3. Collect Node
Meminta user untuk upload dokumen:
```json
{
  "id": "collect_ktp",
  "type": "collect",
  "text": "Silakan upload foto KTP Anda",
  "collect": {
    "field": "ktp",
    "validation": "required"
  }
}
```

### 4. Action Node
Proses tertentu (misal: compile data, submit):
```json
{
  "id": "submit_perpanjangan",
  "type": "action",
  "text": "Memproses permohonan perpanjangan SIM...",
  "action": {
    "type": "submit_form",
    "endpoint": "/api/sim/submit"
  }
}
```

## Flow Example: Perpanjangan SIM A

```
1. welcome (message)
   ↓
2. ask_sim_type (question)
   User pilih: "SIM A"
   ↓
3. ask_ever_had_sim (question)
   User jawab: "Ya, pernah"
   ↓
4. ask_sim_status (question)
   User jawab: "Masih berlaku"
   ↓
5. offer_help_perpanjang (question)
   User jawab: "Ya, tolong bantu"
   ↓
6. collect_ktp (collect)
   User upload KTP
   ↓
7. collect_sim_lama (collect)
   User upload SIM lama
   ↓
8. collect_surat_sehat (collect)
   User upload surat sehat
   ↓
9. collect_tes_psikologi (collect)
   User upload hasil tes psikologi
   ↓
10. compile_perpanjangan (action)
    System compile data
    ↓
11. ready_to_submit (message)
    Konfirmasi siap submit
```

## API Response dengan SIM Flow

Ketika SIM flow aktif, response akan include `sim_flow_info`:

```json
{
  "success": true,
  "response": "Jenis SIM apa yang ingin Anda perpanjang?\n\n1. SIM A\n2. SIM C\n3. SIM D",
  "session_id": "abc123",
  "sim_flow_info": {
    "active": true,
    "current_node": "ask_sim_type",
    "node_type": "question",
    "node_text": "Jenis SIM apa yang ingin Anda perpanjang?",
    "choices": [
      {"id": "sim_a", "label": "SIM A"},
      {"id": "sim_c", "label": "SIM C"},
      {"id": "sim_d", "label": "SIM D"}
    ]
  }
}
```

## Testing

Gunakan test script yang disediakan:

```bash
./test-sim-flow.sh
```

Script ini akan menjalankan full flow dari awal sampai akhir, termasuk simulasi upload dokumen.

## Manual Test dengan curl:

### 1. Trigger SIM Flow
```bash
curl -X POST http://localhost:3000/api/chat \
  -H "Content-Type: application/json" \
  -d '{
    "message": "Saya mau perpanjang SIM",
    "context": {"location": "Jakarta"}
  }'
```

### 2. Pilih SIM Type
```bash
curl -X POST http://localhost:3000/api/chat \
  -H "Content-Type: application/json" \
  -d '{
    "message": "SIM A",
    "session_id": "<session_id_dari_response_sebelumnya>",
    "context": {"location": "Jakarta"}
  }'
```

### 3. Upload Document
```bash
curl -X POST http://localhost:3000/api/chat \
  -H "Content-Type: application/json" \
  -d '{
    "message": "Ini KTP saya",
    "session_id": "<session_id>",
    "context": {"location": "Jakarta"},
    "documents": [
      {
        "file_name": "ktp.jpg",
        "file_type": "image/jpeg",
        "file_url": "https://example.com/ktp.jpg",
        "description": "KTP"
      }
    ]
  }'
```

## Catatan Penting

1. **Session Required**: Flow memerlukan session yang konsisten untuk tracking state
2. **Exact Match**: Choice matching menggunakan fuzzy matching untuk fleksibilitas
3. **AI Guidance**: AI dipandu oleh flow context tapi tetap bisa handle natural conversation
4. **Document Collection**: Setiap node "collect" meminta upload dokumen tertentu
5. **State Persistence**: State flow disimpan di session dan akan expired setelah 24 jam inactive

## Menambah Flow Baru

Untuk menambah flow baru:

1. Buat JSON file baru (misal: `pembuatan_skck.json`)
2. Definisikan nodes, transitions, dan choices
3. Buat service baru atau extend `SIMFlowService`
4. Update chat handler untuk detect flow baru
5. Inject flow context ke AI prompt

## Troubleshooting

### Flow tidak terdeteksi
- Check keyword matching di `DetectSIMIntent()`
- Pastikan message mengandung keyword yang sesuai

### Node tidak berpindah
- Check transition conditions
- Verify choice matching logic
- Check session state (`sim_flow_current_node`)

### AI tidak mengikuti flow
- Verify flow context injection di `buildSystemPrompt()`
- Check `SIMFlowInfo` di context
- Review AI prompt instructions untuk flow

## Future Enhancements

1. **Multiple Flows**: Support untuk multiple flows secara parallel
2. **Flow Analytics**: Track completion rate, drop-off points
3. **Dynamic Flows**: Generate flow dari database/config
4. **Flow Builder**: UI untuk membuat dan edit flow
5. **Validation**: Validate uploaded documents (OCR, format check)
6. **Reminder**: Kirim reminder jika user drop mid-flow
