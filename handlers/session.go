package handlers

import (
	"log"
	"police-assistant-backend/models"
	"police-assistant-backend/services"

	"github.com/gofiber/fiber/v2"
)

type SessionHandler struct {
	sessionStore *services.SessionStore
}

func NewSessionHandler() *SessionHandler {
	return &SessionHandler{
		sessionStore: services.GetSessionStore(),
	}
}

// CreateSession membuat session baru
func (h *SessionHandler) CreateSession(c *fiber.Ctx) error {
	sessionID := h.sessionStore.CreateSession()
	log.Printf("🆕 New session created: %s", sessionID)

	return c.JSON(models.SessionResponse{
		Success:   true,
		SessionID: sessionID,
		Message:   "Session created successfully",
	})
}

// ClearSession menghapus history dari session
func (h *SessionHandler) ClearSession(c *fiber.Ctx) error {
	sessionID := c.Params("session_id")

	if sessionID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(models.SessionResponse{
			Success: false,
			Error:   "Session ID is required",
		})
	}

	h.sessionStore.ClearSession(sessionID)
	log.Printf("🧹 Session cleared: %s", sessionID)

	return c.JSON(models.SessionResponse{
		Success:   true,
		SessionID: sessionID,
		Message:   "Session history cleared",
	})
}

// DeleteSession menghapus session
func (h *SessionHandler) DeleteSession(c *fiber.Ctx) error {
	sessionID := c.Params("session_id")

	if sessionID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(models.SessionResponse{
			Success: false,
			Error:   "Session ID is required",
		})
	}

	h.sessionStore.DeleteSession(sessionID)
	log.Printf("🗑️  Session deleted: %s", sessionID)

	return c.JSON(models.SessionResponse{
		Success: true,
		Message: "Session deleted successfully",
	})
}

// GetSessionInfo menampilkan info session (untuk debugging)
func (h *SessionHandler) GetSessionInfo(c *fiber.Ctx) error {
	sessionID := c.Params("session_id")

	if sessionID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Session ID is required",
		})
	}

	session, exists := h.sessionStore.GetSession(sessionID)
	if !exists {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"error":   "Session not found",
		})
	}

	return c.JSON(fiber.Map{
		"success":       true,
		"session_id":    session.ID,
		"message_count": len(session.History),
		"created_at":    session.CreatedAt,
		"updated_at":    session.UpdatedAt,
	})
}
