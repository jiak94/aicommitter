package main

import (
	"fmt"
	"io/ioutil"
	"os"
)

func registerPreCommitMsgHook(force bool) {
	fmt.Println("Registering hook script...")
	// check if we are in a git repo
	if _, err := os.Stat(".git"); os.IsNotExist(err) {
		fmt.Println(".git directory not found. Are you in a git repository?")
		os.Exit(1)
	}
	
	// check if the hook folder exists
	if _, err := os.Stat(".git/hooks"); os.IsNotExist(err) {
		fmt.Println("Hooks folder not found. Creating it...")
		err := os.Mkdir(".git/hooks", 0755)
		if err != nil {
			fmt.Println("Error creating hooks folder")
			os.Exit(1)
		}
	}
	if _, err := os.Stat(".git/hooks/prepare-commit-msg"); err == nil && !force {
		fmt.Println("Register hook script failed: hook script already exists.")
		os.Exit(1)
	}
	
	if err := generateScript(); err != nil {
		fmt.Printf("Error generating hook script: %s\n", err.Error())
		os.Exit(1)
	}
	fmt.Println("Hook script generated.")
}


func generateScript() error {
	// Define the hook script as a string
    hookScript := `#!/bin/bash

# Get the commit message file path
COMMIT_MSG_FILE=$1


# run aicommit
aicommit $COMMIT_MSG_FILE
	`

	err := ioutil.WriteFile(".git/hooks/prepare-commit-msg", []byte(hookScript), 0755)
    if err != nil {
        return err
    }

    return nil
}

