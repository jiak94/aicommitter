package main

import "flag"

var (
	ShowVersion = flag.Bool("version", false, "Print the version information")

	// config subcommand
	ConfigCmd = flag.NewFlagSet("config", flag.ExitOnError)
	ConfigModel = ConfigCmd.String("model", "", "The model to use")
	ConfigAPIKey = ConfigCmd.String("api-key", "", "The OpenAI API key")
	ShowConfig = ConfigCmd.Bool("show", false, "Show the current configuration")

	// registerHook subcommand
	RegisterHookCmd = flag.NewFlagSet("registerHook", flag.ExitOnError)
	RegisterHookForce = RegisterHookCmd.Bool("f", false, "Overwrite the existing pre-commit hook")

	// help flag
	Help = flag.Bool("help", false, "Print the help information")
)
