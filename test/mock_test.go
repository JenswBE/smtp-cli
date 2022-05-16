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
	_, err := http.DefaultClient.Do(req)
	return err
}

func getMessages() ([]message, error) {
	// Call SMTP mock server
	resp, err := http.Get(smtpMockBaseURL + "/api/Messages")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Parse response
	var messages []message
	err = json.NewDecoder(resp.Body).Decode(&messages)
	return messages, err
}

func getMessageBody(id string) (string, error) {
	// Call SMTP mock server
	url := fmt.Sprintf("%s/api/Messages/%s/raw", smtpMockBaseURL, id)
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Parse response
	var builder strings.Builder
	_, err = io.Copy(&builder, resp.Body)
	return builder.String(), err
}
