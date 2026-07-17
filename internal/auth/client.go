package auth

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type Client struct {
	meURL      string
	loginURL   string
	httpClient *http.Client
}

func NewClient(meURL string, loginURL string) *Client {
	return &Client{
		meURL:    meURL,
		loginURL: loginURL,
		httpClient: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

func (c *Client) Login(loginRequest LoginRequest) (*LoginResponse, int, error) {
	if strings.TrimSpace(c.loginURL) == "" {
		return nil, http.StatusInternalServerError, fmt.Errorf("AUTH_LOGIN_URL is empty")
	}

	email := strings.ToLower(strings.TrimSpace(loginRequest.Email))
	password := loginRequest.Password

	if email == "" || strings.TrimSpace(password) == "" {
		return nil, http.StatusBadRequest, fmt.Errorf("email and password are required")
	}

	payload, err := json.Marshal(map[string]string{
		"email":    email,
		"password": password,
		"role":     "admin",
	})
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	req, err := http.NewRequest(http.MethodPost, c.loginURL, bytes.NewReader(payload))
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, http.StatusBadGateway, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, http.StatusBadGateway, err
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, resp.StatusCode, errors.New(readErrorMessage(body, "login failed"))
	}

	result, err := decodeLoginResponse(body)
	if err != nil {
		return nil, http.StatusBadGateway, err
	}

	return result, http.StatusOK, nil
}

func (c *Client) GetMe(sourceRequest *http.Request) (*MeResponse, int, error) {
	if strings.TrimSpace(c.meURL) == "" {
		return nil, http.StatusInternalServerError, fmt.Errorf("AUTH_ME_URL is empty")
	}

	req, err := http.NewRequest(http.MethodGet, c.meURL, nil)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	if authorization := sourceRequest.Header.Get("Authorization"); authorization != "" {
		req.Header.Set("Authorization", authorization)
	}

	if cookie := sourceRequest.Header.Get("Cookie"); cookie != "" {
		req.Header.Set("Cookie", cookie)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, http.StatusBadGateway, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized {
		return nil, http.StatusUnauthorized, nil
	}

	if resp.StatusCode == http.StatusForbidden {
		return nil, http.StatusForbidden, nil
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, http.StatusBadGateway, fmt.Errorf("auth /me returned status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, http.StatusBadGateway, err
	}

	result, err := decodeMeResponse(body)
	if err != nil {
		return nil, http.StatusBadGateway, err
	}

	return result, http.StatusOK, nil
}

func decodeLoginResponse(body []byte) (*LoginResponse, error) {
	var wrapped struct {
		Data LoginResponse `json:"data"`
	}
	if err := json.Unmarshal(body, &wrapped); err == nil && wrapped.Data.AccessToken != "" {
		return &wrapped.Data, nil
	}

	var direct LoginResponse
	if err := json.Unmarshal(body, &direct); err != nil {
		return nil, err
	}

	if direct.AccessToken == "" {
		return nil, fmt.Errorf("login response does not include accessToken")
	}

	return &direct, nil
}

func decodeMeResponse(body []byte) (*MeResponse, error) {
	var wrapped struct {
		Data MeResponse `json:"data"`
	}
	if err := json.Unmarshal(body, &wrapped); err == nil && wrapped.Data.Role != "" {
		return &wrapped.Data, nil
	}

	var direct MeResponse
	if err := json.Unmarshal(body, &direct); err != nil {
		return nil, err
	}

	if direct.Role == "" {
		return nil, fmt.Errorf("auth /me response does not include role")
	}

	return &direct, nil
}

func readErrorMessage(body []byte, fallback string) string {
	var result struct {
		Error   string `json:"error"`
		Message any    `json:"message"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return fallback
	}

	switch message := result.Message.(type) {
	case string:
		if strings.TrimSpace(message) != "" {
			return message
		}
	case []any:
		parts := make([]string, 0, len(message))
		for _, item := range message {
			if value, ok := item.(string); ok && strings.TrimSpace(value) != "" {
				parts = append(parts, value)
			}
		}
		if len(parts) > 0 {
			return strings.Join(parts, "; ")
		}
	}

	if result.Error != "" {
		return result.Error
	}

	return fallback
}
