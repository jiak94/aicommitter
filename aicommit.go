package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
)

const (
	URL           = "https://api.openai.com/v1/chat/completions"
	DEFAULT_MODEL = "gpt-3.5-turbo"
	VERSION       = "1.0.1"
)

var _config OpenAIConfig

func main() {
	showVersion := flag.Bool("version", false, "Print the version information")

	// config subcommand
	config := flag.NewFlagSet("config", flag.ExitOnError)
	model := config.String("model", "", "The model to use")
	api_key := config.String("api-key", "", "The OpenAI API key")
	show_config := config.Bool("show", false, "Show the current configuration")

	// help flag
	help := flag.Bool("help", false, "Print the help information")

	flag.Parse()

	if *help {
		printHelp()
		return
	}

	if *showVersion {
		fmt.Printf("Version: %s\n", VERSION)
		return
	}

	if len(flag.Args()) > 0 {
		switch flag.Args()[0] {
		case "config":
			config.Parse(os.Args[2:])

			if *show_config {
				getConfig()
				fmt.Printf("Model: %s, API key: %s\n", _config.Model, _config.OpenAIKey)
				return
			}
			if len(*model) == 0 && len(*api_key) == 0 {
				fmt.Println("Please specify a model or an API key")
				os.Exit(1)
			}
			setConfig(*model, *api_key)
			return
		default:
			fmt.Printf("Unknown command: %s", flag.Args()[0])
			os.Exit(1)
		}
	}

	getConfig()

	diff, err := getDiff()
	if err != nil {
		fmt.Printf("Error getting diff: %v\n", err)
		os.Exit(1)
	}

	commitMsg, err := chatWithGPT3(diff)
	if err != nil {
		fmt.Printf("Error getting response from GPT-3: %s\n", err.Error())
		os.Exit(1)
	}

	writeCommitMessage(commitMsg)
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
}
