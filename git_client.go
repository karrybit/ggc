package main

import (
	"fmt"
	"os"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
)

type gitClient struct {
	repository *git.Repository
	worktree   *git.Worktree
	cnf        *config.Config
}

func NewGitClient(workspacePath string) (*gitClient, error) {
	repository, err := git.PlainOpen(workspacePath)
	if err != nil {
		return nil, err
	}
	fmt.Printf("repositoy: %+v\n", repository)

	worktree, err := repository.Worktree()
	if err != nil {
		return nil, err
	}
	fmt.Printf("worktree: %+v\n", worktree)

	cnf, err := repository.Config()
	if err != nil {
		return nil, err
	}
	fmt.Printf("config: %+v\n", cnf)

	return &gitClient{
		repository: repository,
		worktree:   worktree,
		cnf:        cnf,
	}, nil
}

func (gc *gitClient) createRemote(repositoryName string) error {
	if _, err := gc.repository.CreateRemote(&config.RemoteConfig{
		Name: "origin",
		URLs: []string{"https://github.com/" + repositoryName + ".git"},
	}); err != nil {
		return err
	}
	cnf, err := gc.repository.Config()
	if err != nil {
		return err
	}
	gc.cnf = cnf
	return nil
}

func (gc *gitClient) deleteRemote() error {
	if err := gc.repository.DeleteRemote("origin"); err != nil {
		return err
	}
	cnf, err := gc.repository.Config()
	if err != nil {
		return err
	}
	gc.cnf = cnf
	return nil
}

func (gc *gitClient) setInstead(accessToken string, insteadOfURLs []string) error {
	useURL := "https://" + accessToken + "@github.com"
	for _, insteadOfURL := range insteadOfURLs {
		gc.cnf.URLs[insteadOfURL] = &config.URL{Name: useURL, InsteadOf: insteadOfURL}
	}
	if err := gc.repository.SetConfig(gc.cnf); err != nil {
		return err
	}
	cnf, err := gc.repository.Config()
	if err != nil {
		return err
	}
	gc.cnf = cnf
	return nil
}

func (gc *gitClient) checkout(branch string) error {
	return gc.worktree.Checkout(&git.CheckoutOptions{Branch: plumbing.ReferenceName("refs/heads/" + branch), Create: true})
}

func (gc *gitClient) addFiles(filePaths []string) error {
	pathsCount := len(filePaths)
	for i, filePath := range filePaths {
		fmt.Printf("[%d/%d] git add %v\n", i, pathsCount, filePath)
		blobHash, err := gc.worktree.Add(filePath)
		if err != nil {
			return err
		}
		fmt.Printf("blob hash is %s\n", blobHash)
	}
	return nil
}

func (gc *gitClient) commit(message string) error {
	commitHash, err := gc.worktree.Commit(message, &git.CommitOptions{
		All: true,
		Author: &object.Signature{
			Name:  gc.cnf.Author.Name,
			Email: gc.cnf.Author.Email,
		},
		Committer: &object.Signature{
			Name:  gc.cnf.User.Name,
			Email: gc.cnf.User.Email,
		},
	})
	if err != nil {
		return err
	}
	fmt.Printf("commit hash is %s\n", commitHash)
	return nil
}

func (gc *gitClient) pushWithUserNamePassword(password string) error {
	return gc.repository.Push(&git.PushOptions{
		RemoteName: "origin",
		Auth: &http.BasicAuth{
			Username: gc.cnf.User.Name,
			Password: password,
		},
		Progress: os.Stdout,
	})
}

func (gc *gitClient) pushWithAccessToken(accessToken string) error {
	return gc.repository.Push(&git.PushOptions{
		RemoteName: "origin",
		Auth: &http.BasicAuth{
			Username: gc.cnf.User.Name,
			Password: accessToken,
		},
		Progress: os.Stdout,
	})
}

func (gc *gitClient) pushWithSSH(pemFilePath string, password string) error {
	publicKey, err := ssh.NewPublicKeysFromFile(ssh.DefaultUsername, pemFilePath, password)
	if err != nil {
		return err
	}

	return gc.repository.Push(&git.PushOptions{
		RemoteName: "origin",
		Auth:       publicKey,
		Progress:   os.Stdout,
	})
}
