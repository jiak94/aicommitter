package main

import (
	"fmt"

	"os"
	"os/exec"

	"bytes"
	"encoding/json"

	"net/http"

	"github.com/BurntSushi/toml"
)

const (
	URL           = "https://api.openai.com/v1/chat/completions"
	DEFAULT_MODEL = "gpt-3.5-turbo"
)

type OpenAIConfig struct {
	OpenAIKey string
	Model     string
}

var _config OpenAIConfig

func main() {
	configLocation := getConfigLocation()
	createOrGetOpenAIConfig(configLocation)

	diff, err := getDiff()
	if err != nil {
		fmt.Printf("Error getting diff: %v\n", err)
		os.Exit(1)
	}

	response, err := chatWithGPT3(diff)
	if err != nil {
		fmt.Printf("Error getting response from GPT-3: %v\n", err)
		os.Exit(1)
	}

	fmt.Println(response)

	writeCommitMessage(response)
	showCommitMessage()
}

func writeCommitMessage(message string) error {
	_, err := exec.Command("git", "commit", "-m", message).Output()
	if err != nil {
		return err
	}
	return nil
}

func showCommitMessage() {
	cmd := exec.Command("git", "commit", "--amend")
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()
}

func getDiff() (string, error) {
	out, err := exec.Command("git", "diff", "HEAD").Output()
	if err != nil {
		return "", err
	}
	if len(out) == 0 {
		return "", fmt.Errorf("No changes to commit")
	}
	return string(out), nil
}

func getConfigLocation() string {
	configDir := os.Getenv("XDG_CONFIG_HOME")
	if configDir == "" {
		configDir = os.Getenv("HOME") + "/.config"
	}
	return configDir + "/aicommitter"
}

func getConfig() {
	if _, err := toml.DecodeFile(getConfigLocation()+"/config.toml", &_config); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	if _config.OpenAIKey == "" {
		fmt.Printf("Please add your OpenAI API key to the config file at %s\n", getConfigLocation())
		os.Exit(1)
	}
}

func createOrGetOpenAIConfig(location string) {
	if _, err := os.Stat(location); os.IsNotExist(err) {
		if err := os.MkdirAll(location, 0755); err != nil {
			fmt.Printf("Error creating aicommitter config directory: %v\n", err)
			os.Exit(1)
		}

		file, err := os.Create(location + "/config.toml")
		if err != nil {
			fmt.Printf("Error creating aicommitter config file: %v\n", err)
			os.Exit(1)
		}
		defer file.Close()

		config := OpenAIConfig{
			OpenAIKey: "your_api_key",
			Model:     DEFAULT_MODEL,
		}
		if err := toml.NewEncoder(file).Encode(config); err != nil {
			fmt.Printf("Error writing aicommitter config file: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("aicommitter config file created at %s, please update your api key\n", location)
		os.Exit(1)

	} else {
		getConfig()
	}
}

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
	defer resp.Body.Close()

	var responseData Response
	err = json.NewDecoder(resp.Body).Decode(&responseData)
	if err != nil {
		fmt.Println("Error decoding response:", err)
		return "", err
	}
	return responseData.Choices[0].Message.Content, nil
}
