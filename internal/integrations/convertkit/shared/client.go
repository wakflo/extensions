package shared

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	sdkcore "github.com/wakflo/go-sdk/core"
)

const (
	BaseURL = "https://api.convertkit.com/v3"
)

func GetConvertKitClient(path, method string, body io.Reader) (sdkcore.JSON, error) {
	client := &http.Client{}
	fullURL := BaseURL + path

	req, err := http.NewRequest(method, fullURL, body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode >= 400 {
		errMsg := fmt.Sprintf("ConvertKit API Error: %s (Status: %d)", string(bodyBytes), resp.StatusCode)
		return nil, errors.New(errMsg)
	}

	var response map[string]interface{}
	if err := json.Unmarshal(bodyBytes, &response); err != nil {
		return nil, err
	}

	return response, nil
}

func PostConvertKitClient(path string, payload sdkcore.JSON) (sdkcore.JSON, error) {
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	return GetConvertKitClient(path, http.MethodPost, strings.NewReader(string(payloadBytes)))
}
