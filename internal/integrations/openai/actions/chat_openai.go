package actions

import (
	"encoding/json"
	"errors"
	"io"

	"github.com/juicycleff/smartform/v1"
	"github.com/wakflo/go-sdk/v2"
	sdkcontext "github.com/wakflo/go-sdk/v2/context"
	"github.com/wakflo/go-sdk/v2/core"
	sdkcore "github.com/wakflo/go-sdk/v2/core"
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

// Metadata returns metadata about the action
func (a *ChatOpenAIAction) Metadata() sdk.ActionMetadata {
	return sdk.ActionMetadata{
		ID:            "chat_openai",
		DisplayName:   "Chat OpenAI",
		Description:   "Integrate with OpenAI's chatbot API to automate conversations and generate human-like responses within your workflow. This integration enables you to leverage OpenAI's vast language model capabilities to provide personalized support, answer frequent questions, and even create custom workflows based on user input.",
		Type:          core.ActionTypeAction,
		Documentation: chatOpenAIDocs,
		SampleOutput: map[string]any{
			"name":       "openai-prompt-chatgpt",
			"usage_mode": "operation",
			"prompt":     "Who won the Oscar for best actor in 2010?",
			"max_tokens": "500",
			"gpt_answer": "The Oscar for Best Actor at the 82nd Academy Awards, held in 2010, was won by Jeff Bridges for his role as Otis \"Bad\" Blake in the film \"Crazy Heart.\"",
		},
		Settings: core.ActionSettings{},
	}
}

// Properties returns the schema for the action's input configuration
func (a *ChatOpenAIAction) Properties() *smartform.FormSchema {
	form := smartform.NewForm("chat_openai", "Chat OpenAI")

	getGPTModels := func(ctx sdkcontext.DynamicFieldContext) (*sdkcore.DynamicOptionsResponse, error) {
		authCtx, err := ctx.AuthContext()
		if err != nil {
			return nil, err
		}

		models, err := getModels(authCtx.Extra["token"], "gpt")
		if err != nil {
			return nil, err
		}

		return ctx.Respond(models, len(models))
	}

	form.SelectField("model", "Model").
		Required(true).
		HelpText("Choose model which you will interact").
		WithDynamicOptions(
			smartform.NewOptionsBuilder().
				Dynamic().
				WithFunctionOptions(sdk.WithDynamicFunctionCalling(&getGPTModels)).
				WithSearchSupport().
				WithPagination(10).
				End().
				GetDynamicSource(),
		)

	form.TextField("prompt", "Prompt").
		Required(true).
		HelpText("What would you like to ask ChatGPT?")

	form.NumberField("max_tokens", "Max Tokens").
		Required(false).
		HelpText("Max tokens that ChatGPT can use to generate it's answer.")

	form.NumberField("frequency_penalty", "Frequency penalty").
		Required(false).
		HelpText("Penalize or not new tokens (letters/words) based on their frequency. Positive values penalize same verbatim while negative values don't (-2.0 to 2.0)")

	form.NumberField("presence_penalty", "Presence penalty").
		Required(false).
		HelpText("Penalize or not new tokens (letters/words) based on their appearance. Positive values will increase the chance of new topics while negative values don't (-2.0 to 2.0)")

	form.NumberField("seed", "Random seed").
		Required(false).
		HelpText("By using a random seed the model will try its best to replicate the answer that it gave before with the same seed")

	form.NumberField("temperature", "Temperature").
		Required(false).
		HelpText("Control model randomness with higher values and more focused and deterministic with lowest values (0 to 2)")

	form.NumberField("top_p", "Top P").
		Required(false).
		HelpText("An alternative to temperature that considers the probability of the tokens appearing. Lower values consider low probabilities (0 to 2). It's advised to only use this or temperature and not both.")

	schema := form.Build()

	return schema
}

// Auth returns the authentication requirements for the action
func (a *ChatOpenAIAction) Auth() *core.AuthMetadata {
	return nil
}

// Perform executes the action with the given context and input
func (a *ChatOpenAIAction) Perform(ctx sdkcontext.PerformContext) (core.JSON, error) {
	input, err := sdk.InputToTypeSafely[chatOpenAIActionProps](ctx)
	if err != nil {
		return nil, err
	}

	authCtx, err := ctx.AuthContext()
	if err != nil {
		return nil, err
	}

	// Validate inputs
	if err := validateInput(input); err != nil {
		return nil, err
	}

	// Build request body
	requestBody := buildRequestBody(input)

	client, err := getOpenAiClient(authCtx.Secret)
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

func NewChatOpenAIAction() sdk.Action {
	return &ChatOpenAIAction{}
}
