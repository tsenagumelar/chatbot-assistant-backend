package main

import (
	"log"
	"police-assistant-backend/config"
	"police-assistant-backend/handlers"
	"police-assistant-backend/services"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func main() {
	// ASCII Art Banner
	log.Println(`
	╔═══════════════════════════════════════════════════╗
	║   🚓 AI POLICE ASSISTANT API                     ║
	║   Backend: Golang + Fiber                        ║
	║   AI: OpenAI ChatGPT                             ║
	║   Maps: OpenRouteService (FREE!)                 ║
	╚═══════════════════════════════════════════════════╝
	`)

	// Load configuration
	config.LoadConfig()

	// Initialize services
	log.Println("🔧 Initializing services...")
	openaiService := services.NewOpenAIService()
	orsService := services.NewORSService()
	etilangService := services.NewETilangService()
	pelayananService := services.NewPelayananService()
	simFlowService := services.NewSIMFlowService()

	// Initialize handlers
	chatHandler := handlers.NewChatHandler(openaiService, orsService, etilangService, pelayananService, simFlowService)
	trafficHandler := handlers.NewTrafficHandler(orsService)
	routeHandler := handlers.NewRouteHandler(orsService)
	sessionHandler := handlers.NewSessionHandler()

	// Create Fiber app with config
	app := fiber.New(fiber.Config{
		AppName:      "Police Assistant API v1.0",
		ServerHeader: "Fiber",
		ErrorHandler: customErrorHandler,
	})

	// Middleware
	app.Use(recover.New())
	app.Use(logger.New(logger.Config{
		Format:     "${time} | ${status} | ${latency} | ${method} ${path}\n",
		TimeFormat: "15:04:05",
	}))
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
	}))

	// Root endpoint
	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "ok",
			"message": "🚓 AI Police Assistant API is running",
			"version": "1.0.0",
			"endpoints": fiber.Map{
				"health":  "/health",
				"chat":    "/api/v1/chat",
				"session": "/api/v1/session",
				"traffic": "/api/v1/traffic",
				"routes":  "/api/v1/routes",
			},
		})
	})

	// Health check endpoint
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "healthy",
			"service": "police-assistant-api",
			"uptime":  "running",
		})
	})

	// API v1 routes
	api := app.Group("/api/v1")

	// Chat endpoints
	api.Post("/chat", chatHandler.HandleChat)

	// Session management endpoints
	api.Post("/session", sessionHandler.CreateSession)                  // Buat session baru
	api.Delete("/session/:session_id", sessionHandler.DeleteSession)    // Hapus session
	api.Post("/session/:session_id/clear", sessionHandler.ClearSession) // Clear history
	api.Get("/session/:session_id", sessionHandler.GetSessionInfo)      // Info session (debug)

	// Traffic endpoints
	api.Get("/traffic", trafficHandler.GetTraffic)

	// Route endpoints
	api.Post("/routes", routeHandler.GetRoutes)

	// 404 handler
	app.Use(func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error":   "Not Found",
			"message": "The requested endpoint does not exist",
		})
	})

	// Start server
	port := ":" + config.AppConfig.Port
	log.Printf("🚀 Server starting on port %s", config.AppConfig.Port)
	log.Printf("📝 API Documentation: http://localhost%s/", port)
	log.Printf("💚 Health Check: http://localhost%s/health", port)
	log.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	if err := app.Listen(port); err != nil {
		log.Fatalf("❌ Failed to start server: %v", err)
	}
}

// Custom error handler
func customErrorHandler(c *fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError

	if e, ok := err.(*fiber.Error); ok {
		code = e.Code
	}

	return c.Status(code).JSON(fiber.Map{
		"success": false,
		"error":   err.Error(),
		"code":    code,
	})
}
