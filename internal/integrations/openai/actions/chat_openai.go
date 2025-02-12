package actions

import (
	"encoding/json"
	"errors"
	"io"

	"github.com/wakflo/go-sdk/autoform"
	sdkcore "github.com/wakflo/go-sdk/core"
	"github.com/wakflo/go-sdk/sdk"
)

type chatOpenAIActionProps struct {
	Model            string   `json:"model"`
	Prompt           string   `json:"prompt"`
	MaxTokens        *int     `json:"max_tokens,omitempty"`
	FrequencyPenalty *float64 `json:"frequency_penalty,omitempty"` // -2.0 to 2.0
	PresencePenalty  *float64 `json:"presence_penalty,omitempty"`  // -2.0 to 2.0
	Seed             *int     `json:"seed,omitempty"`              // random integer
	Temperature      *float64 `json:"temperature,omitempty"`       // 0 to 2.0
	TopP             *float64 `json:"top_p,omitempty"`             // not said, but it's an alternative to Temperature, so probably 0 to 2.0
}

type ChatOpenAIAction struct{}

func (a *ChatOpenAIAction) Name() string {
	return "Chat OpenAI"
}

func (a *ChatOpenAIAction) Description() string {
	return "Integrate with OpenAI's chatbot API to automate conversations and generate human-like responses within your workflow. This integration enables you to leverage OpenAI's vast language model capabilities to provide personalized support, answer frequent questions, and even create custom workflows based on user input."
}

func (a *ChatOpenAIAction) GetType() sdkcore.ActionType {
	return sdkcore.ActionTypeNormal
}

func (a *ChatOpenAIAction) Documentation() *sdk.OperationDocumentation {
	return &sdk.OperationDocumentation{
		Documentation: &chatOpenAIDocs,
	}
}

func (a *ChatOpenAIAction) Icon() *string {
	return nil
}

func (a *ChatOpenAIAction) Properties() map[string]*sdkcore.AutoFormSchema {
	getGPTModels := func(ctx *sdkcore.DynamicFieldContext) (*sdkcore.DynamicOptionsResponse, error) {
		models, err := getModels(ctx.Auth.Secret, "gpt")
		if err != nil {
			return nil, err
		}

		return ctx.Respond(models, len(models))
	}

	return map[string]*sdkcore.AutoFormSchema{
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
	}
}

func (a *ChatOpenAIAction) Perform(ctx sdk.PerformContext) (sdkcore.JSON, error) {
	input, err := sdk.InputToTypeSafely[chatOpenAIActionProps](ctx.BaseContext)
	if err != nil {
		return nil, err
	}

	// Validate inputs
	if err := validateInput(input); err != nil {
		return nil, err
	}

	// Build request body
	requestBody := buildRequestBody(input)

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
		"prompt":         input.Prompt,
		"model":          input.Model,
		"model_settings": requestBody,
		"gpt_answer":     chatCompletion.Choices[0].Message.Content,
	}, nil
}

func (a *ChatOpenAIAction) Auth() *sdk.Auth {
	return nil
}

func (a *ChatOpenAIAction) SampleData() sdkcore.JSON {
	return map[string]any{
		"name":       "openai-prompt-chatgpt",
		"usage_mode": "operation",
		"prompt":     "Who won the Oscar for best actor in 2010?",
		"max_tokens": "500",
		"gpt_answer": "The Oscar for Best Actor at the 82nd Academy Awards, held in 2010, was won by Jeff Bridges for his role as Otis \"Bad\" Blake in the film \"Crazy Heart.\"",
	}
}

func (a *ChatOpenAIAction) Settings() sdkcore.ActionSettings {
	return sdkcore.ActionSettings{}
}

func NewChatOpenAIAction() sdk.Action {
	return &ChatOpenAIAction{}
}
