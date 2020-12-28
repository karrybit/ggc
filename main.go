package main

import (
	"fmt"
	"os"

	"github.com/go-git/go-git/v5"
)

func main() {
	repository, err := git.PlainOpen(".")
	if err != nil {
		panic(err)
	}
	fmt.Printf("repositoy:\n%+v\n", repository)

	worktree, err := repository.Worktree()
	if err != nil {
		panic(err)
	}
	fmt.Printf("worktree:\n%+v\n", worktree)

	if err := worktree.Checkout(&git.CheckoutOptions{Branch: "refs/heads/hoge", Create: false}); err != nil {
		panic(err)
	}

	blobHash, err := worktree.Add(".")
	if err != nil {
		panic(err)
	}
	fmt.Printf("blob hash is %s\n", blobHash)

	commitHash, err := worktree.Commit("message", &git.CommitOptions{All: true})
	if err != nil {
		panic(err)
	}
	fmt.Printf("commit hash is %s\n", commitHash)

	if err := repository.Push(&git.PushOptions{Progress: os.Stdout}); err != nil {
		panic(err)
	}
}
