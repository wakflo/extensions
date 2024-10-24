package openai

import (
	"encoding/json"
	"errors"
	"io"

	"github.com/wakflo/go-sdk/autoform"
	sdk "github.com/wakflo/go-sdk/connector"
	sdkcore "github.com/wakflo/go-sdk/core"
)

type PromptChatGPTOperation struct {
	options *sdk.OperationInfo
}

type PromptChatGPTOperationInput struct {
	Model            string   `json:"model"`
	Prompt           string   `json:"prompt"`
	MaxTokens        *int     `json:"max_tokens,omitempty"`
	FrequencyPenalty *float64 `json:"frequency_penalty,omitempty"` // -2.0 to 2.0
	PresencePenalty  *float64 `json:"presence_penalty,omitempty"`  // -2.0 to 2.0
	Seed             *int     `json:"seed,omitempty"`              // random integer
	Temperature      *float64 `json:"temperature,omitempty"`       // 0 to 2.0
	TopP             *float64 `json:"top_p,omitempty"`             // not said, but it's an alternative to Temperature, so probably 0 to 2.0
}

func NewPromptChatGPT() *PromptChatGPTOperation {
	getGPTModels := func(ctx *sdkcore.DynamicFieldContext) (interface{}, error) {
		models, err := getModels(ctx.Auth.Secret, "gpt")
		if err != nil {
			return nil, err
		}

		return models, nil
	}

	return &PromptChatGPTOperation{
		options: &sdk.OperationInfo{
			Name:        "Prompt ChatGPT",
			Description: "Ask something to ChatGPT",
			RequireAuth: false,
			Auth:        autoform.NewAuthSecretField().SetDisplayName("OpenAI API key").SetDescription("Your OpenAI api key").Build(),
			Input: map[string]*sdkcore.AutoFormSchema{
				"model": autoform.NewDynamicField(sdkcore.String).
					SetDisplayName("Model").
					SetDescription("Choose model which you will interact").
					SetDynamicOptions(&getGPTModels).
					SetDependsOn([]string{"connection"}).
					SetRequired(true).
					Build(),
				"prompt": autoform.NewLongTextField().
					SetDisplayName("Prompt").
					SetDescription("What would you like to ask ChatGPT?").
					SetRequired(true).
					Build(),
				"max_tokens": autoform.NewNumberField().
					SetDisplayName("Max Tokens").
					SetDisplayName("Max tokens that ChatGPT can use to generate it's answer.").
					SetMinimum(1).
					//nolint:mnd
					SetMaximum(4096).
					Build(),
				"frequency_penalty": autoform.NewNumberField().
					SetDisplayName("Frequency penalty").
					SetDescription("Penalize or not new tokens (letters/words) based on their frequency. Positive values penalize same verbatim while negative values don't (-2.0 to 2.0)").
					SetMinimum(-2.0).
					//nolint:mnd
					SetMaximum(2.0).
					Build(),
				"presence_penalty": autoform.NewNumberField().
					SetDisplayName("Penalize or not new tokens (letters/words) based on their appearance. Positive values will increase the chance of new topics while negative values don't (-2.0 to 2.0)").
					SetDescription("Presence penalty").
					SetMinimum(-2.0).
					//nolint:mnd
					SetMaximum(2.0).
					Build(),
				"seed": autoform.NewNumberField().
					SetDisplayName("Random seed").
					SetDescription("By using a random seed the model will try its best to replicate the answer that it gave before with the same seed").
					Build(),
				"temperature": autoform.NewNumberField().
					SetDisplayName("Temperature").
					SetDescription("Control model randomness with higher values and more focused and deterministic with lowest values (0 to 2)").
					//nolint:mnd
					SetMinimum(.0).
					//nolint:mnd
					SetMaximum(2.0).
					Build(),
				"top_p": autoform.NewNumberField().
					SetDisplayName("Top P").
					SetDescription("An alternative to temperature that considers the probability of the tokens appearing. Lower values consider low probabilities (0 to 2). It's advised to only use this or temperature and not both.").
					SetMinimum(0.0).
					//nolint:mnd
					SetMaximum(2.0).
					Build(),
			},
			SampleOutput: map[string]interface{}{
				"name":       "openai-prompt-chatgpt",
				"usage_mode": "operation",
				"prompt":     "Who won the Oscar for best actor in 2010?",
				"max_tokens": "500",
				"gpt_answer": "The Oscar for Best Actor at the 82nd Academy Awards, held in 2010, was won by Jeff Bridges for his role as Otis \"Bad\" Blake in the film \"Crazy Heart.\"",
			},
			ErrorSettings: sdkcore.StepErrorSettings{
				ContinueOnError: false,
				RetryOnError:    false,
			},
		},
	}
}

func (c *PromptChatGPTOperation) Run(ctx *sdk.RunContext) (sdk.JSON, error) {
	inputPt := sdk.InputToType[PromptChatGPTOperationInput](ctx)

	// Validate inputs
	if err := validateInput(inputPt); err != nil {
		return nil, err
	}

	// Build request body
	requestBody := buildRequestBody(inputPt)

	client, err := getOpenAiClient(ctx.Auth.Secret)
	if err != nil {
		return nil, err
	}

	res, err := client.POST("/chat/completions").
		Header().AddContentType("application/json").
		Body().AsJSON(requestBody).
		Send()
	if err != nil {
		return nil, err
	}

	if res.Status().IsError() {
		return nil, errors.New(res.Status().Text())
	}

	bodyBytes, err := io.ReadAll(res.Body().Raw())
	if err != nil {
		return nil, err
	}

	var chatCompletion ChatCompletionResponse
	if err = json.Unmarshal(bodyBytes, &chatCompletion); err != nil {
		return nil, errors.New("could not read the response body")
	}

	if len(chatCompletion.Choices) == 0 {
		return nil, errors.New("GPT did not give an answer")
	}

	return map[string]interface{}{
		"name":           "openai-prompt-chatgpt",
		"usage_mode":     "operation",
		"prompt":         inputPt.Prompt,
		"model":          inputPt.Model,
		"model_settings": requestBody,
		"gpt_answer":     chatCompletion.Choices[0].Message.Content,
	}, nil
}

func (c *PromptChatGPTOperation) Test(ctx *sdk.RunContext) (sdk.JSON, error) {
	return c.Run(ctx)
}

func (c *PromptChatGPTOperation) GetInfo() *sdk.OperationInfo {
	return c.options
}
