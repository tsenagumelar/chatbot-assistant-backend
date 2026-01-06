package services

import (
	"encoding/json"
	"log"
	"os"
	"police-assistant-backend/models"
	"strings"
)

type FlowNode struct {
	ID          string                   `json:"id"`
	Type        string                   `json:"type"` // message, question, collect, action
	Text        string                   `json:"text"`
	Choices     []FlowChoice             `json:"choices,omitempty"`
	Collect     *FlowCollect             `json:"collect,omitempty"`
	Action      *FlowAction              `json:"action,omitempty"`
	Transitions []FlowTransition         `json:"transitions"`
	OnSelect    []map[string]interface{} `json:"on_select,omitempty"`
}

type FlowChoice struct {
	ID    string      `json:"id"`
	Label string      `json:"label"`
	Value interface{} `json:"value"`
}

type FlowCollect struct {
	Key  string   `json:"key"`
	Mime []string `json:"mime"`
}

type FlowAction struct {
	Type      string   `json:"type"`
	Template  string   `json:"template,omitempty"`
	Inputs    []string `json:"inputs,omitempty"`
	OutputKey string   `json:"output_key,omitempty"`
	Target    string   `json:"target,omitempty"`
	Reason    string   `json:"reason,omitempty"`
}

type FlowTransition struct {
	When string `json:"when"`
	To   string `json:"to"`
}

type SIMFlow struct {
	FlowID    string              `json:"flow_id"`
	Version   string              `json:"version"`
	Locale    string              `json:"locale"`
	EntryNode string              `json:"entry_node"`
	Nodes     []FlowNode          `json:"nodes"`
	nodeMap   map[string]FlowNode // For quick lookup
}

type SIMFlowService struct {
	flow *SIMFlow
}

func NewSIMFlowService() *SIMFlowService {
	service := &SIMFlowService{}

	if err := service.loadFlow(); err != nil {
		log.Printf("⚠️  Failed to load SIM flow: %v", err)
	} else {
		log.Printf("✅ SIM Flow Service initialized with %d nodes", len(service.flow.Nodes))
	}

	return service
}

func (s *SIMFlowService) loadFlow() error {
	file, err := os.ReadFile("perpanjangan_sim.json")
	if err != nil {
		return err
	}

	var flow SIMFlow
	if err := json.Unmarshal(file, &flow); err != nil {
		return err
	}

	// Build node map for quick lookup
	flow.nodeMap = make(map[string]FlowNode)
	for _, node := range flow.Nodes {
		flow.nodeMap[node.ID] = node
	}

	s.flow = &flow
	return nil
}

// DetectSIMIntent checks if user is asking about SIM services
func (s *SIMFlowService) DetectSIMIntent(message string) bool {
	messageLower := strings.ToLower(message)

	simKeywords := []string{
		"perpanjang sim", "perpanjangan sim", "extend sim",
		"buat sim", "bikin sim", "sim baru",
		"sim a", "sim c",
		"proses sim", "urus sim",
	}

	for _, keyword := range simKeywords {
		if strings.Contains(messageLower, keyword) {
			return true
		}
	}

	return false
}

// GetCurrentNode gets the node for current state
func (s *SIMFlowService) GetCurrentNode(nodeID string) *FlowNode {
	if nodeID == "" {
		nodeID = s.flow.EntryNode
	}

	if node, exists := s.flow.nodeMap[nodeID]; exists {
		return &node
	}

	return nil
}

// FormatNodeResponse formats the node response for chat
func (s *SIMFlowService) FormatNodeResponse(node *FlowNode) string {
	if node == nil {
		return ""
	}

	response := node.Text

	// Add choices if available
	if len(node.Choices) > 0 {
		response += "\n\n"
		for i, choice := range node.Choices {
			response += "\n" + string(rune(i+1)) + ". " + choice.Label
		}
	}

	return response
}

// GetFlowContext returns flow context for AI
func (s *SIMFlowService) GetFlowContext(nodeID string) string {
	if s.flow == nil {
		return ""
	}

	node := s.GetCurrentNode(nodeID)
	if node == nil {
		return ""
	}

	context := `
🔄 FLOW PERPANJANGAN SIM AKTIF
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
User sedang dalam proses flow perpanjangan SIM.

Current Node: ` + nodeID + `
Node Type: ` + node.Type + `

INSTRUKSI KHUSUS UNTUK FLOW SIM:
- Gunakan EXACT text dari flow: "` + node.Text + `"`

	if len(node.Choices) > 0 {
		context += "\n- Tampilkan pilihan dengan format:"
		for i, choice := range node.Choices {
			context += "\n  " + string(rune(i+1)) + ". " + choice.Label
		}
	}

	if node.Type == "collect" {
		context += "\n- Ini adalah step upload dokumen"
		context += "\n- Minta user untuk upload file yang diminta"
		context += "\n- Gunakan emoji 📤 untuk upload"
	}

	context += "\n\n⚠️ PENTING: Ikuti EXACT teks dan urutan dari flow. Jangan tambahkan informasi lain."
	context += "\n━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

	return context
}

// ProcessUserChoice processes user's choice and returns next node
func (s *SIMFlowService) ProcessUserChoice(currentNodeID string, userInput string) (string, *FlowNode) {
	node := s.GetCurrentNode(currentNodeID)
	if node == nil {
		return "", nil
	}

	// Match user input with choices
	userInputLower := strings.ToLower(strings.TrimSpace(userInput))

	for _, choice := range node.Choices {
		choiceLower := strings.ToLower(choice.Label)

		// Check if user input matches choice label or ID
		if strings.Contains(userInputLower, choiceLower) ||
			strings.Contains(userInputLower, choice.ID) ||
			userInputLower == strings.ToLower(choice.Label) {

			// Find transition
			for _, transition := range node.Transitions {
				if strings.Contains(transition.When, choice.ID) {
					nextNode := s.GetCurrentNode(transition.To)
					return transition.To, nextNode
				}
			}
		}
	}

	// Check for number input (1, 2, 3, etc.)
	for i, choice := range node.Choices {
		if userInputLower == string(rune(i+49)) { // '1', '2', '3'
			for _, transition := range node.Transitions {
				if strings.Contains(transition.When, choice.ID) {
					nextNode := s.GetCurrentNode(transition.To)
					return transition.To, nextNode
				}
			}
		}
	}

	return "", nil
}

// GetSIMFlowInfo returns info about SIM flow for context
func (s *SIMFlowService) GetSIMFlowInfo(nodeID string) *models.SIMFlowInfo {
	if s.flow == nil {
		return nil
	}

	node := s.GetCurrentNode(nodeID)
	if node == nil {
		return nil
	}

	info := &models.SIMFlowInfo{
		Active:      true,
		CurrentNode: nodeID,
		NodeType:    node.Type,
		NodeText:    node.Text,
		Choices:     []models.SIMFlowChoice{},
	}

	for _, choice := range node.Choices {
		info.Choices = append(info.Choices, models.SIMFlowChoice{
			ID:    choice.ID,
			Label: choice.Label,
		})
	}

	return info
}
