package handlers

import (
	"log"
	"police-assistant-backend/models"
	"police-assistant-backend/services"
	"strings"

	"github.com/gofiber/fiber/v2"
)

type ChatHandler struct {
	openaiService    *services.OpenAIService
	orsService       *services.ORSService
	etilangService   *services.ETilangService
	pelayananService *services.PelayananService
	simFlowService   *services.SIMFlowService
}

func NewChatHandler(openaiService *services.OpenAIService, orsService *services.ORSService, etilangService *services.ETilangService, pelayananService *services.PelayananService, simFlowService *services.SIMFlowService) *ChatHandler {
	return &ChatHandler{
		openaiService:    openaiService,
		orsService:       orsService,
		etilangService:   etilangService,
		pelayananService: pelayananService,
		simFlowService:   simFlowService,
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

	// Check if user uploaded documents
	if len(req.Documents) > 0 {
		req.Context.HasUploadedDocuments = true
		req.Context.UploadedDocumentCount = len(req.Documents)
		log.Printf("📤 User uploaded %d document(s)", len(req.Documents))

		// Log document details
		for i, doc := range req.Documents {
			log.Printf("   %d. %s (%s) - %s", i+1, doc.FileName, doc.FileType, doc.Description)
		}
	}

	// Check if user is asking about pelayanan (services)
	messageLower := strings.ToLower(req.Message)
	pelayananKeywords := []string{
		"pelayanan", "layanan", "sim", "stnk", "pajak", "balik nama", "mutasi",
		"perpanjang", "buat", "bikin", "ganti", "hilang", "kehilangan",
		"pengesahan", "dokumen", "syarat", "persyaratan",
	}

	shouldCheckPelayanan := false
	for _, keyword := range pelayananKeywords {
		if strings.Contains(messageLower, keyword) {
			shouldCheckPelayanan = true
			break
		}
	}

	if shouldCheckPelayanan {
		log.Printf("📋 Pelayanan check requested")
		pelayananInfo := h.pelayananService.SearchPelayanan(req.Message)
		if pelayananInfo.Found {
			req.Context.PelayananInfo = pelayananInfo
			log.Printf("✅ Pelayanan info attached: %s", pelayananInfo.Pelayanan.JenisPelayanan)
		}
	}

	// Check if user is talking about SIM renewal (perpanjangan/pembuatan SIM)
	if h.simFlowService.DetectSIMIntent(req.Message) {
		log.Printf("🪪 SIM flow detected")

		// Get current node from session (if exists)
		currentNodeID := sessionStore.GetData(req.SessionID, "sim_flow_current_node")
		if currentNodeID == "" {
			// First time, start from entry node
			currentNodeID = "entry_node"
			sessionStore.SetData(req.SessionID, "sim_flow_current_node", currentNodeID)
			log.Printf("🆕 Starting SIM flow from entry_node")
		} else {
			log.Printf("📍 Continuing SIM flow from node: %s", currentNodeID)
		}

		// Get current node
		currentNode := h.simFlowService.GetCurrentNode(currentNodeID)
		if currentNode != nil {
			// Process user choice to get next node
			nextNodeID, nextNode := h.simFlowService.ProcessUserChoice(currentNode, req.Message)

			if nextNode != nil {
				// Update session to next node
				sessionStore.SetData(req.SessionID, "sim_flow_current_node", nextNodeID)
				log.Printf("➡️  Moving to next node: %s (type: %s)", nextNodeID, nextNode.Type)

				// Get flow info for this node
				flowInfo := h.simFlowService.GetSIMFlowInfo(nextNodeID)
				req.Context.SIMFlowInfo = flowInfo
			} else {
				// No transition matched, stay on current node
				log.Printf("⏸️  No transition matched, staying on node: %s", currentNodeID)
				flowInfo := h.simFlowService.GetSIMFlowInfo(currentNodeID)
				req.Context.SIMFlowInfo = flowInfo
			}
		}
	}

	// Check if user is asking about e-tilang
	messageLower = strings.ToLower(req.Message)
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
		Success:       true,
		Response:      response,
		SessionID:     req.SessionID, // Return session ID ke frontend
		ETilangInfo:   req.Context.ETilangInfo,
		PelayananInfo: req.Context.PelayananInfo,
		SIMFlowInfo:   req.Context.SIMFlowInfo,
	})
}
