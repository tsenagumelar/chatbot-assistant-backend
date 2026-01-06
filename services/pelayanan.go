package services

import (
	"encoding/json"
	"log"
	"os"
	"police-assistant-backend/models"
	"strings"
)

type PelayananService struct {
	dataset models.PelayananDataset
}

func NewPelayananService() *PelayananService {
	service := &PelayananService{}

	// Load data from JSON file
	if err := service.loadData(); err != nil {
		log.Printf("⚠️  Failed to load pelayanan data: %v", err)
	} else {
		log.Printf("✅ Pelayanan Service initialized with %d flows", len(service.dataset.Flows))
	}

	return service
}

func (s *PelayananService) loadData() error {
	// Read JSON file
	file, err := os.ReadFile("data_pelayanan.json")
	if err != nil {
		return err
	}

	// Parse JSON into new structure
	if err := json.Unmarshal(file, &s.dataset); err != nil {
		return err
	}

	return nil
}

// SearchPelayanan searches for service flow based on user query
func (s *PelayananService) SearchPelayanan(query string) *models.PelayananInfo {
	queryLower := strings.ToLower(query)

	// Keywords untuk berbagai jenis pelayanan
	keywords := map[string][]string{
		"SIM":        {"sim", "surat izin mengemudi", "bikin sim", "buat sim", "perpanjang sim", "perpanjangan sim", "sim hilang", "sim rusak", "sim internasional"},
		"STNK":       {"stnk", "pajak kendaraan", "pajak motor", "pajak mobil", "pengesahan stnk", "stnk hilang", "ganti data stnk", "lapor kehilangan stnk"},
		"Tilang":     {"tilang", "etle", "e-tilang", "cek tilang", "pelanggaran"},
		"BPN":        {"balik nama", "mutasi kendaraan", "mutasi"},
		"Kehilangan": {"hilang", "kehilangan", "lapor kehilangan"},
	}

	// Try to match with keywords first
	matchedType := ""
	for serviceType, words := range keywords {
		for _, word := range words {
			if strings.Contains(queryLower, word) {
				matchedType = serviceType
				break
			}
		}
		if matchedType != "" {
			break
		}
	}

	// Search in flows based on matched type and specific keywords
	for _, flow := range s.dataset.Flows {
		titleLower := strings.ToLower(flow.Title)

		// Direct match with flow title
		if strings.Contains(titleLower, queryLower) {
			log.Printf("🔍 Found flow: %s (direct match)", flow.Title)
			return &models.PelayananInfo{
				Found:       true,
				Flow:        flow,
				Query:       query,
				CurrentTurn: 0,
			}
		}

		// Specific matching logic based on keywords
		switch matchedType {
		case "SIM":
			if strings.Contains(queryLower, "buat") || strings.Contains(queryLower, "bikin") {
				if strings.Contains(titleLower, "buat") && strings.Contains(titleLower, "sim") {
					log.Printf("🔍 Found flow: %s (create SIM)", flow.Title)
					return &models.PelayananInfo{
						Found:       true,
						Flow:        flow,
						Query:       query,
						CurrentTurn: 0,
					}
				}
			} else if strings.Contains(queryLower, "perpanjang") {
				if strings.Contains(titleLower, "perpanjangan") && strings.Contains(titleLower, "sim") {
					log.Printf("🔍 Found flow: %s (renew SIM)", flow.Title)
					return &models.PelayananInfo{
						Found:       true,
						Flow:        flow,
						Query:       query,
						CurrentTurn: 0,
					}
				}
			} else if strings.Contains(queryLower, "internasional") {
				if strings.Contains(titleLower, "internasional") {
					log.Printf("🔍 Found flow: %s (international SIM)", flow.Title)
					return &models.PelayananInfo{
						Found:       true,
						Flow:        flow,
						Query:       query,
						CurrentTurn: 0,
					}
				}
			}
		case "STNK":
			if strings.Contains(queryLower, "pajak") {
				if strings.Contains(titleLower, "pajak") {
					log.Printf("🔍 Found flow: %s (vehicle tax)", flow.Title)
					return &models.PelayananInfo{
						Found:       true,
						Flow:        flow,
						Query:       query,
						CurrentTurn: 0,
					}
				}
			} else if strings.Contains(queryLower, "pengesahan") || strings.Contains(queryLower, "5 tahun") {
				if strings.Contains(titleLower, "pengesahan") {
					log.Printf("🔍 Found flow: %s (STNK validation)", flow.Title)
					return &models.PelayananInfo{
						Found:       true,
						Flow:        flow,
						Query:       query,
						CurrentTurn: 0,
					}
				}
			} else if strings.Contains(queryLower, "ganti data") {
				if strings.Contains(titleLower, "ganti data") {
					log.Printf("🔍 Found flow: %s (change STNK data)", flow.Title)
					return &models.PelayananInfo{
						Found:       true,
						Flow:        flow,
						Query:       query,
						CurrentTurn: 0,
					}
				}
			} else if strings.Contains(queryLower, "hilang") {
				if strings.Contains(titleLower, "kehilangan") && strings.Contains(titleLower, "stnk") {
					log.Printf("🔍 Found flow: %s (lost STNK)", flow.Title)
					return &models.PelayananInfo{
						Found:       true,
						Flow:        flow,
						Query:       query,
						CurrentTurn: 0,
					}
				}
			} else if strings.Contains(queryLower, "cek status") {
				if strings.Contains(titleLower, "cek status") {
					log.Printf("🔍 Found flow: %s (check status)", flow.Title)
					return &models.PelayananInfo{
						Found:       true,
						Flow:        flow,
						Query:       query,
						CurrentTurn: 0,
					}
				}
			}
		case "BPN":
			if strings.Contains(titleLower, "balik nama") || strings.Contains(titleLower, "mutasi") {
				log.Printf("🔍 Found flow: %s (vehicle transfer)", flow.Title)
				return &models.PelayananInfo{
					Found:       true,
					Flow:        flow,
					Query:       query,
					CurrentTurn: 0,
				}
			}
		case "Kehilangan":
			if strings.Contains(titleLower, "kehilangan") || strings.Contains(titleLower, "hilang") {
				log.Printf("🔍 Found flow: %s (report lost)", flow.Title)
				return &models.PelayananInfo{
					Found:       true,
					Flow:        flow,
					Query:       query,
					CurrentTurn: 0,
				}
			}
		}
	}

	// If no specific match, try fuzzy matching with all flows
	for _, flow := range s.dataset.Flows {
		titleLower := strings.ToLower(flow.Title)
		words := strings.Fields(titleLower)

		matchCount := 0
		for _, word := range words {
			if strings.Contains(queryLower, word) {
				matchCount++
			}
		}

		// If at least 2 words match, consider it a match
		if matchCount >= 2 {
			log.Printf("🔍 Found flow: %s (fuzzy match)", flow.Title)
			return &models.PelayananInfo{
				Found:       true,
				Flow:        flow,
				Query:       query,
				CurrentTurn: 0,
			}
		}
	}

	log.Printf("📝 No flow found for query: %s", query)
	return &models.PelayananInfo{
		Found: false,
		Query: query,
	}
}

// GetAllFlows returns all available service flows
func (s *PelayananService) GetAllFlows() []models.PelayananFlow {
	return s.dataset.Flows
}
