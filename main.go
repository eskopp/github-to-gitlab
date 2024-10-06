package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"

	git "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
)

func main() {
	// Get environment variables for Git configuration and authentication
	gitUsername := os.Getenv("INPUT_GIT_USERNAME")
	gitEmail := os.Getenv("INPUT_GIT_EMAIL")
	gitlabRepo := os.Getenv("INPUT_GITLAB_REPO")
	gitlabToken := os.Getenv("INPUT_GITLAB_TOKEN")

	// Validate environment variables
	if gitUsername == "" || gitEmail == "" || gitlabRepo == "" || gitlabToken == "" {
		log.Fatal("Environment variables INPUT_GIT_USERNAME, INPUT_GIT_EMAIL, INPUT_GITLAB_REPO, or INPUT_GITLAB_TOKEN are not set")
	}

	// Configure Git username and email
	fmt.Println("Configuring Git username and email...")
	err := exec.Command("git", "config", "--global", "user.name", gitUsername).Run()
	if err != nil {
		log.Fatalf("Failed to set Git username: %s", err)
	}
	err = exec.Command("git", "config", "--global", "user.email", gitEmail).Run()
	if err != nil {
		log.Fatalf("Failed to set Git email: %s", err)
	}

	// Clone the GitHub repository to a temporary directory
	fmt.Println("Cloning GitHub repository...")
	repo, err := git.PlainClone("./repo", false, &git.CloneOptions{
		URL:      "https://github.com/yourusername/yourrepo.git", // Replace with actual GitHub URL
		Progress: os.Stdout,
	})
	if err != nil {
		log.Fatalf("Failed to clone GitHub repository: %s", err)
	}

	// Add the GitLab remote
	fmt.Println("Adding GitLab remote...")
	_, err = repo.CreateRemote(&config.RemoteConfig{
		Name: "gitlab",
		URLs: []string{gitlabRepo},
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

	// Authenticate using the GitLab token
	auth := &http.BasicAuth{
		Username: "gitlab-ci-token", // This username is required for GitLab CI tokens
		Password: gitlabToken,
	}

	// Push all branches to GitLab
	fmt.Println("Pushing all branches to GitLab...")
	err = repo.Push(&git.PushOptions{
		RemoteName: "gitlab",
		Auth:       auth,
		Progress:   os.Stdout,
	})
	if err != nil && err != git.NoErrAlreadyUpToDate {
		log.Fatalf("Failed to push branches to GitLab: %s", err)
	}

	// Push all tags to GitLab
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
