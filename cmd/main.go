package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"autopuller/docker"
	"autopuller/env"
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

		// Check if last run was successful
		if passed, err := gitHub.CheckLastRun(ctx, masterSum); err == nil && passed {
			log.Println("Last run passed, proceeding with update.")

			// Check for file differences
			diffs, err := gitHub.CheckDifferences(ctx, currentSum, masterSum)
			if err != nil {
				return err
			}

			if len(diffs) == 0 {
				log.Println("No files changed. Exiting.")
				return nil
			}

			// Run git pull to update the repository
			repoDir := os.Getenv("REPODIR")
			if err := gitHub.RunGitPull(ctx, repoDir); err != nil {
				return err
			}

			// Restart services using Docker Compose
			err = dockerMgr.RestartServices(ctx)
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

	for _, arg := range os.Args {
		if strings.Contains(arg, "systemd") {
			log.Println("Generating systemd file...")
			log.Println("\tRemember to run `sudo systemctl daemon-reload`")
			log.Println("\tTo enable the service: `sudo systemctl enable autopuller`")
			log.Println("\tTo start the service: `sudo systemctl start autopuller`")
			err := GenerateSystemdService("autopuller", "root")
			if err != nil {
				fmt.Printf("Error: %v\n", err)
			}
			return
		}
	}

	// Initialize logger
	logger.InitLogger("autopuller.log")

	// Load Dot Env
	if err := env.LoadEnv(); err != nil {
		log.Fatalf("Error loading ENV: %v", err)
	}

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
