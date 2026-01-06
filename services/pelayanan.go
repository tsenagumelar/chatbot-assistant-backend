package services

import (
	"encoding/json"
	"log"
	"os"
	"police-assistant-backend/models"
	"strings"
)

type PelayananService struct {
	data []models.PelayananData
}

func NewPelayananService() *PelayananService {
	service := &PelayananService{
		data: []models.PelayananData{},
	}

	// Load data from JSON file
	if err := service.loadData(); err != nil {
		log.Printf("⚠️  Failed to load pelayanan data: %v", err)
	} else {
		log.Printf("✅ Pelayanan Service initialized with %d services", len(service.data))
	}

	return service
}

func (s *PelayananService) loadData() error {
	// Read JSON file
	file, err := os.ReadFile("data_pelayanan.json")
	if err != nil {
		return err
	}

	// Parse JSON
	if err := json.Unmarshal(file, &s.data); err != nil {
		return err
	}

	return nil
}

// SearchPelayanan searches for service information based on user query
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

	// Search in data based on matched type and specific keywords
	for _, pelayanan := range s.data {
		jenisPelayanan := strings.ToLower(pelayanan.JenisPelayanan)
		jenisPelayananReadable := strings.ReplaceAll(jenisPelayanan, "_", " ")

		// Direct match with service name
		if strings.Contains(jenisPelayananReadable, queryLower) {
			log.Printf("🔍 Found pelayanan: %s (direct match)", pelayanan.JenisPelayanan)
			return &models.PelayananInfo{
				Found:     true,
				Pelayanan: pelayanan,
				Query:     query,
			}
		}

		// Specific matching logic
		switch matchedType {
		case "SIM":
			if strings.Contains(queryLower, "buat") || strings.Contains(queryLower, "bikin") {
				if strings.Contains(jenisPelayanan, "buat_sim") {
					log.Printf("🔍 Found pelayanan: %s (create SIM)", pelayanan.JenisPelayanan)
					return &models.PelayananInfo{
						Found:     true,
						Pelayanan: pelayanan,
						Query:     query,
					}
				}
			} else if strings.Contains(queryLower, "perpanjang") {
				if strings.Contains(jenisPelayanan, "perpanjangan_sim") {
					log.Printf("🔍 Found pelayanan: %s (renew SIM)", pelayanan.JenisPelayanan)
					return &models.PelayananInfo{
						Found:     true,
						Pelayanan: pelayanan,
						Query:     query,
					}
				}
			} else if strings.Contains(queryLower, "hilang") || strings.Contains(queryLower, "rusak") {
				if strings.Contains(jenisPelayanan, "sim_hilang") {
					log.Printf("🔍 Found pelayanan: %s (lost/damaged SIM)", pelayanan.JenisPelayanan)
					return &models.PelayananInfo{
						Found:     true,
						Pelayanan: pelayanan,
						Query:     query,
					}
				}
			} else if strings.Contains(queryLower, "internasional") {
				if strings.Contains(jenisPelayanan, "sim_internasional") {
					log.Printf("🔍 Found pelayanan: %s (international SIM)", pelayanan.JenisPelayanan)
					return &models.PelayananInfo{
						Found:     true,
						Pelayanan: pelayanan,
						Query:     query,
					}
				}
			}
		case "STNK":
			if strings.Contains(queryLower, "pajak") {
				if strings.Contains(jenisPelayanan, "pajak_kendaraan") {
					log.Printf("🔍 Found pelayanan: %s (vehicle tax)", pelayanan.JenisPelayanan)
					return &models.PelayananInfo{
						Found:     true,
						Pelayanan: pelayanan,
						Query:     query,
					}
				}
			} else if strings.Contains(queryLower, "pengesahan") || strings.Contains(queryLower, "5 tahun") {
				if strings.Contains(jenisPelayanan, "pengesahan_stnk") {
					log.Printf("🔍 Found pelayanan: %s (STNK validation)", pelayanan.JenisPelayanan)
					return &models.PelayananInfo{
						Found:     true,
						Pelayanan: pelayanan,
						Query:     query,
					}
				}
			} else if strings.Contains(queryLower, "ganti data") {
				if strings.Contains(jenisPelayanan, "ganti_data_stnk") {
					log.Printf("🔍 Found pelayanan: %s (change STNK data)", pelayanan.JenisPelayanan)
					return &models.PelayananInfo{
						Found:     true,
						Pelayanan: pelayanan,
						Query:     query,
					}
				}
			} else if strings.Contains(queryLower, "hilang") {
				if strings.Contains(jenisPelayanan, "laporan_kehilangan_stnk") {
					log.Printf("🔍 Found pelayanan: %s (lost STNK)", pelayanan.JenisPelayanan)
					return &models.PelayananInfo{
						Found:     true,
						Pelayanan: pelayanan,
						Query:     query,
					}
				}
			} else if strings.Contains(queryLower, "cek status") {
				if strings.Contains(jenisPelayanan, "cek_status") {
					log.Printf("🔍 Found pelayanan: %s (check status)", pelayanan.JenisPelayanan)
					return &models.PelayananInfo{
						Found:     true,
						Pelayanan: pelayanan,
						Query:     query,
					}
				}
			}
		case "Tilang":
			if strings.Contains(jenisPelayanan, "cek_tilang") || strings.Contains(jenisPelayanan, "etle") {
				log.Printf("🔍 Found pelayanan: %s (e-tilang check)", pelayanan.JenisPelayanan)
				return &models.PelayananInfo{
					Found:     true,
					Pelayanan: pelayanan,
					Query:     query,
				}
			}
		case "BPN":
			if strings.Contains(jenisPelayanan, "balik_nama") {
				log.Printf("🔍 Found pelayanan: %s (vehicle transfer)", pelayanan.JenisPelayanan)
				return &models.PelayananInfo{
					Found:     true,
					Pelayanan: pelayanan,
					Query:     query,
				}
			} else if strings.Contains(jenisPelayanan, "mutasi") {
				log.Printf("🔍 Found pelayanan: %s (vehicle mutation)", pelayanan.JenisPelayanan)
				return &models.PelayananInfo{
					Found:     true,
					Pelayanan: pelayanan,
					Query:     query,
				}
			}
		}
	}

	// If no specific match, try fuzzy matching with all services
	for _, pelayanan := range s.data {
		jenisPelayanan := strings.ToLower(pelayanan.JenisPelayanan)
		parts := strings.Split(jenisPelayanan, "_")

		matchCount := 0
		for _, part := range parts {
			if strings.Contains(queryLower, part) {
				matchCount++
			}
		}

		// If at least 2 parts match, consider it a match
		if matchCount >= 2 {
			log.Printf("🔍 Found pelayanan: %s (fuzzy match)", pelayanan.JenisPelayanan)
			return &models.PelayananInfo{
				Found:     true,
				Pelayanan: pelayanan,
				Query:     query,
			}
		}
	}

	log.Printf("📝 No pelayanan found for query: %s", query)
	return &models.PelayananInfo{
		Found: false,
		Query: query,
	}
}

// GetAllPelayanan returns all available services
func (s *PelayananService) GetAllPelayanan() []models.PelayananData {
	return s.data
}
