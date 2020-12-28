package main

import (
	"context"
	"fmt"
	"os"
	"strings"
)

func main() {
	gc, err := NewGitClient(".")
	if err != nil {
		panic(err)
	}
	ref, err := gc.repository.Head()
	if err != nil {
		panic(err)
	}
	refs := strings.Split(ref.Name().String(), "/")

	token := os.Getenv("GITHUB_AUTH_TOKEN")
	ghc := NewGitHubClientWithAccessToken(token)
	ctx := context.Background()
	if err := ghc.createPR(
		ctx,
		"title",
		"description",
		refs[len(refs)-1],
		"main",
		gc.cnf.User.Name,
		"ggc",
	); err != nil {
		panic(err)
	}
	fmt.Println("finish")
}
