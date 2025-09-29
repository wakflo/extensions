package actions

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Error message constants
const (
	BillingIssueMessage = `Error Occurred: 429 - Rate Limit / Billing Issue
1. Ensure that billing is enabled on your OpenAI platform.
2. Generate a new API key.
3. Attempt the process again.
For guidance, visit: https://platform.openai.com/account/billing`

	UnauthorizedMessage = `Error Occurred: 401 - Unauthorized
Ensure that your API key is valid.
For guidance, visit: https://platform.openai.com/api-keys`

	InvalidModelMessage = `Error Occurred: 404 - Model Not Found
The specified model does not exist or you don't have access to it.
Please check the model name and try again.`

	BadRequestMessage = `Error Occurred: 400 - Bad Request
The request was invalid. Please check your input parameters.`

	ServerErrorMessage = `Error Occurred: 500 - Server Error
OpenAI is experiencing issues. Please try again later.`

	TimeoutMessage = `Error Occurred: Request Timeout
The request took too long to process. Please try again with a smaller input or simpler request.`
)

// OpenAIErrorResponse represents the error structure from OpenAI API
type OpenAIErrorResponse struct {
	Error struct {
		Message string `json:"message"`
		Type    string `json:"type"`
		Param   string `json:"param,omitempty"`
		Code    string `json:"code,omitempty"`
	} `json:"error"`
}

// ParseOpenAIError parses the OpenAI error response and returns a user-friendly message
func ParseOpenAIError(statusCode int, responseBody []byte) error {
	var openAIError OpenAIErrorResponse

	// Try to parse the OpenAI error structure
	if err := json.Unmarshal(responseBody, &openAIError); err == nil && openAIError.Error.Message != "" {
		// Create detailed error based on status code
		switch statusCode {
		case 401:
			return fmt.Errorf("%s\nDetails: %s", UnauthorizedMessage, openAIError.Error.Message)
		case 429:
			// Check if it's a rate limit or quota issue
			if strings.Contains(strings.ToLower(openAIError.Error.Message), "quota") ||
				strings.Contains(strings.ToLower(openAIError.Error.Message), "billing") ||
				strings.Contains(strings.ToLower(openAIError.Error.Message), "exceeded") {
				return fmt.Errorf("%s\nDetails: %s", BillingIssueMessage, openAIError.Error.Message)
			}
			return fmt.Errorf("Rate limit exceeded. Please wait and try again.\nDetails: %s", openAIError.Error.Message)
		case 404:
			return fmt.Errorf("%s\nDetails: %s", InvalidModelMessage, openAIError.Error.Message)
		case 400:
			// Check for specific parameter errors
			if openAIError.Error.Param != "" {
				return fmt.Errorf("%s\nParameter: %s\nDetails: %s",
					BadRequestMessage, openAIError.Error.Param, openAIError.Error.Message)
			}
			return fmt.Errorf("%s\nDetails: %s", BadRequestMessage, openAIError.Error.Message)
		case 500, 502, 503:
			return fmt.Errorf("%s\nDetails: %s", ServerErrorMessage, openAIError.Error.Message)
		default:
			return fmt.Errorf("OpenAI API Error (Status %d): %s", statusCode, openAIError.Error.Message)
		}
	}

	// Fallback to status-based messages if we can't parse the error
	switch statusCode {
	case 401:
		return fmt.Errorf(UnauthorizedMessage)
	case 429:
		return fmt.Errorf(BillingIssueMessage)
	case 404:
		return fmt.Errorf(InvalidModelMessage)
	case 400:
		return fmt.Errorf(BadRequestMessage)
	case 500, 502, 503:
		return fmt.Errorf(ServerErrorMessage)
	default:
		return fmt.Errorf("OpenAI API Error (Status %d): %s", statusCode, string(responseBody))
	}
}

// HandleOpenAIResponse processes the response and returns appropriate errors
func HandleOpenAIResponse(res interface{}, hasStatus func() interface{}, hasBody func() interface{}) error {
	// This is a generic handler that works with your fastshot response
	// You'll need to adapt this based on your actual response type

	// Example implementation (adjust based on your actual types):
	// if res.Status().IsError() {
	//     bodyBytes, _ := io.ReadAll(res.Body().Raw())
	//     return ParseOpenAIError(res.Status().Code(), bodyBytes)
	// }

	return nil
}
