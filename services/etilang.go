package services

import (
	"log"
	"police-assistant-backend/models"
	"strings"
)

type ETilangService struct {
	// Data dummy untuk testing
	dummyData map[string]*models.ETilangInfo
}

func NewETilangService() *ETilangService {
	service := &ETilangService{
		dummyData: make(map[string]*models.ETilangInfo),
	}

	// Initialize dummy data
	service.initializeDummyData()

	log.Println("✅ E-Tilang Service initialized with dummy data")
	return service
}

func (s *ETilangService) initializeDummyData() {
	// Data dummy 1: Ada pelanggaran
	s.dummyData["B1234SV"] = &models.ETilangInfo{
		PlateNumber:   "B 1234 SV",
		ChassisNumber: "MH1RP6701FK123456",
		OwnerName:     "Budi Santoso",
		VehicleType:   "Motor Honda Beat",
		HasViolation:  true,
		Violations: []models.ETilangViolation{
			{
				Date:        "2025-12-15",
				Violation:   "Melanggar lampu merah",
				Location:    "Jl. Sudirman - Jakarta Pusat",
				Fine:        500000,
				OfficerName: "Brigadir Joko Widodo",
				Status:      "unpaid",
			},
			{
				Date:        "2025-12-20",
				Violation:   "Tidak menggunakan helm SNI",
				Location:    "Jl. Gatot Subroto - Jakarta Selatan",
				Fine:        250000,
				OfficerName: "Aipda Siti Nurhaliza",
				Status:      "unpaid",
			},
		},
		TotalFine: 750000,
	}

	// Data dummy 2: Ada pelanggaran parkir
	s.dummyData["B5678XY"] = &models.ETilangInfo{
		PlateNumber:   "B 5678 XY",
		ChassisNumber: "MHKA42BA7JK098765",
		OwnerName:     "Siti Rahayu",
		VehicleType:   "Mobil Toyota Avanza",
		HasViolation:  true,
		Violations: []models.ETilangViolation{
			{
				Date:        "2026-01-02",
				Violation:   "Parkir di tempat terlarang",
				Location:    "Jl. MH Thamrin - Jakarta Pusat",
				Fine:        300000,
				OfficerName: "Bripka Ahmad Dahlan",
				Status:      "unpaid",
			},
		},
		TotalFine: 300000,
	}

	// Data dummy 3: Tidak ada pelanggaran
	s.dummyData["B9999ZZ"] = &models.ETilangInfo{
		PlateNumber:   "B 9999 ZZ",
		ChassisNumber: "MH1JC5101FK234567",
		OwnerName:     "Ahmad Fauzi",
		VehicleType:   "Motor Yamaha NMAX",
		HasViolation:  false,
		Violations:    []models.ETilangViolation{},
		TotalFine:     0,
	}

	// Data dummy 4: Pelanggaran kecepatan
	s.dummyData["D1111AA"] = &models.ETilangInfo{
		PlateNumber:   "D 1111 AA",
		ChassisNumber: "MHRGN81235K876543",
		OwnerName:     "Rina Kartika",
		VehicleType:   "Mobil Honda CR-V",
		HasViolation:  true,
		Violations: []models.ETilangViolation{
			{
				Date:        "2025-12-28",
				Violation:   "Melebihi batas kecepatan (120 km/jam di tol)",
				Location:    "Tol Jagorawi KM 15",
				Fine:        500000,
				OfficerName: "Aiptu Bambang Suryono",
				Status:      "paid",
			},
		},
		TotalFine: 500000,
	}

	// Data dummy 5: Pelanggaran penggunaan HP
	s.dummyData["E7777BB"] = &models.ETilangInfo{
		PlateNumber:   "E 7777 BB",
		ChassisNumber: "MH1JFJ110FK345678",
		OwnerName:     "Dedi Gunawan",
		VehicleType:   "Motor Kawasaki Ninja",
		HasViolation:  true,
		Violations: []models.ETilangViolation{
			{
				Date:        "2026-01-05",
				Violation:   "Menggunakan handphone saat berkendara",
				Location:    "Jl. Asia Afrika - Bandung",
				Fine:        750000,
				OfficerName: "Brigadir Eka Prasetya",
				Status:      "unpaid",
			},
		},
		TotalFine: 750000,
	}
}

// CheckETilang checks e-tilang by plate number
func (s *ETilangService) CheckETilang(plateNumber string) *models.ETilangInfo {
	// Normalize plate number (remove spaces, uppercase)
	normalized := strings.ToUpper(strings.ReplaceAll(plateNumber, " ", ""))

	log.Printf("🔍 Checking E-Tilang for plate: %s (normalized: %s)", plateNumber, normalized)

	// Check in dummy data
	if info, exists := s.dummyData[normalized]; exists {
		log.Printf("✅ E-Tilang data found for %s", plateNumber)
		return info
	}

	// If not found, return no violation
	log.Printf("📝 No E-Tilang data found for %s, returning clean record", plateNumber)
	return &models.ETilangInfo{
		PlateNumber:   plateNumber,
		ChassisNumber: "XXXXXXXXXXXX",
		OwnerName:     "-",
		VehicleType:   "-",
		HasViolation:  false,
		Violations:    []models.ETilangViolation{},
		TotalFine:     0,
	}
}

// ExtractPlateNumber tries to extract plate number from user message
func (s *ETilangService) ExtractPlateNumber(message string) string {
	// Simple extraction logic
	// Look for patterns like: B1234SV, B 1234 SV, etc.
	message = strings.ToUpper(message)

	// Common patterns for Indonesian plate numbers
	// Format: [Letter(s)] [Numbers] [Letter(s)]
	words := strings.Fields(message)

	for i, word := range words {
		// Remove common punctuation
		word = strings.Trim(word, ".,!?;:")

		// Check if word looks like plate number (has letters and numbers)
		hasLetter := false
		hasNumber := false

		for _, char := range word {
			if char >= 'A' && char <= 'Z' {
				hasLetter = true
			}
			if char >= '0' && char <= '9' {
				hasNumber = true
			}
		}

		// If this word has both letters and numbers, might be plate number
		if hasLetter && hasNumber {
			// Try to combine with next words if they exist (for spaced plate numbers)
			candidate := word
			if i+1 < len(words) {
				nextWord := strings.Trim(words[i+1], ".,!?;:")
				// Check if next word is numbers
				isNumber := true
				for _, char := range nextWord {
					if char < '0' || char > '9' {
						isNumber = false
						break
					}
				}
				if isNumber && i+2 < len(words) {
					thirdWord := strings.Trim(words[i+2], ".,!?;:")
					// Check if third word is letters
					isLetter := true
					for _, char := range thirdWord {
						if char < 'A' || char > 'Z' {
							isLetter = false
							break
						}
					}
					if isLetter {
						candidate = word + nextWord + thirdWord
					}
				}
			}

			// Validate length (Indonesian plate numbers are typically 5-8 chars without spaces)
			normalized := strings.ReplaceAll(candidate, " ", "")
			if len(normalized) >= 5 && len(normalized) <= 10 {
				return normalized
			}
		}
	}

	return ""
}
