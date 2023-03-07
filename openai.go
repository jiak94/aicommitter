package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type Request struct {
	Messages []Message `json:"messages"`
	Model    string    `json:"model"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type Response struct {
	Choices []struct {
		Message Message `json:"message"`
	} `json:"choices"`
}

func chatWithGPT3(diff string) (string, error) {
	systemMessage := "You are a helpful assistant writes short git commit messages."
	userMessage := fmt.Sprintf("%s\n\nWrite the commit message.", diff)

	messages := []Message{
		{
			Role:    "system",
			Content: systemMessage,
		},
		{
			Role:    "user",
			Content: userMessage,
		},
	}
	requestData := Request{Messages: messages, Model: _config.Model}
	requestBody, err := json.Marshal(requestData)

	if err != nil {
		fmt.Println("Error marshalling request:", err)
		return "", err
	}

	req, err := http.NewRequest("POST", URL, bytes.NewBuffer(requestBody))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+_config.OpenAIKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return "", err
	}
	if resp.StatusCode != 200 {
		err := fmt.Errorf(resp.Status)
		return "", err
	}
	defer resp.Body.Close()

	var responseData Response
	err = json.NewDecoder(resp.Body).Decode(&responseData)
	if err != nil {
		fmt.Println("Error decoding response:", err)
		return "", err
	}
	return responseData.Choices[0].Message.Content, nil
}
