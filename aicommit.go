package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"

	. "github.com/jiak94/aicommit/config"
	. "github.com/jiak94/aicommit/openai"
)

const (
	VERSION = "1.0.2"
)

var (
	ShowVersion = flag.Bool("version", false, "Print the version information")

	// config subcommand
	ConfigCmd    = flag.NewFlagSet("config", flag.ExitOnError)
	ConfigModel  = ConfigCmd.String("model", "", "The model to use")
	ConfigAPIKey = ConfigCmd.String("api-key", "", "The OpenAI API key")
	ConfigTimeout = ConfigCmd.Int("timeout", 5, "The timeout for the OpenAI API call (in seconds)")
	ShowConfig   = ConfigCmd.Bool("show", false, "Show the current configuration")

	// registerHook subcommand
	RegisterHookCmd   = flag.NewFlagSet("registerHook", flag.ExitOnError)
	RegisterHookForce = RegisterHookCmd.Bool("f", false, "Overwrite the existing prepare-commit-msg hook")

	// help flag
	Help = flag.Bool("help", false, "Print the help information")
)

func main() {
	flag.Parse()

	if *Help {
		printHelp()
		return
	}

	if *ShowVersion {
		fmt.Printf("Version: %s\n", VERSION)
		return
	}

	if len(flag.Args()) > 0 {
		switch flag.Args()[0] {
		case "config":
			ConfigCmd.Parse(os.Args[2:])

			if *ShowConfig {
				config, err := GetConfig()
				if err != nil {
					fmt.Printf("Error getting config: %s\n", err.Error())
					os.Exit(1)
				}

				fmt.Printf("Model: %s, API key: %s\n", config.Model, config.OpenAIKey)
				return
			}
			if len(*ConfigModel) == 0 && len(*ConfigAPIKey) == 0 {
				fmt.Println("Please specify a model or an API key")
				os.Exit(1)
			}
			SetConfig(*ConfigModel, *ConfigAPIKey, *ConfigTimeout)
			return
		case "registerHook":
			RegisterHookCmd.Parse(os.Args[2:])
			registerPreCommitMsgHook(*RegisterHookForce)
			return
		default:
			// if the first argument is a file, process it
			if _, err := os.Stat(flag.Args()[0]); err == nil {
				// process file
				processFile(flag.Args()[0])
				return
			}
			fmt.Printf("Unknown command: %s", flag.Args()[0])
			os.Exit(1)
		}
	}

	diff, err := getDiff()
	if err != nil {
		fmt.Printf("Error getting diff: %v\n", err)
		return
	}

	commitMsg, err := ChatWithGPT3(diff)
	if err != nil {
		fmt.Printf("Error getting response from GPT-3: %s\n", err.Error())
		return
	}
	fmt.Println(commitMsg)
}

func processFile(file string) {
	diff, err := getDiff()
	if err != nil {
		return
	}
	commitMsg, err := ChatWithGPT3(diff)

	if err != nil {
		return
	}

	WriteCommitMessage(commitMsg, file)
}

func WriteCommitMessage(message, file string) error {
	// Read the contents of the file
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}

	// Combine the string to add and the file contents into a single byte slice
	newContents := append([]byte(message), data...)

	// Write the new contents back to the file
	err = ioutil.WriteFile(file, newContents, 0644)
	if err != nil {
		return err
	}
	return nil
}

func getDiff() (string, error) {
	out, err := exec.Command("git", "diff", "--staged").Output()
	if err != nil {
		return "", err
	}
	if len(out) == 0 {
		return "", fmt.Errorf("No changes to commit")
	}
	return string(out), nil
}

func printHelp() {
	fmt.Println("Usage: aicommit [options] [command]")
	fmt.Println("Options:")
	fmt.Println("  -version\tPrint the version information")
	fmt.Println("  -help\t\tPrint the help information")
	fmt.Println("Commands:")
	fmt.Println("  config --api-key <your_api_key>\tSet the OpenAI API key")
	fmt.Println("  config --model <model_name>\t\tSet the model to use")
	fmt.Println("  config --timeout <timeout>\t\tSet the timeout for the OpenAI API call (in seconds)")
	fmt.Println("  config --show\t\tshow the current configuration")
	fmt.Println("  registerHook\t\tregister the prepare-commit-msg hook")
}
