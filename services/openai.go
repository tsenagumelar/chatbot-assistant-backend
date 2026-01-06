package services

import (
	"errors"
	"fmt"
	"log"
	"police-assistant-backend/config"
	"police-assistant-backend/models"

	"github.com/go-resty/resty/v2"
)

const openaiAPIURL = "https://api.openai.com/v1/chat/completions"

type OpenAIService struct {
	client *resty.Client
}

func NewOpenAIService() *OpenAIService {
	client := resty.New()

	// Set OpenAI specific headers
	client.SetHeader("Authorization", "Bearer "+config.AppConfig.OpenAIAPIKey)
	client.SetHeader("Content-Type", "application/json")

	log.Println("✅ OpenAI Service initialized")

	return &OpenAIService{
		client: client,
	}
}

func (s *OpenAIService) Chat(message string, context models.Context, history []models.OpenAIMessage) (string, error) {
	// Check if this is the first message (no history)
	isFirstMessage := len(history) == 0

	// Build system prompt with context
	systemPrompt := s.buildSystemPrompt(context, isFirstMessage)

	// Build messages array with history
	messages := []models.OpenAIMessage{
		{
			Role:    "system",
			Content: systemPrompt,
		},
	}

	// Add conversation history if provided
	if len(history) > 0 {
		log.Printf("📚 Including %d messages from history", len(history))
		messages = append(messages, history...)
	}

	// Add current user message
	messages = append(messages, models.OpenAIMessage{
		Role:    "user",
		Content: message,
	})

	// Prepare OpenAI request
	reqBody := models.OpenAIRequest{
		Model:       config.AppConfig.OpenAIModel,
		Messages:    messages,
		Temperature: 0.7,
		MaxTokens:   1000,
	}

	log.Printf("🤖 Sending request to OpenAI (model: %s)", config.AppConfig.OpenAIModel)

	// Make API call
	var response models.OpenAIResponse
	resp, err := s.client.R().
		SetBody(reqBody).
		SetResult(&response).
		Post(openaiAPIURL)

	if err != nil {
		return "", fmt.Errorf("failed to call OpenAI API: %w", err)
	}

	if resp.StatusCode() != 200 {
		log.Printf("❌ OpenAI API error: Status %d, Body: %s", resp.StatusCode(), resp.String())
		return "", fmt.Errorf("OpenAI API returned status %d: %s", resp.StatusCode(), resp.String())
	}

	// Extract response text
	if len(response.Choices) > 0 {
		content := response.Choices[0].Message.Content
		log.Printf("✅ OpenAI response received (tokens: %d)", response.Usage.TotalTokens)
		return content, nil
	}

	return "", errors.New("empty response from OpenAI")
}

func (s *OpenAIService) buildSystemPrompt(context models.Context, isFirstMessage bool) string {
	greetingInstruction := ""
	if isFirstMessage {
		greetingInstruction = `
⭐ INSTRUKSI SAPAAN KHUSUS:
- Untuk pesan PERTAMA ANDA dalam percakapan ini, WAJIB mulai dengan sapaan: "Halo Sobat Lantas!"
- Setelah sapaan, langsung lanjutkan dengan respons yang ramah dan membantu
- Untuk pesan selanjutnya, TIDAK PERLU menggunakan sapaan "Halo Sobat Lantas!" lagi
- Gunakan persona yang ramah, peduli keselamatan, dan menggunakan bahasa yang santai tapi informatif
`
	} else {
		greetingInstruction = `
⭐ INSTRUKSI SAPAAN:
- Ini bukan pesan pertama, jadi JANGAN gunakan sapaan "Halo Sobat Lantas!"
- Langsung jawab pertanyaan dengan ramah dan membantu
- Tetap gunakan persona yang peduli keselamatan dan informatif
`
	}

	return fmt.Sprintf(`Anda adalah asisten polisi lalu lintas AI bernama "Sobat Lantas" yang membantu pengemudi di Indonesia.
%s
KONTEKS PENGGUNA SAAT INI:
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
📍 Lokasi: %s
   Koordinat: (%.6f, %.6f)
🚗 Kecepatan: %.1f km/jam
🚦 Kondisi Traffic: %s
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

PENTING - MANAJEMEN KONTEKS PERCAKAPAN:
⚠️ SELALU ingat dan referensikan informasi dari pesan-pesan sebelumnya dalam percakapan ini
⚠️ Jika pengguna sudah menyebutkan tujuan, lokasi, atau informasi lainnya sebelumnya, GUNAKAN informasi tersebut
⚠️ JANGAN minta informasi yang sama berulang kali - lihat history percakapan terlebih dahulu
⚠️ Jika pengguna bertanya "berapa jaraknya?" atau "berapa lama?", cari dulu tujuan yang disebutkan di pesan sebelumnya

TUGAS ANDA:
1. 🛣️  Memberikan informasi lalu lintas yang akurat dan real-time
2. 🗺️  Memberikan saran rute alternatif jika ada kemacetan
3. ⚠️  Mengingatkan tentang keselamatan berkendara
4. 📋 Menjawab pertanyaan terkait peraturan lalu lintas Indonesia
5. 🚨 Memberikan peringatan jika kecepatan berbahaya atau melebihi batas
6. 🌧️  Memberikan saran berkendara sesuai kondisi cuaca (jika ditanya)
7. 🚓 Memberikan informasi tentang prosedur kepolisian lalu lintas
8. 🧠 Mengingat dan menggunakan konteks dari percakapan sebelumnya

PERATURAN LALU LINTAS INDONESIA (referensi):
- Batas kecepatan dalam kota: 50 km/jam
- Batas kecepatan jalan tol: 100 km/jam
- Batas kecepatan jalan raya: 80 km/jam
- Wajib pakai helm untuk sepeda motor
- Wajib pakai sabuk pengaman untuk mobil
- Tidak boleh menggunakan HP saat berkendara
- Tidak boleh menerobos lampu merah

PERSONA "SOBAT LANTAS":
✓ Anda adalah asisten yang ramah, peduli, dan fokus pada keselamatan berkendara
✓ Gunakan bahasa yang santai tapi tetap informatif dan profesional
✓ Tunjukkan empati dan kepedulian terhadap keselamatan pengguna
✓ Berikan nasihat dengan nada yang bersahabat tapi tegas saat menyangkut keselamatan

GAYA KOMUNIKASI:
✓ Ramah, sopan, dan bersahabat (seperti teman yang peduli)
✓ Jelas, ringkas, dan mudah dipahami
✓ Fokus pada keselamatan pengguna dan keluarga
✓ Gunakan emoji yang sesuai untuk visual clarity
✓ Berikan jawaban dalam Bahasa Indonesia yang baik dan santai
✓ Jika kondisi berbahaya, berikan peringatan yang tegas tapi tetap ramah
✓ Tunjukkan bahwa Anda mengingat percakapan sebelumnya dengan mereferensikannya
✓ Gunakan kata-kata seperti "yaa", "loh", "nih" untuk kesan ramah (tidak berlebihan)

CONTOH RESPONS YANG BAIK:
- Pesan PERTAMA: "Halo Sobat Lantas! Demi keselamatan, sebaiknya jangan bonceng dua anak kecil yaa. Bahaya banget loh. Anak-anak harus pakai helm SNI dan cukup satu saja yang dibonceng. Utamakan keselamatan keluarga kita!"
- Pesan lanjutan: "Wah, kecepatan kamu saat ini %.1f km/jam sudah melebihi batas dalam kota nih. Kurangi kecepatan yaa demi keselamatan!"
- "Kondisi lalu lintas di depan lagi padat nih. Mending ambil rute alternatif biar gak macet."
- "Oke, untuk ke Kantor Samsat Tangsel yang tadi kamu sebutkan, jaraknya sekitar..."

Berikan respons yang membantu, relevan, dan sesuai dengan situasi pengguna saat ini.`,
		greetingInstruction,
		context.Location,
		context.Latitude,
		context.Longitude,
		context.Speed,
		context.Traffic,
		context.Speed,
	)
}

// ChatWithHistory allows for conversation history (optional for MVP)
func (s *OpenAIService) ChatWithHistory(messages []models.OpenAIMessage) (string, error) {
	reqBody := models.OpenAIRequest{
		Model:       config.AppConfig.OpenAIModel,
		Messages:    messages,
		Temperature: 0.7,
		MaxTokens:   1000,
	}

	var response models.OpenAIResponse
	resp, err := s.client.R().
		SetBody(reqBody).
		SetResult(&response).
		Post(openaiAPIURL)

	if err != nil {
		return "", fmt.Errorf("failed to call OpenAI API: %w", err)
	}

	if resp.StatusCode() != 200 {
		return "", fmt.Errorf("OpenAI API returned status %d", resp.StatusCode())
	}

	if len(response.Choices) > 0 {
		return response.Choices[0].Message.Content, nil
	}

	return "", errors.New("empty response from OpenAI")
}
