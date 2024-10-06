#!/bin/bash
set -e

# Check if required environment variables are set, using GitHub Actions error annotation
if [ -z "$INPUT_GIT_USERNAME" ]; then
  echo "::error::git_username is not set."
  exit 1
fi

if [ -z "$INPUT_GIT_EMAIL" ]; then
  echo "::error::git_email is not set."
  exit 1
fi

if [ -z "$INPUT_GITLAB_REPO" ]; then
  echo "::error::gitlab_repo is not set."
  exit 1
fi

if [ -z "$gitlab_token" ]; then
  echo "::error::gitlab_token is not set."
  exit 1
fi

# Extract the branch name from the GitHub reference
branch_name=$(echo "${GITHUB_REF#refs/heads/}")

# Configure Git with error handling
if ! git config --global user.name "$INPUT_GIT_USERNAME"; then
  echo "::error::Failed to configure Git username."
  exit 1
fi

if ! git config --global user.email "$INPUT_GIT_EMAIL"; then
  echo "::error::Failed to configure Git email."
  exit 1
fi

# Add GitLab remote using the GitLab token for authentication, with error handling
if ! git remote add gitlab https://gitlab-ci-token:${gitlab_token}@${INPUT_GITLAB_REPO}; then
  echo "::error::Failed to add GitLab remote."
  exit 1
fi

# Fetch from GitLab with error handling
if ! git fetch gitlab; then
  echo "::error::Failed to fetch from GitLab."
  exit 1
fi

# Merge changes from GitLab, allowing unrelated histories, with error handling
if ! git merge "gitlab/$branch_name" --allow-unrelated-histories || true; then
  echo "::error::Failed to merge changes from GitLab."
  exit 1
fi

# Push changes to GitLab with error handling
if ! git push gitlab "refs/heads/$branch_name:refs/heads/$branch_name"; then
  echo "::error::Failed to push changes to GitLab."
  exit 1
fi

# Push tags to GitLab with error handling
if ! git push --tags gitlab; then
  echo "::error::Failed to push tags to GitLab."
  exit 1
fi
