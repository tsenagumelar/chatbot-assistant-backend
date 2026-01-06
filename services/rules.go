package services

import (
	"encoding/json"
	"log"
	"os"
	"strings"
)

type ResponseRule struct {
	No                        int    `json:"No"`
	JenisPelayanan            string `json:"Jenis Pelayanan"`
	DokumenYangPerluDisiapkan string `json:"Dokumen yang Perlu Disiapkan"`
	Pertanyaan1               string `json:"pertanyaan 1"`
	Response1                 string `json:"response 1"`
	Pertanyaan2               string `json:"pertanyaan 2"`
	Response2                 string `json:"response 2 "`
	Pertanyaan3               string `json:"pertanyaan 3"`
	Response3                 string `json:"response 3"`
}

type ResponseRules struct {
	Sheet1 []ResponseRule `json:"Sheet1"`
}

type LocationRule struct {
	No             int    `json:"No"`
	JenisPelayanan string `json:"Jenis Pelayanan"`
	Input          string `json:"Input"`
	ArahkanKe      string `json:"Arahkan ke"`
}

type RulesService struct {
	responseRules ResponseRules
	locationRules []LocationRule
}

func NewRulesService() *RulesService {
	service := &RulesService{}

	// Load response rules
	if err := service.loadResponseRules(); err != nil {
		log.Printf("⚠️  Failed to load response rules: %v", err)
	} else {
		log.Printf("✅ Response Rules loaded: %d rules", len(service.responseRules.Sheet1))
	}

	// Load location rules
	if err := service.loadLocationRules(); err != nil {
		log.Printf("⚠️  Failed to load location rules: %v", err)
	} else {
		log.Printf("✅ Location Rules loaded: %d rules", len(service.locationRules))
	}

	return service
}

func (s *RulesService) loadResponseRules() error {
	file, err := os.ReadFile("response-rules.json")
	if err != nil {
		return err
	}

	if err := json.Unmarshal(file, &s.responseRules); err != nil {
		return err
	}

	return nil
}

func (s *RulesService) loadLocationRules() error {
	file, err := os.ReadFile("location-rules.json")
	if err != nil {
		return err
	}

	if err := json.Unmarshal(file, &s.locationRules); err != nil {
		return err
	}

	return nil
}

// GetResponseRule mencari response rule berdasarkan jenis pelayanan
func (s *RulesService) GetResponseRule(jenisPelayanan string) *ResponseRule {
	jenisPelayananLower := strings.ToLower(jenisPelayanan)

	for _, rule := range s.responseRules.Sheet1 {
		ruleLower := strings.ToLower(rule.JenisPelayanan)
		if strings.Contains(ruleLower, jenisPelayananLower) || strings.Contains(jenisPelayananLower, ruleLower) {
			return &rule
		}
	}

	return nil
}

// GetLocationRule mencari location rule berdasarkan jenis pelayanan
func (s *RulesService) GetLocationRule(jenisPelayanan string) *LocationRule {
	jenisPelayananLower := strings.ToLower(jenisPelayanan)

	for _, rule := range s.locationRules {
		ruleLower := strings.ToLower(rule.JenisPelayanan)
		if strings.Contains(ruleLower, jenisPelayananLower) || strings.Contains(jenisPelayananLower, ruleLower) {
			return &rule
		}
	}

	return nil
}

// FormatResponseRuleForPrompt mengformat response rule untuk di-inject ke prompt
func (s *RulesService) FormatResponseRuleForPrompt(rule *ResponseRule, userName string, userLocation string) string {
	if rule == nil {
		return ""
	}

	prompt := "📋 ALUR PERCAKAPAN YANG HARUS DIIKUTI:\n"
	prompt += "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n"
	prompt += "⚠️⚠️⚠️ INSTRUKSI WAJIB - IKUTI ALUR INI ⚠️⚠️⚠️\n\n"

	// Replace placeholders
	nameToUse := userName
	if nameToUse == "" {
		nameToUse = "Sobat Lantas"
	}

	locationContext := ""
	if userLocation != "" {
		locationContext = "Saya lihat Anda sedang di " + userLocation
	}

	prompt += "Turn 1 (User bertanya):\n"
	prompt += "  Contoh: \"" + rule.Pertanyaan1 + "\"\n"
	prompt += "  Anda HARUS menjawab:\n"
	response1 := strings.ReplaceAll(rule.Response1, "<name>", nameToUse)
	response1 = strings.ReplaceAll(response1, "<konteks>", locationContext)
	prompt += "  \"" + response1 + "\"\n\n"

	if rule.Pertanyaan2 != "" && rule.Response2 != "" {
		prompt += "Turn 2 (User lanjut):\n"
		prompt += "  Contoh: \"" + rule.Pertanyaan2 + "\"\n"
		prompt += "  Anda HARUS menjawab:\n"
		response2 := strings.ReplaceAll(rule.Response2, "<name>", nameToUse)
		response2 = strings.ReplaceAll(response2, "<konteks>", locationContext)
		prompt += "  \"" + response2 + "\"\n\n"
	}

	if rule.Pertanyaan3 != "" && rule.Response3 != "" {
		prompt += "Turn 3 (User lanjut):\n"
		prompt += "  Contoh: \"" + rule.Pertanyaan3 + "\"\n"
		prompt += "  Anda HARUS menjawab:\n"
		response3 := strings.ReplaceAll(rule.Response3, "<name>", nameToUse)
		response3 = strings.ReplaceAll(response3, "<konteks>", locationContext)
		prompt += "  \"" + response3 + "\"\n\n"
	}

	prompt += "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n\n"
	prompt += "ATURAN PENTING:\n"
	prompt += "1. ✅ IKUTI alur percakapan di atas dengan KETAT\n"
	prompt += "2. ✅ Gunakan TEKS yang sudah ditentukan, jangan improvisasi berlebihan\n"
	prompt += "3. ✅ Ganti <name> dengan \"" + nameToUse + "\"\n"
	prompt += "4. ✅ Ganti <konteks> dengan informasi lokasi user\n"
	prompt += "5. ✅ Track turn/giliran percakapan, jangan skip atau loncat\n"
	prompt += "6. ✅ Jika user upload dokumen, konfirmasi dan lanjut ke turn berikutnya\n"
	prompt += "7. ✅ Tetap ramah dan natural, tapi WAJIB ikuti alur\n\n"

	return prompt
}

// FormatLocationRuleForPrompt mengformat location rule untuk di-inject ke prompt
func (s *RulesService) FormatLocationRuleForPrompt(rule *LocationRule) string {
	if rule == nil {
		return ""
	}

	prompt := "📍 ATURAN LOKASI UNTUK LAYANAN INI:\n"
	prompt += "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n"
	prompt += "Jenis Input: " + rule.Input + "\n"
	prompt += "Instruksi: " + rule.ArahkanKe + "\n"
	prompt += "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n\n"

	return prompt
}
