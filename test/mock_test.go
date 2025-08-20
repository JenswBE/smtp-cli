package test

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/rs/zerolog/log"
)

const (
	smtpMockImplictTLSBaseURL = "http://localhost:8080"
	smtpMockSTARTTLSBaseURL   = "http://localhost:8081"
)

type getMessagesResponse struct {
	Messages []message `json:"results"`
}

type message struct {
	ID      string   `json:"id"`
	From    string   `json:"from"`
	To      []string `json:"to"`
	Subject string   `json:"subject"`
}

func pingServer(baseURL string) error {
	resp, err := http.Get(baseURL + "/api/Version")
	if err != nil {
		return fmt.Errorf("failed to ping server: %w", err)
	}
	_ = resp.Body.Close()
	return nil
}

func clearMessages(baseURL string) error {
	// /api/Messages/*
	req, _ := http.NewRequest(http.MethodDelete, baseURL+"/api/Messages/*", nil)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to clear messages: %w", err)
	}
	_ = res.Body.Close()
	return nil
}

func getMessages(baseURL string) ([]message, error) {
	// Call SMTP mock server
	resp, err := http.Get(baseURL + "/api/Messages")
	if err != nil {
		return nil, fmt.Errorf("failed to get messages: %w", err)
	}
	defer silentClose(resp.Body)

	// Read response
	rawResponse, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Parse response
	response := getMessagesResponse{}
	if err = json.Unmarshal(rawResponse, &response); err != nil {
		return nil, fmt.Errorf("failed to parse get messages response: %w. Body: %v", err, string(rawResponse))
	}
	return response.Messages, nil
}

func getMessageBody(baseURL string, id string) (string, error) {
	// Call SMTP mock server
	url := fmt.Sprintf("%s/api/Messages/%s/raw", baseURL, id)
	resp, err := http.Get(url) // #nosec G107
	if err != nil {
		return "", fmt.Errorf("failed to get message body: %w", err)
	}
	defer silentClose(resp.Body)

	// Parse response
	var builder strings.Builder
	if _, err = io.Copy(&builder, resp.Body); err != nil {
		return "", fmt.Errorf("failed to parse get message body response: %w", err)
	}
	return builder.String(), nil
}

func silentClose(closer io.Closer) {
	err := closer.Close()
	if err != nil {
		log.Printf("Failed to close body: %v", err)
	}
}
