#!/bin/bash
set -e

# Assign GitHub Action inputs to local variables
git_username="$INPUT_GIT_USERNAME"
git_email="$INPUT_GIT_EMAIL"
gitlab_repo="$INPUT_GITLAB_REPO"
gitlab_token="$INPUT_GITLAB_TOKEN"  # Not printed for security reasons

# Output the values of the variables (excluding the token)
echo "::debug::git_username: $git_username"
echo "::debug::git_email: $git_email"
echo "::debug::gitlab_repo (before check): $gitlab_repo"

# Check if required environment variables are set, using GitHub Actions error annotation
if [ -z "$git_username" ]; then
  echo "::error::git_username is not set."
  exit 1
fi

if [ -z "$git_email" ]; then
  echo "::error::git_email is not set."
  exit 1
fi

if [ -z "$gitlab_repo" ]; then
  echo "::error::gitlab_repo is not set."
  exit 1
fi

if [ -z "$gitlab_token" ]; then
  echo "::error::gitlab_token is not set."
  exit 1
fi

# Mark the GitHub Actions workspace as a safe directory
git config --global --add safe.directory /github/workspace

# Extract the branch name from the GitHub reference
branch_name=$(echo "${GITHUB_REF#refs/heads/}")
echo "::debug::branch_name: $branch_name"

# Ensure the GitLab repository URL does not have duplicate https://
if [[ "$gitlab_repo" == https://* ]]; then
  full_gitlab_repo="$gitlab_repo"
else
  full_gitlab_repo="https://$gitlab_repo"
fi

# Output the full GitLab repository URL (debugging)
echo "::debug::full_gitlab_repo (after check): $full_gitlab_repo"

# Configure Git with error handling
echo "::debug::Configuring Git with username: $git_username and email: $git_email"
if ! git config --global user.name "$git_username"; then
  echo "::error::Failed to configure Git username."
  exit 1
fi

if ! git config --global user.email "$git_email"; then
  echo "::error::Failed to configure Git email."
  exit 1
fi

# Add GitLab remote using the GitLab token for authentication, with error handling
echo "::debug::Adding GitLab remote for repository: $full_gitlab_repo"
if ! git remote add gitlab "https://gitlab-ci-token:${gitlab_token}@${full_gitlab_repo}"; then
  echo "::error::Failed to add GitLab remote."
  exit 1
fi

# Debugging: List the remotes to verify the remote URL
echo "::debug::Listing remotes:"
git remote -v

# Fetch from GitLab with error handling
echo "::debug::Fetching from GitLab repository"
if ! git fetch gitlab; then
  echo "::error::Failed to fetch from GitLab."
  echo "::debug::Retrying fetch command with verbose output"
  git fetch gitlab --verbose
  exit 1
fi

# Merge changes from GitLab, allowing unrelated histories, with error handling
echo "::debug::Merging changes from branch: $branch_name"
if ! git merge "gitlab/$branch_name" --allow-unrelated-histories || true; then
  echo "::error::Failed to merge changes from GitLab."
  exit 1
fi

# Push changes to GitLab with error handling
echo "::debug::Pushing changes to GitLab branch: $branch_name"
if ! git push gitlab "refs/heads/$branch_name:refs/heads/$branch_name"; then
  echo "::error::Failed to push changes to GitLab."
  exit 1
fi

# Push tags to GitLab with error handling
echo "::debug::Pushing tags to GitLab"
if ! git push --tags gitlab; then
  echo "::error::Failed to push tags to GitLab."
  exit 1
fi
