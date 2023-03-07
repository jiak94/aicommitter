# AI Committer

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

AI Committer is a command-line tool that uses OpenAI's GPT-3 language model to generate short git commit messages based on the changes you made to your code.

## Usage

```bash
aicommit
```

When you run AI Committer, it will use the diff between your latest commit and the current state of your code to generate a short commit message.

AI Committer will prompt you to confirm or edit the generated commit message before making a commit.

## Installation

To install AI Committer, you need to have Go installed on your system. Then run:

```bash
git clone git@github.com:jiak94/aicommitter.git
cd aicommitter
go build -v -o aicommit
```

This will generate a binary called `aicommit`. After successfully compiling the binary, move it to `$PATH` folder.

## Configuration

AI Committer needs an OpenAI API key to work. You can create a free account and get an API key from the [OpenAI website](https://beta.openai.com/signup/).

Once you have an API key, you need to create a configuration file for AI Committer. By default, the configuration file should be located at `~/.config/aicommitter/config.toml`.

If the configuration file does not exist, AI Committer will create it for you and prompt you to enter your OpenAI API key.

You will need to replace the default value for the api-key field with your own API key.

## License

AI Committer is released under the [MIT License](LICENSE).
