# GitHub to GitLab Sync Action

This GitHub Action synchronizes all branches from a GitHub repository to a specified GitLab repository. It uses a Docker container that has Git pre-installed, ensuring compatibility across environments.

## Inputs

The Action requires the following inputs:

| Input          | Description                          | Required |
| -------------- | ------------------------------------ | -------- |
| `git_username` | Git username to configure Git.       | true     |
| `git_email`    | Git email to configure Git.          | true     |
| `gitlab_repo`  | The GitLab repository URL.           | true     |
| `gitlab_token` | GitLab token for authentication.     | true     |

## Example Usage

```yaml
name: Sync all branches to GitLab

on:
  push:
    branches:
      - '**'

jobs:
  sync:
    runs-on: ubuntu-latest

    steps:
    - name: Checkout code from GitHub
      uses: actions/checkout@v4
      with:
        fetch-depth: 0

    - name: Run GitHub to GitLab Sync Action
      uses: eskopp/github-to-gitlab-sync@v0.0.3
      with:
        git_username: "Your Git Username"
        git_email: "your-email@example.com"
        gitlab_repo: "gitlab.example.com/your-repo.git"
        gitlab_token: ${{ secrets.GITLAB }}

```
# Setup
1. Ensure you have created a GitLab personal access token and stored it in your GitHub repository as a secret called ``GITLAB``.

2. In your workflow, use this action by referencing its location and passing the required inputs.


# License
This project is licensed under the MIT License. [MIT-Licence](LICENSE)