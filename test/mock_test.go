package test

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

const smtpMockBaseURL = "http://localhost:8080"

type message struct {
	ID      string `json:"id"`
	From    string `json:"from"`
	To      string `json:"to"`
	Subject string `json:"subject"`
}

func clearMessages() error {
	// /api/Messages/*
	req, _ := http.NewRequest(http.MethodDelete, smtpMockBaseURL+"/api/Messages/*", nil)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		_ = res.Body.Close()
	}
	return fmt.Errorf("failed to clear messages: %w", err)
}

func getMessages() ([]message, error) {
	// Call SMTP mock server
	resp, err := http.Get(smtpMockBaseURL + "/api/Messages")
	if err != nil {
		return nil, fmt.Errorf("failed to get messages: %w", err)
	}
	defer resp.Body.Close()

	// Parse response
	var messages []message
	err = json.NewDecoder(resp.Body).Decode(&messages)
	return messages, fmt.Errorf("failed to parse get messages response: %w", err)
}

func getMessageBody(id string) (string, error) {
	// Call SMTP mock server
	url := fmt.Sprintf("%s/api/Messages/%s/raw", smtpMockBaseURL, id)
	resp, err := http.Get(url) // #nosec G107
	if err != nil {
		return "", fmt.Errorf("failed to get message body: %w", err)
	}
	defer resp.Body.Close()

	// Parse response
	var builder strings.Builder
	_, err = io.Copy(&builder, resp.Body)
	return builder.String(), fmt.Errorf("failed to parse get message body response: %w", err)
}
