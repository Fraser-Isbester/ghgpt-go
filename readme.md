# GHGPT-GO (GitHub GPT - Go)

This is a simple Golang program that automates the creation of GitHub pull requests. It uses both the Go GitHub client and the OpenAI Language Model API to generate Pull Request titles and descriptions based on the `git diff` from your local repository.

## Dependencies

This project uses the following Golang libraries:
* `github.com/google/go-github/v41/github`
* `github.com/tmc/langchaingo/llms/openai`
* `golang.org/x/oauth2`

## Features

The program does the following:

* Authenticates with GitHub using a personal access token.
* Fetches the GitHub username and repository name from the local git configuration.
* Pushes the current branch to the remote repository.
* Retrieves the `git diff` between the current branch and `origin/main`.
* Generates a PR title and body using the OpenAI Language Model.
* Creates a new pull request with the generated title and body.

## Usage

You can run this program as follows:

```bash
go run main.go
```

## Notes

* You'll need to have GitHub CLI (`gh`) installed and authenticated on your local machine.
* Ensure that the local repository has an `origin` remote that corresponds to a GitHub repository.
* The program uses the OpenAI API for generating the pull request title and body. You will need to supply your OpenAI API key and ensure that the `langchaingo/llms/openai` package is configured correctly.
* Errors from external commands are logged and cause the program to exit.
