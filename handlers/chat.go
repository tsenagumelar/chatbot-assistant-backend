package handlers

import (
	"log"
	"police-assistant-backend/models"
	"police-assistant-backend/services"
	"strings"

	"github.com/gofiber/fiber/v2"
)

type ChatHandler struct {
	openaiService  *services.OpenAIService
	orsService     *services.ORSService
	etilangService *services.ETilangService
}

func NewChatHandler(openaiService *services.OpenAIService, orsService *services.ORSService, etilangService *services.ETilangService) *ChatHandler {
	return &ChatHandler{
		openaiService:  openaiService,
		orsService:     orsService,
		etilangService: etilangService,
	}
}

func (h *ChatHandler) HandleChat(c *fiber.Ctx) error {
	var req models.ChatRequest

	// Parse request body
	if err := c.BodyParser(&req); err != nil {
		log.Printf("❌ Failed to parse request: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(models.ChatResponse{
			Success: false,
			Error:   "Invalid request body",
		})
	}

	// Validate message
	if req.Message == "" {
		return c.Status(fiber.StatusBadRequest).JSON(models.ChatResponse{
			Success: false,
			Error:   "Message is required",
		})
	}

	// Get session store
	sessionStore := services.GetSessionStore()

	// Jika tidak ada session_id, buat session baru
	if req.SessionID == "" {
		req.SessionID = sessionStore.CreateSession()
		log.Printf("🆕 Created new session: %s", req.SessionID)
	} else {
		log.Printf("📝 Using existing session: %s", req.SessionID)
	}

	log.Printf("💬 Chat request: %s", req.Message)
	log.Printf("📍 Context: Location=%s, Speed=%.1f km/h, Traffic=%s",
		req.Context.Location, req.Context.Speed, req.Context.Traffic)

	// Check if user is asking about e-tilang
	messageLower := strings.ToLower(req.Message)
	if strings.Contains(messageLower, "tilang") || strings.Contains(messageLower, "pelanggaran") ||
		strings.Contains(messageLower, "denda") || strings.Contains(messageLower, "cek") && (strings.Contains(messageLower, "polisi") || strings.Contains(messageLower, "nopol")) {

		// Try to extract plate number
		plateNumber := h.etilangService.ExtractPlateNumber(req.Message)
		if plateNumber != "" {
			log.Printf("🚗 E-Tilang check requested for plate: %s", plateNumber)

			// Get e-tilang info
			etilangInfo := h.etilangService.CheckETilang(plateNumber)
			req.Context.ETilangInfo = etilangInfo

			log.Printf("📋 E-Tilang info attached: HasViolation=%v, TotalFine=%d",
				etilangInfo.HasViolation, etilangInfo.TotalFine)
		}
	}

	// If location is empty but coordinates are provided, do reverse geocoding
	if req.Context.Location == "" && req.Context.Latitude != 0 && req.Context.Longitude != 0 {
		address, err := h.orsService.ReverseGeocode(req.Context.Latitude, req.Context.Longitude)
		if err == nil {
			req.Context.Location = address
			log.Printf("🗺️  Reverse geocoded location: %s", address)
		}
	}

	// Ambil history dari session (prioritas: backend storage > request body)
	var history []models.OpenAIMessage
	if req.SessionID != "" {
		history = sessionStore.GetHistory(req.SessionID)
		if len(history) > 0 {
			log.Printf("📚 Using %d messages from session history", len(history))
		}
	} else if len(req.History) > 0 {
		// Backward compatibility: jika tidak ada session_id tapi ada history di request
		history = req.History
		log.Printf("📚 Using %d messages from request history", len(history))
	}

	// Call OpenAI API with history
	response, err := h.openaiService.Chat(req.Message, req.Context, history)
	if err != nil {
		log.Printf("❌ OpenAI error: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(models.ChatResponse{
			Success: false,
			Error:   "Failed to get AI response: " + err.Error(),
		})
	}

	// Simpan pesan user dan response ke session history
	sessionStore.AddMessage(req.SessionID, "user", req.Message)
	sessionStore.AddMessage(req.SessionID, "assistant", response)

	log.Printf("✅ Chat response generated successfully (session: %s)", req.SessionID)

	return c.JSON(models.ChatResponse{
		Success:   true,
		Response:  response,
		SessionID: req.SessionID, // Return session ID ke frontend
	})
}
