package main

import (
	"fmt"
	"os"
	"path/filepath"
	"text/template"
)

// ServiceFileTemplate is a systemd service file template
const ServiceFileTemplate = `[Unit]
Description={{.ServiceName}} Service
After=network.target

[Service]
Type=simple
WorkingDirectory={{.WorkingDir}}
ExecStart={{.ExecPath}}
Restart=always
RestartSec=10
User={{.User}}

[Install]
WantedBy=multi-user.target
`

// SystemdServiceConfig holds the information needed to generate a systemd service file
type SystemdServiceConfig struct {
	ServiceName string
	ExecPath    string
	WorkingDir  string
	User        string
}

// GenerateSystemdService creates and saves a systemd service file based on the current directory
func GenerateSystemdService(serviceName, user string) error {
	// Get the current working directory
	workingDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("could not get working directory: %v", err)
	}

	// Get the current executable path
	execPath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("could not get current executable path: %v", err)
	}

	// Create a SystemdServiceConfig with the gathered information
	config := SystemdServiceConfig{
		ServiceName: serviceName,
		ExecPath:    execPath,
		WorkingDir:  workingDir,
		User:        user,
	}

	// Define the output file path for the systemd service
	outputPath := filepath.Join("/etc/systemd/system", serviceName+".service")

	// Create or open the file
	file, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("could not create service file: %v", err)
	}
	defer file.Close()

	// Parse and execute the template
	tmpl, err := template.New("systemdService").Parse(ServiceFileTemplate)
	if err != nil {
		return fmt.Errorf("could not parse template: %v", err)
	}

	// Write the generated systemd service to the file
	err = tmpl.Execute(file, config)
	if err != nil {
		return fmt.Errorf("could not execute template: %v", err)
	}

	// Set proper permissions for the service file
	err = os.Chmod(outputPath, 0644)
	if err != nil {
		return fmt.Errorf("could not set file permissions: %v", err)
	}

	fmt.Printf("Systemd service file created: %s\n", outputPath)
	return nil
}
