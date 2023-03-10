package config

import (
	"io/ioutil"
	"os"
	"reflect"
	"testing"

	"github.com/BurntSushi/toml"
)

func TestGetConfigLocation(t *testing.T) {
	// Set XDG_CONFIG_HOME environment variable to a temporary directory
	tempDir, err := ioutil.TempDir("", "config_test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tempDir)
	os.Setenv("XDG_CONFIG_HOME", tempDir)

	// Call the getConfigLocation function and verify the result
	expectedLocation := tempDir + "/aicommitter"
	actualLocation := getConfigLocation()
	if actualLocation != expectedLocation {
		t.Errorf("getConfigLocation() = %q, want %q", actualLocation, expectedLocation)
	}

	// Unset XDG_CONFIG_HOME environment variable and call the function again
	os.Unsetenv("XDG_CONFIG_HOME")
	expectedLocation = os.Getenv("HOME") + "/.config/aicommitter"
	actualLocation = getConfigLocation()
	if actualLocation != expectedLocation {
		t.Errorf("getConfigLocation() = %q, want %q", actualLocation, expectedLocation)
	}
}

func TestGetConfig(t *testing.T) {
	// Create a temporary config file with known contents
	tempDir, err := ioutil.TempDir("", "config_test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tempDir)
	os.Mkdir(tempDir+"/aicommitter", 0755)
	configFile := tempDir + "/aicommitter/config.toml"
	err = ioutil.WriteFile(configFile, []byte(`OpenAIKey = "test_key"`), 0644)
	if err != nil {
		t.Fatal(err)
	}

	// Call the GetConfig function with the temporary config file
	os.Setenv("XDG_CONFIG_HOME", tempDir)
	config, err := GetConfig()

	// Verify that the function returns the expected OpenAIConfig object and no error
	expectedConfig := &OpenAIConfig{OpenAIKey: "test_key"}
	if !reflect.DeepEqual(config, expectedConfig) {
		t.Errorf("GetConfig() = %v, want %v", config, expectedConfig)
	}
	if err != nil {
		t.Errorf("GetConfig() returned an unexpected error: %v", err)
	}

	// Modify the config file to have an empty OpenAIKey
	err = ioutil.WriteFile(configFile, []byte(`OpenAIKey = ""`), 0644)
	if err != nil {
		t.Fatal(err)
	}

	// Call the function again and verify that it returns an error
	config, err = GetConfig()
	if config != nil {
		t.Errorf("GetConfig() returned %v, want nil", config)
	}
	expectedError := "Please update your OpenAI API key using\naicommit config --api-key <your_api_key>"
	if err == nil || err.Error() != expectedError {
		t.Errorf("GetConfig() returned error %q, want %q", err, expectedError)
	}
}

func TestSetConfig(t *testing.T) {
	// Create a temporary directory for the test
	tmpDir, err := ioutil.TempDir("", "aicommitter")
	if err != nil {
		t.Fatalf("Failed to create temporary directory: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Set the environment variable for XDG_CONFIG_HOME to the temporary directory
	os.Setenv("XDG_CONFIG_HOME", tmpDir)

	// Test with default values
	SetConfig("", "", 0)
	expectedConfig := OpenAIConfig{
		Timeout:   DEFAULT_TIMEOUT,
		OpenAIKey: "your_api_key",
		Model:     DEFAULT_MODEL,
	}
	assertConfigEquals(t, tmpDir+"/aicommitter/config.toml", &expectedConfig)

	// Test with custom values
	SetConfig("davinci-codex", "abc123", 30)
	expectedConfig = OpenAIConfig{
		Timeout:   30,
		OpenAIKey: "abc123",
		Model:     "davinci-codex",
	}
	assertConfigEquals(t, tmpDir+"/aicommitter/config.toml", &expectedConfig)

	// Test updating an existing config file
	err = ioutil.WriteFile(tmpDir+"/aicommitter/config.toml", []byte(`Timeout = 10
OpenAIKey = "old_key"
Model = "davinci"`), 0644)
	if err != nil {
		t.Fatalf("Failed to create config file: %v", err)
	}
	SetConfig("davinci-codex", "new_key", 0)
	expectedConfig = OpenAIConfig{
		Timeout:   10,
		OpenAIKey: "new_key",
		Model:     "davinci-codex",
	}
	assertConfigEquals(t, tmpDir+"/aicommitter/config.toml", &expectedConfig)
}

// Helper function to check if the config file matches the expected config
func assertConfigEquals(t *testing.T, filename string, expected *OpenAIConfig) {
	actual := &OpenAIConfig{}
	_, err := toml.DecodeFile(filename, actual)
	if err != nil {
		t.Fatalf("Failed to decode config file: %v", err)
	}
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("Config mismatch:\nExpected: %+v\nActual:   %+v", expected, actual)
	}
}