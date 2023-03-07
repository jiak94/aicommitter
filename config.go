package main

import (
	"fmt"
	"os"

	"github.com/BurntSushi/toml"
)

type OpenAIConfig struct {
	OpenAIKey string `toml:"OpenAIKey"`
	Model     string `toml:"Model"`
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
		fmt.Println("Please update your OpenAI API key using\naicommit config --api-key <your_api_key>")
		os.Exit(1)
	}
	if _config.OpenAIKey == "" || _config.OpenAIKey == "your_api_key" {
		fmt.Println("Please update your OpenAI API key using\naicommit config --api-key <your_api_key>")
		os.Exit(1)
	}
}

func setConfig(model, api_key string) {
	configLocation := getConfigLocation()
	// Check if config directory exists
	if _, err := os.Stat(configLocation); os.IsNotExist(err) {
		if err := os.MkdirAll(configLocation, 0755); err != nil {
			fmt.Printf("Error creating aicommitter config directory: %v\n", err)
			os.Exit(1)
		}
	}

	config := OpenAIConfig{}
	// Check if config file exists
	if _, err := os.Stat(configLocation + "/config.toml"); os.IsNotExist(err) {
		if model == "" {
			config.Model = DEFAULT_MODEL
		} else {
			config.Model = model
		}
		
		if api_key == "" {
			config.OpenAIKey = "your_api_key"
		} else {
			config.OpenAIKey = api_key
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