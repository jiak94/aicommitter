package openai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	. "github.com/jiak94/aicommit/config"
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
	Error struct {
		Message string `json:"message"`
	}
}

const (
	URL = "https://api.openai.com/v1/chat/completions"
)

var ()

func ChatWithGPT3(diff string) (string, error) {
	config, err := GetConfig()
	if err != nil {
		return "", err
	}

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
	requestData := Request{Messages: messages, Model: config.Model}
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
	req.Header.Set("Authorization", "Bearer "+config.OpenAIKey)

	client := &http.Client{
		Timeout: time.Duration(config.Timeout) * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return "", err
	}

	defer resp.Body.Close()
	var responseData Response
	err = json.NewDecoder(resp.Body).Decode(&responseData)
	if err != nil {
		fmt.Println("Error decoding response:", err)
		return "", err
	}
	if resp.StatusCode != 200 {
		fmt.Println(responseData.Error.Message)
		err := fmt.Errorf(resp.Status)
		return "", err
	}

	return responseData.Choices[0].Message.Content, nil
}
