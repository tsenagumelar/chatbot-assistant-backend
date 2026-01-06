package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port         string
	OpenAIAPIKey string
	ORSAPIKey    string // OpenRouteService API Key
	OpenAIModel  string
}

var AppConfig *Config

func LoadConfig() {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: .env file not found, using environment variables")
	}

	AppConfig = &Config{
		Port:         getEnv("PORT", "8080"),
		OpenAIAPIKey: getEnv("OPENAI_API_KEY", ""),
		ORSAPIKey:    getEnv("OPENROUTESERVICE_API_KEY", ""),
		OpenAIModel:  getEnv("OPENAI_MODEL", "gpt-5.1"),
	}

	// Validate required keys
	if AppConfig.OpenAIAPIKey == "" {
		log.Fatal("❌ OPENAI_API_KEY is required in .env file")
	}
	if AppConfig.ORSAPIKey == "" {
		log.Fatal("❌ OPENROUTESERVICE_API_KEY is required in .env file")
	}

	log.Println("✅ Configuration loaded successfully")
	log.Printf("📝 Using OpenAI Model: %s", AppConfig.OpenAIModel)
	log.Println("🗺️  Using OpenRouteService (Free Maps API)")
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
