package models

// Request & Response structures for Chat
type ChatRequest struct {
	Message   string             `json:"message" validate:"required"`
	Context   Context            `json:"context"`
	SessionID string             `json:"session_id,omitempty"` // Session ID untuk backend-managed history
	History   []OpenAIMessage    `json:"history,omitempty"`    // Optional: untuk backward compatibility
	Documents []UploadedDocument `json:"documents,omitempty"`  // Dokumen yang diupload (base64 atau URL)
}

type Context struct {
	Location              string         `json:"location"`
	Speed                 float64        `json:"speed"`
	Traffic               string         `json:"traffic"`
	Latitude              float64        `json:"latitude"`
	Longitude             float64        `json:"longitude"`
	ETilangInfo           *ETilangInfo   `json:"e_tilang_info,omitempty"`  // Info tilang jika dicek
	PelayananInfo         *PelayananInfo `json:"pelayanan_info,omitempty"` // Info pelayanan jika ditanyakan
	HasUploadedDocuments  bool           `json:"has_uploaded_documents"`   // Flag jika user upload dokumen
	UploadedDocumentCount int            `json:"uploaded_document_count"`  // Jumlah dokumen yang diupload
}

type ChatResponse struct {
	Success       bool           `json:"success"`
	Response      string         `json:"response"`
	SessionID     string         `json:"session_id,omitempty"`     // Return session ID ke frontend
	ETilangInfo   *ETilangInfo   `json:"e_tilang_info,omitempty"`  // Info tilang jika dicek
	PelayananInfo *PelayananInfo `json:"pelayanan_info,omitempty"` // Info pelayanan jika ditanyakan
	Error         string         `json:"error,omitempty"`
}

// Session structures
type SessionResponse struct {
	Success   bool   `json:"success"`
	SessionID string `json:"session_id"`
	Message   string `json:"message,omitempty"`
	Error     string `json:"error,omitempty"`
}

// Traffic structures
type TrafficRequest struct {
	Latitude  float64 `json:"latitude" validate:"required"`
	Longitude float64 `json:"longitude" validate:"required"`
}

type TrafficResponse struct {
	Success bool                   `json:"success"`
	Traffic map[string]interface{} `json:"traffic"`
	Error   string                 `json:"error,omitempty"`
}

// Route structures
type RouteRequest struct {
	Origin      string `json:"origin" validate:"required"`
	Destination string `json:"destination" validate:"required"`
}

type RouteResponse struct {
	Success bool                     `json:"success"`
	Routes  []map[string]interface{} `json:"routes"`
	Error   string                   `json:"error,omitempty"`
}

// OpenAI API structures
type OpenAIMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type OpenAIRequest struct {
	Model               string          `json:"model"`
	Messages            []OpenAIMessage `json:"messages"`
	Temperature         float64         `json:"temperature,omitempty"`
	MaxCompletionTokens int             `json:"max_completion_tokens,omitempty"`
}

type OpenAIResponse struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	Created int64  `json:"created"`
	Model   string `json:"model"`
	Choices []struct {
		Index   int `json:"index"`
		Message struct {
			Role    string `json:"role"`
			Content string `json:"content"`
		} `json:"message"`
		FinishReason string `json:"finish_reason"`
	} `json:"choices"`
	Usage struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens      int `json:"total_tokens"`
	} `json:"usage"`
}

// E-Tilang structures
type ETilangViolation struct {
	Date        string `json:"date"`
	Violation   string `json:"violation"`
	Location    string `json:"location"`
	Fine        int    `json:"fine"`
	OfficerName string `json:"officer_name"`
	Status      string `json:"status"` // "unpaid", "paid", "processed"
}

type ETilangInfo struct {
	PlateNumber   string             `json:"plate_number"`
	ChassisNumber string             `json:"chassis_number"`
	OwnerName     string             `json:"owner_name"`
	VehicleType   string             `json:"vehicle_type"`
	HasViolation  bool               `json:"has_violation"`
	Violations    []ETilangViolation `json:"violations,omitempty"`
	TotalFine     int                `json:"total_fine"`
}

// Pelayanan structures
type PelayananData struct {
	No                        int      `json:"no"`
	JenisPelayanan            string   `json:"jenis_pelayanan"`
	DokumenYangPerluDisiapkan []string `json:"dokumen_yang_perlu_disiapkan"`
}

type PelayananInfo struct {
	Found     bool          `json:"found"`
	Pelayanan PelayananData `json:"pelayanan,omitempty"`
	Query     string        `json:"query"`
}

// Document upload structures
type UploadedDocument struct {
	FileName    string `json:"file_name"`
	FileType    string `json:"file_type"`   // "image/jpeg", "image/png", "application/pdf", etc.
	Base64Data  string `json:"base64_data"` // Base64 encoded file data
	URL         string `json:"url"`         // Or URL if file is hosted elsewhere
	UploadedAt  string `json:"uploaded_at"` // Timestamp
	Description string `json:"description"` // Optional: "KTP", "Surat Kehilangan", etc.
}
