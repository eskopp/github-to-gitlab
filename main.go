package main

import (
	"fmt"
	"log"
	"os"

	git "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
)

func main() {
	// Get environment variables for Git configuration
	gitUsername := os.Getenv("GIT_USERNAME")
	gitEmail := os.Getenv("GIT_EMAIL")
	gitlabToken := os.Getenv("GITLAB_TOKEN")
	githubURL := "https://github.com/yourusername/yourrepo.git"                                     // Replace with your GitHub repo URL
	gitlabURL := "https://gitlab-ci-token:" + gitlabToken + "@gitlab.com/yourusername/yourrepo.git" // Replace with your GitLab repo URL

	// Validate environment variables
	if gitUsername == "" || gitEmail == "" || gitlabToken == "" {
		log.Fatal("Environment variables GIT_USERNAME, GIT_EMAIL, or GITLAB_TOKEN are not set")
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

	// Configure Git user
	fmt.Println("Configuring Git username and email...")
	config, err := repo.Config()
	if err != nil {
		log.Fatalf("Failed to retrieve Git config: %s", err)
	}
	config.User.Name = gitUsername
	config.User.Email = gitEmail
	err = repo.Storer.SetConfig(config)
	if err != nil {
		log.Fatalf("Failed to set Git config: %s", err)
	}

	// Add the GitLab remote
	fmt.Println("Adding GitLab remote...")
	_, err = repo.CreateRemote(&git.Config{
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

	// Push all tags to GitLab
	fmt.Println("Pushing all tags to GitLab...")
	err = repo.Push(&git.PushOptions{
		RemoteName: "gitlab",
		RefSpecs: []git.RefSpec{
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
