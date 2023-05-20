package main

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"os/exec"
	"strings"

	"github.com/tmc/langchaingo/llms/openai"
	"golang.org/x/oauth2"

	"github.com/google/go-github/v41/github"
)

func main() {

	ctx := context.Background()
	token := getToken()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	pathParts := strings.Split(getRepo(), "/")
	gitUsername := pathParts[1]
	gitRepo := strings.TrimSuffix(pathParts[2], ".git")

	gitHead := getHead()
	gitRemote := "origin"

	// Push the branch to remote
	err := push(pushOpts{head: gitHead, remote: gitRemote})
	if err != nil {
		fmt.Println(err)
	}

	// Get the diff
	gitDiff := getDiff()

	// Create a PR
	pullRequest := getAutoPullRequest(ctx, gitDiff)
	pullRequest.Base = github.String("main")
	// pullRequest.Head = github.String(gitHead)
	pullRequest.Draft = github.Bool(true)

	newPullRequest, r, err := client.PullRequests.Create(ctx, gitUsername, gitRepo, &pullRequest)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}
	fmt.Println(newPullRequest, r)

}

// gets the github token from the gh cli
func getToken() string {
	cmd := exec.Command("gh", "auth", "token")

	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatalf("cmd.Run() failed with %s\n", err)
	}
	return strings.TrimSpace(string(output))
}

// gets the repo name from the git remote
func getRepo() string {
	cmd := exec.Command("git", "remote", "get-url", "origin")

	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatalf("cmd.Run() failed with %s\n", err)
	}
	rawId := strings.TrimSpace(string(output))
	rawId = strings.Replace(rawId, ":", "/", 1)

	u, err := url.Parse(rawId)
	if err != nil {
		log.Fatal(err)
	}
	return u.Path
}

// gets the current branch name
func getHead() string {
	cmd := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD")

	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatalf("cmd.Run() failed with %s\n", err)
	}
	return strings.TrimSpace(string(output))
}

func getDiff() string {
	cmd := exec.Command("git", "diff", "origin/main")

	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatalf("cmd.Run() failed with %s\n", err)
	}
	return strings.TrimSpace(string(output))
}

// pushOpts are the options for the push command
type pushOpts struct {
	head   string
	remote string
}

// push pushes the current branch to remote
func push(opts pushOpts) error {
	cmd := exec.Command("git", "push", "-u", opts.remote, opts.head)

	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatalf("cmd.Run() failed with %s\n", err)
		return err
	}
	fmt.Println(string(output))
	return nil
}

func getAutoPullRequest(ctx context.Context, gitDiff string) github.NewPullRequest {

	llm, err := openai.New()
	if err != nil {
		log.Fatal(err)
	}

	titleString := "Create a GitHub PR Title for the following diff:\n" + gitDiff

	completion, err := llm.Call(ctx, titleString)
	if err != nil {
		log.Fatal(err)
		panic(err)
	}

	fmt.Println("completion:", completion)

	return github.NewPullRequest{
		Title: github.String(completion),
		Body:  github.String("This is a test PR"),
	}
}
