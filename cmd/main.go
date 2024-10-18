package main

import (
	"context"
	"log"
	"time"

	"autopuller/docker"
	"autopuller/github"
	"autopuller/logger"
)

func checkForUpdates(ctx context.Context, gitHub github.GitHubAPI, dockerMgr docker.DockerManager) error {
	// Get the master commit from GitHub
	masterSum, err := gitHub.GetMasterSum(ctx)
	if err != nil {
		return err
	}

	// Get the current commit (locally)
	currentSum, err := gitHub.GetCurrentSum()
	if err != nil {
		return err
	}

	// Check if there's a new commit
	if masterSum != currentSum {
		log.Printf("Differences found between master (%s) and current (%s)", masterSum, currentSum)

		// Check if the last GitHub Actions run succeeded
		if passed, err := gitHub.CheckLastRun(ctx, masterSum); err == nil && passed {
			log.Println("Last run passed, proceeding with restart...")

			// Restart services using Docker Compose
			err = dockerMgr.RestartServices(ctx)
			if err != nil {
				return err
			}

			// Update the local record of the latest commit
			err = gitHub.UpdateCurrentSum(masterSum)
			if err != nil {
				return err
			}
		} else {
			log.Println("Last run failed or not completed yet. Skipping restart.")
		}
	} else {
		log.Println("No differences found. Nothing to do.")
	}

	return nil
}

func main() {
	// Initialize logger
	logger.InitLogger()

	// Context for Docker and GitHub operations
	ctx := context.Background()

	// Create real GitHub and Docker implementations
	gitHub := &github.RealGitHubAPI{}
	dockerMgr := &docker.RealDockerManager{}

	// Main loop
	for {
		err := checkForUpdates(ctx, gitHub, dockerMgr)
		if err != nil {
			log.Fatalf("Error in checking updates: %v", err)
		}

		// Sleep between checks
		interval := 60 // Hardcoded for simplicity, you can load this from env if needed
		log.Printf("Sleeping for %d seconds...", interval)
		time.Sleep(time.Duration(interval) * time.Second)
	}
}
