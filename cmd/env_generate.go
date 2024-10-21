package main

import (
	"fmt"
	"os"
	"path/filepath"
)

// Content for the .env.sample file
const envSampleContent = `# .env.sample

# GitHub API token with permissions to access the repository
GITHUBKEY=

# GitHub repository name, without the www. or https://github.com/
# Example: if the repository is https://github.com/user/repo, set this to user/repo
REPONAME=amunchet/autopuller-go

# Directory where the GitHub repository is cloned locally
# Example: /path/to/local/repo
REPODIR=.

# Directory where Docker Compose is configured
# Example: /path/to/docker-compose
DOCKERDIR=./docker/sample

# Interval in seconds between checks for new commits (default: 60 seconds)
INTERVAL=60

# Optional: Command for sending email notifications (default: 'mail -s')
SENDMAIL_CMD=mail -s

# Optional: Commit message used for automatic linting fixes (default: 'Automatic linting fix')
LINTING_COMMIT_MSG=Automatic linting fix

# Optional: Force pulling new images when running docker-compose (set to any value to enable)
# If set to any value, it will add --pull to docker-compose build
FORCEPULL=
`

// GenerateEnvSample generates the .env.sample file in the current directory
func GenerateEnvSample() error {
	// Get the current working directory
	currentDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("could not get current directory: %v", err)
	}

	// Define the file path for the .env.sample file
	filePath := filepath.Join(currentDir, ".env.sample")

	// Create or open the file
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("could not create .env.sample file: %v", err)
	}
	defer file.Close()

	// Write the content to the .env.sample file
	_, err = file.WriteString(envSampleContent)
	if err != nil {
		return fmt.Errorf("could not write to .env.sample file: %v", err)
	}

	fmt.Printf(".env.sample file created at: %s\n", filePath)
	return nil
}

