package main

import (
	"context"
	"log"
	"time"

	"autopuller/docker"
	"autopuller/env"
	"autopuller/github"
	"autopuller/logger"
)

func main() {
	// Initialize logger
	logger.InitLogger()

	// Load environment variables
	err := env.LoadEnv()
	if err != nil {
		log.Fatalf("Failed to load environment variables: %v", err)
	}
	// Context for Docker and GitHub operations
	ctx := context.Background()

	// Main loop
	for {
		log.Println("Checking for updates...")

		// Step 1: Check differences between master and current commit
		masterSum, err := github.GetMasterSum(ctx)
		if err != nil {
			log.Fatalf("Failed to fetch master commit: %v", err)
		}
		log.Printf("Current masterSum: %s", masterSum)
		currentSum, err := github.GetCurrentSum()
		if err != nil {
			log.Fatalf("Failed to read current commit: %v", err)
		}

		if masterSum != currentSum {
			log.Printf("Differences found between master (%s) and current (%s)", masterSum, currentSum)

			// Check if the last GitHub Actions run succeeded
			if passed, err := github.CheckLastRun(ctx, masterSum); err == nil && passed {
				log.Println("Last run passed, proceeding with restart...")

				// Pull the latest changes and restart services
				err = docker.RestartServices(ctx)
				if err != nil {
					log.Fatalf("Failed to restart services: %v", err)
				}

				// Update the local record of the latest commit
				err = github.UpdateCurrentSum(masterSum)
				if err != nil {
					log.Fatalf("Failed to update current commit: %v", err)
				}
			} else {
				log.Println("Last run failed or not completed yet. Skipping restart.")
			}
		} else {
			log.Println("No differences found. Nothing to do.")
		}

		// Sleep between checks
		interval := env.GetInterval()
		log.Printf("Sleeping for %d seconds...", interval)
		time.Sleep(time.Duration(interval) * time.Second)
	}
}
