package models

// Request & Response structures for Chat
type ChatRequest struct {
	Message   string          `json:"message" validate:"required"`
	Context   Context         `json:"context"`
	SessionID string          `json:"session_id,omitempty"` // Session ID untuk backend-managed history
	History   []OpenAIMessage `json:"history,omitempty"`    // Optional: untuk backward compatibility
}

type Context struct {
	Location  string  `json:"location"`
	Speed     float64 `json:"speed"`
	Traffic   string  `json:"traffic"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type ChatResponse struct {
	Success   bool   `json:"success"`
	Response  string `json:"response"`
	SessionID string `json:"session_id,omitempty"` // Return session ID ke frontend
	Error     string `json:"error,omitempty"`
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
	Model       string          `json:"model"`
	Messages    []OpenAIMessage `json:"messages"`
	Temperature float64         `json:"temperature,omitempty"`
	MaxTokens   int             `json:"max_tokens,omitempty"`
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
