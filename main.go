package main

import (
	"context"
	"fmt"
	"os"
)

func main() {
	_, err := NewGitClient(".")
	if err != nil {
		panic(err)
	}
	token := os.Getenv("GITHUB_AUTH_TOKEN")
	ghc := NewGitHubClientWithAccessToken(token)
	ctx := context.Background()
	if err := ghc.repositoryList(ctx); err != nil {
		panic(err)
	}
	if err := ghc.createPR(ctx); err != nil {
		panic(err)
	}
	fmt.Println("finish")
}
