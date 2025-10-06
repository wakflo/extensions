package shared

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/juicycleff/smartform/v1"
	sdkcontext "github.com/wakflo/go-sdk/v2/context"
)

func CreateClaudeClient(ctx context.Context, auth *sdkcontext.AuthContext) (*http.Client, string, error) {
	return &http.Client{}, auth.Extra["apiKey"], nil
}

func CallClaudeAPI(ctx context.Context, auth *sdkcontext.AuthContext, request ClaudeRequest) (*ClaudeResponse, error) {
	client, apiKey, err := CreateClaudeClient(ctx, auth)
	if err != nil {
		return nil, err
	}

	jsonData, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", ClaudeAPIURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", apiKey)
	req.Header.Set("anthropic-version", ClaudeVersion)

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		var errorResp ClaudeErrorResponse
		if err := json.Unmarshal(body, &errorResp); err != nil {
			return nil, fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(body))
		}
		return nil, fmt.Errorf("API error: %s - %s", errorResp.Error.Type, errorResp.Error.Message)
	}

	var claudeResp ClaudeResponse
	if err := json.Unmarshal(body, &claudeResp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &claudeResp, nil
}

func RegisterModelProps(form *smartform.FormBuilder) *smartform.FieldBuilder {
	return form.SelectField("model", "Model").
		Placeholder("Select a Claude model").
		Required(true).
		AddOptions(
			// Claude 4 Family (Latest Generation)
			smartform.NewOption("claude-sonnet-4-5-20250929", "Claude Sonnet 4.5 (Best for Agents & Coding)"),
			smartform.NewOption("claude-sonnet-4-5", "Claude Sonnet 4.5 (Alias)"),
			smartform.NewOption("claude-opus-4-1-20250805", "Claude Opus 4.1 (Specialized Complex Tasks)"),
			smartform.NewOption("claude-opus-4-1", "Claude Opus 4.1 (Alias)"),
			smartform.NewOption("claude-sonnet-4-20250514", "Claude Sonnet 4 (High Performance)"),
			smartform.NewOption("claude-sonnet-4-0", "Claude Sonnet 4 (Alias)"),
			smartform.NewOption("claude-opus-4-20250514", "Claude Opus 4 (Previous Flagship)"),
			smartform.NewOption("claude-opus-4-0", "Claude Opus 4 (Alias)"),
			
			// Claude 3.7 Family
			smartform.NewOption("claude-3-7-sonnet-20250219", "Claude Sonnet 3.7 (Extended Thinking)"),
			smartform.NewOption("claude-3-7-sonnet-latest", "Claude Sonnet 3.7 (Alias)"),
			
			// Claude 3.5 Family
			smartform.NewOption("claude-3-5-haiku-20241022", "Claude Haiku 3.5 (Fastest)"),
			smartform.NewOption("claude-3-5-haiku-latest", "Claude Haiku 3.5 (Alias)"),
			
			// Claude 3 Family
			smartform.NewOption("claude-3-haiku-20240307", "Claude Haiku 3 (Fast & Compact)"),
			
			// Legacy Models (if still supported)
			smartform.NewOption("claude-2.1", "Claude 2.1 (Legacy)"),
			smartform.NewOption("claude-2.0", "Claude 2.0 (Legacy)"),
			smartform.NewOption("claude-instant-1.2", "Claude Instant 1.2 (Legacy)"),
		).
		HelpText("Claude Sonnet 4.5 offers the best performance for complex agents and coding. Use aliases for automatic updates to the latest snapshot.")
}

// ExtractResponseText extracts text from Claude response
func ExtractResponseText(response *ClaudeResponse) string {
	var text string
	for _, content := range response.Content {
		if content.Type == "text" {
			text += content.Text
		}
	}
	return text
}

func getLanguageName(code string) string {
	languages := map[string]string{
		"en": "English", "es": "Spanish", "fr": "French",
		"de": "German", "it": "Italian", "pt": "Portuguese",
		"ru": "Russian", "ja": "Japanese", "ko": "Korean",
		"zh": "Chinese", "ar": "Arabic", "hi": "Hindi",
		"nl": "Dutch", "pl": "Polish", "tr": "Turkish",
		"auto": "Auto-detect",
	}
	if name, ok := languages[code]; ok {
		return name
	}
	return code
}
