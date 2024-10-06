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
	// Holen Sie sich die Umgebungsvariablen mit dem INPUT_ Präfix
	githubURL := os.Getenv("INPUT_GITHUB_REPO_URL")
	gitlabURL := os.Getenv("INPUT_GITLAB_REPO_URL")
	gitlabToken := os.Getenv("INPUT_GITLAB_TOKEN")

	// Überprüfen Sie, ob die erforderlichen Umgebungsvariablen gesetzt sind
	if githubURL == "" || gitlabURL == "" || gitlabToken == "" {
		log.Fatal("Umgebungsvariablen INPUT_GITHUB_REPO_URL, INPUT_GITLAB_REPO_URL oder INPUT_GITLAB_TOKEN sind nicht gesetzt")
	}

	// Klonen Sie das GitHub-Repository in ein temporäres Verzeichnis
	fmt.Println("Cloning GitHub repository...")
	repo, err := git.PlainClone("./repo", false, &git.CloneOptions{
		URL:      githubURL,
		Progress: os.Stdout,
	})
	if err != nil {
		log.Fatalf("Failed to clone GitHub repository: %s", err)
	}

	// Fügen Sie das GitLab-Remote hinzu
	fmt.Println("Adding GitLab remote...")
	_, err = repo.CreateRemote(&config.RemoteConfig{
		Name: "gitlab",
		URLs: []string{gitlabURL},
	})
	if err != nil {
		log.Fatalf("Failed to add GitLab remote: %s", err)
	}

	// Holen Sie alle Branches und Tags von GitHub
	fmt.Println("Fetching all branches and tags from GitHub...")
	err = repo.Fetch(&git.FetchOptions{
		RemoteName: "origin",
		Progress:   os.Stdout,
		// Hier könnten Sie weitere Optionen hinzufügen, falls benötigt
	})
	if err != nil && err != git.NoErrAlreadyUpToDate {
		log.Fatalf("Failed to fetch branches: %s", err)
	}

	// Authentifizierung für GitLab
	auth := &http.BasicAuth{
		Username: "gitlab-ci-token", // Dieser Benutzername ist für GitLab-Tokens erforderlich
		Password: gitlabToken,
	}

	// Pushen Sie alle Branches zu GitLab
	fmt.Println("Pushing all branches to GitLab...")
	err = repo.Push(&git.PushOptions{
		RemoteName: "gitlab",
		Auth:       auth,
		Progress:   os.Stdout,
	})
	if err != nil && err != git.NoErrAlreadyUpToDate {
		log.Fatalf("Failed to push branches to GitLab: %s", err)
	}

	// Pushen Sie alle Tags zu GitLab
	fmt.Println("Pushing all tags to GitLab...")
	err = repo.Push(&git.PushOptions{
		RemoteName: "gitlab",
		RefSpecs: []config.RefSpec{
			config.RefSpec("refs/tags/*:refs/tags/*"),
		},
		Auth:     auth,
		Progress: os.Stdout,
	})
	if err != nil && err != git.NoErrAlreadyUpToDate {
		log.Fatalf("Failed to push tags to GitLab: %s", err)
	}

	fmt.Println("Mirroring complete!")
}
