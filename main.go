package main

import (
	"fmt"
	"log"
	"os"

	git "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
)

func main() {
	// Get environment variables for Git configuration
	gitlabToken := os.Getenv("GITLAB_TOKEN")
	githubURL := "https://github.com/yourusername/yourrepo.git"                                     // Replace with your GitHub repo URL
	gitlabURL := "https://gitlab-ci-token:" + gitlabToken + "@gitlab.com/yourusername/yourrepo.git" // Replace with your GitLab repo URL

	// Validate environment variables
	if gitlabToken == "" {
		log.Fatal("Environment variable GITLAB_TOKEN is not set")
	}

	// Clone the GitHub repository to a temporary directory
	fmt.Println("Cloning GitHub repository...")
	repo, err := git.PlainClone("./repo", false, &git.CloneOptions{
		URL:      githubURL,
		Progress: os.Stdout,
	})
	if err != nil {
		log.Fatalf("Failed to clone GitHub repository: %s", err)
	}

	// Add the GitLab remote
	fmt.Println("Adding GitLab remote...")
	_, err = repo.CreateRemote(&config.RemoteConfig{
		Name: "gitlab",
		URLs: []string{gitlabURL},
	})
	if err != nil {
		log.Fatalf("Failed to add GitLab remote: %s", err)
	}

	// Fetch all branches and tags from GitHub
	fmt.Println("Fetching all branches and tags from GitHub...")
	err = repo.Fetch(&git.FetchOptions{
		RemoteName: "origin",
		Progress:   os.Stdout,
	})
	if err != nil && err != git.NoErrAlreadyUpToDate {
		log.Fatalf("Failed to fetch branches: %s", err)
	}

	// Push all branches to GitLab
	fmt.Println("Pushing all branches to GitLab...")
	auth := &http.BasicAuth{
		Username: "gitlab-ci-token",
		Password: gitlabToken,
	}
	err = repo.Push(&git.PushOptions{
		RemoteName: "gitlab",
		Auth:       auth,
		Progress:   os.Stdout,
	})
	if err != nil && err != git.NoErrAlreadyUpToDate {
		log.Fatalf("Failed to push branches to GitLab: %s", err)
	}

	// Push all tags to GitLab using RefSpec as a string
	fmt.Println("Pushing all tags to GitLab...")
	err = repo.Push(&git.PushOptions{
		RemoteName: "gitlab",
		RefSpecs: []config.RefSpec{
			"refs/tags/*:refs/tags/*",
		},
		Auth:     auth,
		Progress: os.Stdout,
	})
	if err != nil && err != git.NoErrAlreadyUpToDate {
		log.Fatalf("Failed to push tags to GitLab: %s", err)
	}

	fmt.Println("Mirroring complete!")
}
