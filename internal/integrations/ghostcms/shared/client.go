package shared

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type GhostClient struct {
	SiteURL string
	APIKey  string
	client  *http.Client
}

func NewGhostClient(authData map[string]string) (*GhostClient, error) {
	siteURL, ok := authData["site_url"]
	if !ok || siteURL == "" {
		return nil, fmt.Errorf("site URL is required")
	}

	apiKey, ok := authData["admin_api_key"]
	if !ok || apiKey == "" {
		return nil, fmt.Errorf("admin API key is required")
	}

	// Remove trailing slash from site URL
	siteURL = strings.TrimSuffix(siteURL, "/")

	return &GhostClient{
		SiteURL: siteURL,
		APIKey:  apiKey,
		client:  &http.Client{Timeout: 30 * time.Second},
	}, nil
}

func (c *GhostClient) generateToken() (string, error) {
	// Split the API key into ID and secret
	parts := strings.Split(c.APIKey, ":")
	if len(parts) != 2 {
		return "", fmt.Errorf("invalid API key format")
	}

	keyID := parts[0]
	secret := parts[1]

	// Decode the secret from hex
	decodedSecret, err := hex.DecodeString(secret)
	if err != nil {
		return "", fmt.Errorf("failed to decode secret: %w", err)
	}

	// Create the token
	now := time.Now()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"aud": "/admin/",
		"exp": now.Add(5 * time.Minute).Unix(),
		"iat": now.Unix(),
	})

	// Add the key ID to the header
	token.Header["kid"] = keyID

	// Sign the token
	tokenString, err := token.SignedString(decodedSecret)
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return tokenString, nil
}

func (c *GhostClient) makeRequest(method, endpoint string, body interface{}) (map[string]interface{}, error) {
	url := fmt.Sprintf("%s/ghost/api/admin%s", c.SiteURL, endpoint)

	var reqBody io.Reader
	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
		reqBody = bytes.NewBuffer(jsonBody)
	}

	req, err := http.NewRequest(method, url, reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Generate JWT token for authentication
	token, err := c.generateToken()
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Ghost %s", token))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept-Version", "v5.0")

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode >= 400 {
		var errorResp map[string]interface{}
		if err := json.Unmarshal(respBody, &errorResp); err == nil {
			if errors, ok := errorResp["errors"].([]interface{}); ok && len(errors) > 0 {
				if firstError, ok := errors[0].(map[string]interface{}); ok {
					if message, ok := firstError["message"].(string); ok {
						return nil, fmt.Errorf("ghost API error: %s", message)
					}
				}
			}
		}
		return nil, fmt.Errorf("request failed with status %d: %s", resp.StatusCode, string(respBody))
	}

	if len(respBody) == 0 {
		return map[string]interface{}{}, nil
	}

	var result map[string]interface{}
	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return result, nil
}

func (c *GhostClient) Get(endpoint string) (map[string]interface{}, error) {
	return c.makeRequest("GET", endpoint, nil)
}

func (c *GhostClient) Post(endpoint string, body interface{}) (map[string]interface{}, error) {
	return c.makeRequest("POST", endpoint, body)
}

func (c *GhostClient) Put(endpoint string, body interface{}) (map[string]interface{}, error) {
	return c.makeRequest("PUT", endpoint, body)
}

func (c *GhostClient) Delete(endpoint string) (map[string]interface{}, error) {
	return c.makeRequest("DELETE", endpoint, nil)
}

func (c *GhostClient) UploadImage(imageData []byte, filename string) (map[string]interface{}, error) {
	url := fmt.Sprintf("%s/ghost/api/admin/images/upload", c.SiteURL)

	// Create multipart form data
	body := &bytes.Buffer{}
	boundary := fmt.Sprintf("----WebKitFormBoundary%d", time.Now().Unix())

	// Write file field
	body.WriteString(fmt.Sprintf("--%s\r\n", boundary))
	body.WriteString(fmt.Sprintf("Content-Disposition: form-data; name=\"file\"; filename=\"%s\"\r\n", filename))
	body.WriteString("Content-Type: application/octet-stream\r\n\r\n")
	body.Write(imageData)
	body.WriteString("\r\n")

	// Write purpose field
	body.WriteString(fmt.Sprintf("--%s\r\n", boundary))
	body.WriteString("Content-Disposition: form-data; name=\"purpose\"\r\n\r\n")
	body.WriteString("image")
	body.WriteString("\r\n")

	// Close boundary
	body.WriteString(fmt.Sprintf("--%s--\r\n", boundary))

	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Generate JWT token for authentication
	token, err := c.generateToken()
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Ghost %s", token))
	req.Header.Set("Content-Type", fmt.Sprintf("multipart/form-data; boundary=%s", boundary))
	req.Header.Set("Accept-Version", "v5.0")

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("upload failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("upload failed with status %d: %s", resp.StatusCode, string(respBody))
	}

	var result map[string]interface{}
	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return result, nil
}

// Helper function to validate webhook signature
func ValidateWebhookSignature(secret, signature string, body []byte) bool {
	h := hmac.New(sha256.New, []byte(secret))
	h.Write(body)
	expectedSig := hex.EncodeToString(h.Sum(nil))
	return hmac.Equal([]byte(expectedSig), []byte(signature))
}
