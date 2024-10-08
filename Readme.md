[![Sync all branches to GitLab](https://github.com/eskopp/github-to-gitlab/actions/workflows/test.yml/badge.svg)](https://github.com/eskopp/github-to-gitlab/actions/workflows/test.yml)

# GitHub to GitLab Sync Action

This action is designed to automatically sync a GitHub repository to a GitLab repository. It ensures that any changes made to the GitHub repository are mirrored on GitLab, keeping both repositories in sync.

## Features

- Automatically syncs all branches and tags from GitHub to GitLab.
- Works with both public and private repositories.
- Supports authentication using GitHub and GitLab tokens.
- Optionally decodes Base64-encoded credentials and repository URLs for security.

## Inputs

| Input Name     | Description                                        | Required | Default |
| -------------- | -------------------------------------------------- | -------- | ------- |
| `git_username` | Git username to set up Git.                        | Yes      | N/A     |
| `git_email`    | Git email to set up Git.                           | Yes      | N/A     |
| `gitlab_repo`  | The GitLab repository URL where to sync.           | Yes      | N/A     |
| `gitlab_token` | GitLab token for authentication.                   | Yes      | N/A     |
| `github_token` | GitHub token for authentication when cloning.      | Yes      | N/A     |
| `base64`       | Whether to decode username, email, and GitLab repo | No       | `false` |

## Usage

To use this action, add the following to your GitHub Actions workflow file:

```yaml
name: Mirror GitHub to GitLab

on:
  push:
    branches:
      - main # Adjust to the branches you want to sync

jobs:
  sync:
    runs-on: ubuntu-latest

    steps:
      - name: Checkout repository
        uses: actions/checkout@v3

      - name: Sync GitHub to GitLab
        uses: eskopp/github-to-gitlab@v1
        with:
          git_username: ${{ secrets.GIT_USERNAME }}
          git_email: ${{ secrets.GIT_EMAIL }}
          gitlab_repo: ${{ secrets.GITLAB_REPO }}
          gitlab_token: ${{ secrets.GITLAB_TOKEN }}
          github_token: ${{ secrets.GITHUB_TOKEN }}
          base64: "false" # Set to true if using Base64-encoded values
```

### Input Descriptions

- **git_username**: The username used for Git configuration (must be set correctly to push changes).
- **git_email**: The email address used for Git configuration.
- **gitlab_repo**: The URL of the GitLab repository where the GitHub repository will be mirrored.
- **gitlab_token**: The GitLab access token used to authenticate when pushing branches and tags to the GitLab repository.
- **github_token**: The GitHub access token used to authenticate when cloning the GitHub repository.
- **base64**: If set to `true`, this flag will decode `git_username`, `git_email`, and `gitlab_repo` from Base64 before using them.

## Example Setup

### Prerequisites

You will need to store your GitHub and GitLab tokens as secrets in your repository settings:

- `GIT_USERNAME`: Your Git username.
- `GIT_EMAIL`: Your Git email.
- `GITLAB_REPO`: The URL of the GitLab repository.
- `GITLAB_TOKEN`: The GitLab access token.
- `GITHUB_TOKEN`: The GitHub access token (this is usually provided by GitHub Actions itself).

### Example GitHub Workflow

```yaml
name: Mirror GitHub to GitLab

on:
  push:
    branches:
      - main
      - feature/* # Sync specific branches

jobs:
  sync:
    runs-on: ubuntu-latest

    steps:
      - name: Checkout repository
        uses: actions/checkout@v3

      - name: Sync GitHub to GitLab
        uses: eskopp/github-to-gitlab@v1
        with:
          git_username: ${{ secrets.GIT_USERNAME }}
          git_email: ${{ secrets.GIT_EMAIL }}
          gitlab_repo: ${{ secrets.GITLAB_REPO }}
          gitlab_token: ${{ secrets.GITLAB_TOKEN }}
          github_token: ${{ secrets.GITHUB_TOKEN }}
          base64: "false"
```

## License

This project is licensed under the MIT License. [MIT License](./LICENSE)
