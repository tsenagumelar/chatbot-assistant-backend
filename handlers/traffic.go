package handlers

import (
	"log"
	"police-assistant-backend/models"
	"police-assistant-backend/services"

	"github.com/gofiber/fiber/v2"
)

type TrafficHandler struct {
	mapsService *services.ORSService
}

func NewTrafficHandler(mapsService *services.ORSService) *TrafficHandler {
	return &TrafficHandler{
		mapsService: mapsService,
	}
}

// GetTraffic handles GET /api/v1/traffic
// Query params: latitude, longitude
func (h *TrafficHandler) GetTraffic(c *fiber.Ctx) error {
	var req models.TrafficRequest

	// Parse query parameters
	if err := c.QueryParser(&req); err != nil {
		log.Printf("❌ Failed to parse query: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(models.TrafficResponse{
			Success: false,
			Error:   "Invalid query parameters",
		})
	}

	// Validate coordinates
	if req.Latitude == 0 || req.Longitude == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(models.TrafficResponse{
			Success: false,
			Error:   "Latitude and longitude are required",
		})
	}

	log.Printf("🚦 Getting traffic info for: %.6f, %.6f", req.Latitude, req.Longitude)

	// Get traffic information from Google Maps
	traffic, err := h.mapsService.GetTrafficInfo(req.Latitude, req.Longitude)
	if err != nil {
		log.Printf("❌ Failed to get traffic: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(models.TrafficResponse{
			Success: false,
			Error:   "Failed to get traffic information: " + err.Error(),
		})
	}

	log.Printf("✅ Traffic info retrieved successfully")

	return c.JSON(models.TrafficResponse{
		Success: true,
		Traffic: traffic,
	})
}
