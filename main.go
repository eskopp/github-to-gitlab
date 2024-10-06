package main

import (
	"encoding/base64"
	"fmt"
	"log"
	"os"
	"os/exec"

	git "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
)

// Helper function to decode Base64 strings
func decodeBase64(encodedStr string) (string, error) {
	decodedBytes, err := base64.StdEncoding.DecodeString(encodedStr)
	if err != nil {
		return "", err
	}
	return string(decodedBytes), nil
}

func main() {
	// Get environment variables for Git configuration and authentication
	gitUsername := os.Getenv("INPUT_GIT_USERNAME")
	gitEmail := os.Getenv("INPUT_GIT_EMAIL")
	gitlabRepo := os.Getenv("INPUT_GITLAB_REPO")
	gitlabToken := os.Getenv("INPUT_GITLAB_TOKEN")
	githubToken := os.Getenv("INPUT_GITHUB_TOKEN")
	base64Flag := os.Getenv("INPUT_BASE64")

	// Validate environment variables
	if gitUsername == "" || gitEmail == "" || gitlabRepo == "" || gitlabToken == "" || githubToken == "" {
		log.Fatal("Environment variables INPUT_GIT_USERNAME, INPUT_GIT_EMAIL, INPUT_GITLAB_REPO, INPUT_GITLAB_TOKEN, or INPUT_GITHUB_TOKEN are not set")
	}

	// Decode the username, email, and repo URL from Base64 if the base64Flag is set to true
	if base64Flag == "true" {
		fmt.Println("Base64 flag is true, decoding username, email, and repository URL...")
		var err error
		gitUsername, err = decodeBase64(gitUsername)
		if err != nil {
			log.Fatalf("Failed to decode Git username from Base64: %s", err)
		}
		gitEmail, err = decodeBase64(gitEmail)
		if err != nil {
			log.Fatalf("Failed to decode Git email from Base64: %s", err)
		}
		gitlabRepo, err = decodeBase64(gitlabRepo)
		if err != nil {
			log.Fatalf("Failed to decode GitLab repository URL from Base64: %s", err)
		}
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

	// Authenticate using the GitHub token (without encoding/decoding)
	auth := &http.BasicAuth{
		Username: gitUsername,
		Password: githubToken, // Use the GitHub token for authentication
	}

	// Clone the GitHub repository to a temporary directory
	fmt.Println("Cloning GitHub repository...")
	repo, err := git.PlainClone("./repo", false, &git.CloneOptions{
		URL:      "https://github.com/yourusername/yourrepo.git", // Replace with actual GitHub URL
		Auth:     auth,                                           // Pass the auth struct for authentication
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

	// Push all branches to GitLab
	fmt.Println("Pushing all branches to GitLab...")
	err = repo.Push(&git.PushOptions{
		RemoteName: "gitlab",
		Auth: &http.BasicAuth{
			Username: "gitlab-ci-token", // GitLab token username
			Password: gitlabToken,       // GitLab token for authentication
		},
		Progress: os.Stdout,
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
		Auth: &http.BasicAuth{
			Username: "gitlab-ci-token",
			Password: gitlabToken,
		},
		Progress: os.Stdout,
	})
	if err != nil && err != git.NoErrAlreadyUpToDate {
		log.Fatalf("Failed to push tags to GitLab: %s", err)
	}

	fmt.Println("Mirroring complete!")
}
