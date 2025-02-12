package actions

type ModelResponse struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	Created int64  `json:"created"`
	OwnedBy string `json:"owned_by"`
}

type ModelListResponse struct {
	Object string          `json:"object"`
	Data   []ModelResponse `json:"data"`
}

type ChatCompletionUsageResponse struct {
	CompletionTokens int64 `json:"completion_tokens"`
	PromptTokens     int64 `json:"prompt_tokens"`
	TotalTokens      int64 `json:"total_tokens"`
}

type ChatCompletionMessageResponse struct {
	Content string `json:"content"`
	Role    string `json:"role"`
}

type ChatCompletionChoiceResponse struct {
	FinishReason string                        `json:"finish_reason"`
	Index        int64                         `json:"index"`
	Message      ChatCompletionMessageResponse `json:"message"`
	Logprobs     interface{}                   `json:"logprobs"`
}

type ChatCompletionResponse struct {
	ID                string                         `json:"id"`
	Object            string                         `json:"object"`
	Created           int64                          `json:"created"`
	Model             string                         `json:"model"`
	SystemFingerprint string                         `json:"system_fingerprint"`
	Choices           []ChatCompletionChoiceResponse `json:"choices"`
	Usage             ChatCompletionUsageResponse    `json:"usage"`
}

type GeneratedImageResponse struct {
	URL string `json:"url"`
}

type ImageGenerationResponse struct {
	Created int64                    `json:"created"`
	Data    []GeneratedImageResponse `json:"data"`
}
