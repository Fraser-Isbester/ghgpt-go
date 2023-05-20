package main

import (
	"context"
	"log"
	"os/exec"
	"strings"

	"github.com/google/go-github/v41/github"
	"golang.org/x/oauth2"
)

func main() {
	ctx := context.Background()
	token := getToken()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)

	owner := "fraser-isbester"
	repo := "ghgpt-go"

	pr := &github.NewPullRequest{
		Title: github.String("Test PR"),
		Body:  github.String("This is a test PR"),
		Draft: github.Bool(true),
	}

	client.PullRequests.Create(ctx, owner, repo, pr)

	// // list all repositories for the authenticated user
	// repos, _, err := client.Repositories.List(ctx, "", nil)
	// if err != nil {
	// 	fmt.Printf("Error: %v\n", err)
	// 	return
	// }

	// for _, repo := range repos {
	// 	fmt.Printf("Name: %s, URL: %s\n", *repo.Name, *repo.HTMLURL)
	// }
}

func getToken() string {
	cmd := exec.Command("gh", "auth", "token")

	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatalf("cmd.Run() failed with %s\n", err)
	}
	return strings.TrimSpace(string(output))
}
