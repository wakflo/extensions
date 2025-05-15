package actions

import (
	"encoding/json"
	"errors"
	"io"
	"strings"

	"github.com/juicycleff/smartform/v1"
	fastshot "github.com/opus-domini/fast-shot"
)

var (
	form = smartform.NewAuthForm("openai-auth", "OpenAI API Authentication", smartform.AuthStrategyCustom)
	_    = form.TextField("token", "OpenAI Token (Required*)").
		Required(true).
		HelpText("The api token used to authenticate OpenAI")

	OpenAISharedAuth = form.Build()
)

func getOpenAiClient(token string) (fastshot.ClientHttpMethods, error) {
	return fastshot.NewClient("https://api.openai.com/v1").
		Auth().BearerToken(token).
		Header().AddAccept("application/json").
		Build(), nil
}

func getModels(token string, modelPrefix string) ([]ModelResponse, error) {
	client, err := getOpenAiClient(token)
	if err != nil {
		return nil, err
	}

	res, err := client.GET("/models").Send()
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

	var modelsRes ModelListResponse
	if err = json.Unmarshal(bodyBytes, &modelsRes); err != nil {
		return nil, err
	}

	var models []ModelResponse
	for _, model := range modelsRes.Data {
		if strings.HasPrefix(model.ID, modelPrefix) {
			models = append(models, model)
		}
	}

	return models, nil
}

func validateInput(inputPt *chatOpenAIActionProps) error {
	if len(inputPt.Model) == 0 || !strings.HasPrefix(inputPt.Model, "gpt") {
		return errors.New("invalid open ai model")
	}

	if len(inputPt.Prompt) == 0 {
		return errors.New("avoiding requests, you are sending an empty prompt")
	}

	if inputPt.MaxTokens != nil && *inputPt.MaxTokens <= 0 {
		return errors.New("model max tokens must be greater than zero")
	}

	if inputPt.FrequencyPenalty != nil && (*inputPt.FrequencyPenalty < -2.0 || *inputPt.FrequencyPenalty > 2.0) {
		return errors.New("model frequency penalty value must be between -2.0 and 2.0")
	}

	if inputPt.PresencePenalty != nil && (*inputPt.PresencePenalty < -2.0 || *inputPt.PresencePenalty > 2.0) {
		return errors.New("model presence penalty value must be between -2.0 and 2.0")
	}

	if inputPt.Temperature != nil && (*inputPt.Temperature < 0.0 || *inputPt.Temperature > 2.0) {
		return errors.New("model temperature value must be between 0.0 and 2.0")
	}

	if inputPt.TopP != nil && (*inputPt.TopP < 0.0 || *inputPt.TopP > 2.0) {
		return errors.New("model temperature value must be between 0.0 and 2.0")
	}

	return nil
}

func buildRequestBody(inputPt *chatOpenAIActionProps) map[string]interface{} {
	requestBody := map[string]interface{}{
		"model": inputPt.Model,
		"messages": []interface{}{
			map[string]interface{}{
				"role":    "user",
				"content": inputPt.Prompt,
			},
		},
	}

	if inputPt.MaxTokens != nil {
		requestBody["max_tokens"] = *inputPt.MaxTokens
	}

	if inputPt.Seed != nil {
		requestBody["seed"] = *inputPt.Seed
	}

	if inputPt.FrequencyPenalty != nil {
		requestBody["frequency_penalty"] = *inputPt.FrequencyPenalty
	}

	if inputPt.PresencePenalty != nil {
		requestBody["presence_penalty"] = *inputPt.PresencePenalty
	}

	if inputPt.Temperature != nil {
		requestBody["temperature"] = *inputPt.Temperature
	}

	if inputPt.TopP != nil {
		requestBody["top_p"] = *inputPt.TopP
	}

	return requestBody
}
