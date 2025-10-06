package shared

const (
	ClaudeAPIURL  = "https://api.anthropic.com/v1/messages"
	ClaudeVersion = "2023-06-01"
)

// ClaudeMessage represents a message in the conversation
type ClaudeMessage struct {
	Role    string        `json:"role"`
	Content []interface{} `json:"content"`
}

// ClaudeTextContent represents text content
type ClaudeTextContent struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

// ClaudeImageContent represents image content
type ClaudeImageContent struct {
	Type   string            `json:"type"`
	Source ClaudeImageSource `json:"source"`
}

// ClaudeImageSource represents the image source
type ClaudeImageSource struct {
	Type      string `json:"type"`
	MediaType string `json:"media_type"`
	Data      string `json:"data"` // base64 encoded
}

// ClaudeRequest represents the API request
type ClaudeRequest struct {
	Model       string          `json:"model"`
	Messages    []ClaudeMessage `json:"messages"`
	MaxTokens   int             `json:"max_tokens"`
	Temperature float64         `json:"temperature,omitempty"`
	System      string          `json:"system,omitempty"`
}

// ClaudeResponse represents the API response
type ClaudeResponse struct {
	ID      string `json:"id"`
	Type    string `json:"type"`
	Role    string `json:"role"`
	Content []struct {
		Type string `json:"type"`
		Text string `json:"text"`
	} `json:"content"`
	Model        string `json:"model"`
	StopReason   string `json:"stop_reason"`
	StopSequence string `json:"stop_sequence"`
	Usage        struct {
		InputTokens  int `json:"input_tokens"`
		OutputTokens int `json:"output_tokens"`
	} `json:"usage"`
}

// ClaudeErrorResponse represents an error response
type ClaudeErrorResponse struct {
	Type  string `json:"type"`
	Error struct {
		Type    string `json:"type"`
		Message string `json:"message"`
	} `json:"error"`
}
