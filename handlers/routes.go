package handlers

import (
	"log"
	"police-assistant-backend/models"
	"police-assistant-backend/services"

	"github.com/gofiber/fiber/v2"
)

type RouteHandler struct {
	orsService *services.ORSService
}

func NewRouteHandler(orsService *services.ORSService) *RouteHandler {
	return &RouteHandler{
		orsService: orsService,
	}
}

// GetRoutes handles POST /api/v1/routes
// Body: { "origin": "Location A", "destination": "Location B" }
func (h *RouteHandler) GetRoutes(c *fiber.Ctx) error {
	var req models.RouteRequest

	// Parse request body
	if err := c.BodyParser(&req); err != nil {
		log.Printf("❌ Failed to parse request: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(models.RouteResponse{
			Success: false,
			Error:   "Invalid request body",
		})
	}

	// Validate origin and destination
	if req.Origin == "" || req.Destination == "" {
		return c.Status(fiber.StatusBadRequest).JSON(models.RouteResponse{
			Success: false,
			Error:   "Origin and destination are required",
		})
	}

	log.Printf("🗺️  Finding routes from '%s' to '%s'", req.Origin, req.Destination)

	// Get alternative routes with traffic from OpenRouteService
	routes, err := h.orsService.GetAlternativeRoutes(req.Origin, req.Destination)
	if err != nil {
		log.Printf("❌ Failed to get routes: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(models.RouteResponse{
			Success: false,
			Error:   "Failed to get routes: " + err.Error(),
		})
	}

	log.Printf("✅ Found %d alternative route(s)", len(routes))

	return c.JSON(models.RouteResponse{
		Success: true,
		Routes:  routes,
	})
}
