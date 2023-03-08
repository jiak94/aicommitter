# AI Committer

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

AI Committer is a command-line tool that uses OpenAI's GPT-3 language model to generate short git commit messages based on the changes you made to your code.

## Installation

### From Source

To install AI Committer, you need to have Go installed on your system. Then run:

```bash
go install github.com/jiak94/aicommitter@latest
```

### Download Prebuilt Binary

Prebuilt binary are available [GitHub releases page](https://github.com/jiak94/aicommitter/releases)

Choose the appropriate binary for your operating system and architecture:

`aicommit-darwin-amd64`: for macOS using Intel chipset

`aicommit-darwin-arm64`: for macOS using M1/M2

`aicommit-linux-amd64`: for linux using amd64 architecture

`aicommit-linux-arm64`: for linux using arm64 architecture

Please remember to rename it to `aicommit` and put it into `$PATH` folder

## Configuration

AI Committer needs an OpenAI API key to work. You can create a free account and get an API key from the [OpenAI website](https://beta.openai.com/signup/).

Once you have an API key, you need to create a configuration file for AI Committer.

### Show Config

To show the configuration file for AI Committer, run the following command:

```bash
aicommit config --show
```

### Set API Key

To set your OpenAI API key, run the following command:

```bash
aicommit config --api-key <your-apk-key>
```

### Set Model

To set the OpenAI GPT-3 language model to use, run the following command:

```bash
aicommit config --model <model_name>
```

## Usage

### Register AI committer as [`prepare-commit-msg`](https://git-scm.com/docs/githooks#_prepare_commit_msg) hook

Navigate to the root directory of your git repository. Then run the following command:

```bash
aicommit registerHook
```

After registering the `prepare-commit-msg` hook, the AI committer will generate a commit message for you every time you run `git commit`.

### Show AI generated commit message

Add the changes to staging area. Then run the following command:

```bash
aicommit
```

AI Committer uses the difference between your latest commit and the current state of your code in the staging area to generate a concise commit message when it is run.

## License

AI Committer is released under the [MIT License](LICENSE).
