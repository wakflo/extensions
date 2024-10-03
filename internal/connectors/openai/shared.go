package openai

import (
	"encoding/json"
	"errors"
	"io"
	"strings"

	fastshot "github.com/opus-domini/fast-shot"
)

func getOpenAiClient(token string) (fastshot.ClientHttpMethods, error) {
	return fastshot.NewClient("https://api.openai.com/v1").
		Auth().BearerToken(token).
		Header().AddAccept("application/json").
		Build(), nil
}

func getModels(token string, modelPrefix string) (interface{}, error) {
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
		if strings.HasPrefix(model.Id, modelPrefix) {
			models = append(models, model)
		}
	}

	return models, nil
}
