package config

import (
	"fmt"
	"os"

	"github.com/BurntSushi/toml"
)

const (
	DEFAULT_MODEL   = "gpt-3.5-turbo"
	DEFAULT_TIMEOUT = 5
)

type OpenAIConfig struct {
	OpenAIKey string `toml:"OpenAIKey"`
	Model     string `toml:"Model"`
	Timeout   int    `toml:"Timeout"`
}

func getConfigLocation() string {
	configDir := os.Getenv("XDG_CONFIG_HOME")
	if configDir == "" {
		configDir = os.Getenv("HOME") + "/.config"
	}
	return configDir + "/aicommitter"
}

func GetConfig() (*OpenAIConfig, error) {
	config := OpenAIConfig{}
	if _, err := toml.DecodeFile(getConfigLocation()+"/config.toml", &config); err != nil {
		fmt.Println(err)
		return nil, fmt.Errorf("Please update your OpenAI API key using\naicommit config --api-key <your_api_key>")
	}
	if config.OpenAIKey == "" || config.OpenAIKey == "your_api_key" {
		return nil, fmt.Errorf("Please update your OpenAI API key using\naicommit config --api-key <your_api_key>")
	}
	return &config, nil
}

func SetConfig(model, api_key string, timeout int) {
	configLocation := getConfigLocation()
	// Check if config directory exists
	if _, err := os.Stat(configLocation); os.IsNotExist(err) {
		if err := os.MkdirAll(configLocation, 0755); err != nil {
			fmt.Printf("Error creating aicommitter config directory: %v\n", err)
			os.Exit(1)
		}
	}

	config := OpenAIConfig{
		Timeout:   DEFAULT_TIMEOUT,
		OpenAIKey: "your_api_key",
		Model:     DEFAULT_MODEL,
	}
	// Check if config file exists
	if _, err := os.Stat(configLocation + "/config.toml"); os.IsNotExist(err) {
		if model != "" {
			config.Model = model
		}

		if api_key != "" {
			config.OpenAIKey = api_key
		}

		if timeout > 0 {
			config.Timeout = timeout
		}

	} else {
		if _, err := toml.DecodeFile(configLocation+"/config.toml", &config); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		if model != "" {
			config.Model = model
		}

		if api_key != "" {
			config.OpenAIKey = api_key
		}

		if timeout > 0 {
			config.Timeout = timeout
		}
	}
	file, err := os.Create(configLocation + "/config.toml")
	defer file.Close()

	if err != nil {
		fmt.Printf("Error creating aicommitter config file: %v\n", err)
		os.Exit(1)
	}
	if err := toml.NewEncoder(file).Encode(config); err != nil {
		fmt.Printf("Error writing aicommitter config file: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("aicommitter config file updated")
}
