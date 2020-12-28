package main

import (
	"context"
	"fmt"
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

func NewGitHubClientWithBasic(username string, password string, oneTimePassword *string) *gitHubClient {
	transport := github.BasicAuthTransport{
		Username: strings.TrimSpace(username),
		Password: strings.TrimSpace(password),
		OTP:      *oneTimePassword,
	}
	client := github.NewClient(transport.Client())
	return &gitHubClient{client: client}
}

func (gc *gitHubClient) repositoryList(ctx context.Context) error {
	wrapped, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	repo, resp, err := gc.client.Repositories.List(wrapped, "", nil)
	if err != nil {
		return err
	}
	fmt.Printf("repo %v\n", repo)
	fmt.Printf("resp %v\n", resp)
	return nil
}

func (gc *gitHubClient) createPR(ctx context.Context) error {
	wrapped, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	title := "title"
	commitBranch := "commitBranch"
	prBranch := "prBranch"
	prDescription := "prDescription"
	newPR := &github.NewPullRequest{
		Title:               &title,
		Head:                &commitBranch,
		Base:                &prBranch,
		Body:                &prDescription,
		MaintainerCanModify: github.Bool(true),
	}

	pr, _, err := gc.client.PullRequests.Create(wrapped, "prRepoOwner", "prRepo", newPR)
	if err != nil {
		return err
	}

	fmt.Printf("PR created: %s\n", pr.GetHTMLURL())
	return nil
}
