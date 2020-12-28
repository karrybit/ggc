package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

type gitHubClient struct {
	client *github.Client
}

func NewGitHubClientWithAccessToken(accessToken string) *gitHubClient {
	ctx := context.Background()
	tokenSource := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: accessToken},
	)
	oauthClient := oauth2.NewClient(ctx, tokenSource)
	client := github.NewClient(oauthClient)
	return &gitHubClient{client: client}
}

func NewGitHubClientWithBasic(username string, password string) (*gitHubClient, error) {
	fmt.Print("one time token> ")
	if err := bufio.NewWriter(os.Stdout).Flush(); err != nil {
		return nil, err
	}

	var oneTimePassword string
	if _, err := fmt.Scan(&oneTimePassword); err != nil {
		return nil, err
	}

	transport := github.BasicAuthTransport{
		Username: strings.TrimSpace(username),
		Password: strings.TrimSpace(password),
		OTP:      oneTimePassword,
	}

	client := github.NewClient(transport.Client())
	return &gitHubClient{client: client}, nil
}

func (gc *gitHubClient) createPR(
	ctx context.Context,
	prTitle string,
	prDescription string,
	commitBranch string,
	prBranch string,
	prRepoOwnerName string,
	prRepoName string,
) error {
	wrapped, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	newPR := &github.NewPullRequest{
		Title:               &prTitle,
		Head:                &commitBranch,
		Base:                &prBranch,
		Body:                &prDescription,
		MaintainerCanModify: github.Bool(true),
	}

	pr, _, err := gc.client.PullRequests.Create(wrapped, prRepoOwnerName, prRepoName, newPR)
	if err != nil {
		return err
	}

	fmt.Printf("PR created: %s\n", pr.GetHTMLURL())
	return nil
}
