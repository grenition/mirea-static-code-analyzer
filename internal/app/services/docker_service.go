package services

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

type DockerService struct {
	containerName string
}

func NewDockerService() *DockerService {
	containerName := os.Getenv("LINTER_CONTAINER_NAME")
	if containerName == "" {
		containerName = "static_analyzer_linters"
	}
	return &DockerService{
		containerName: containerName,
	}
}

func (d *DockerService) RunLinter(ctx context.Context, linterCmd []string, workDir string) (string, string, error) {
	// Prepare docker exec command
	cmd := exec.CommandContext(ctx, "docker", append([]string{
		"exec",
		"-i",
		"-w", workDir,
		d.containerName,
	}, linterCmd...)...)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	return stdout.String(), stderr.String(), err
}

func (d *DockerService) EnsureContainerRunning() error {
	// Check if container is running
	cmd := exec.Command("docker", "ps", "--filter", fmt.Sprintf("name=%s", d.containerName), "--format", "{{.Names}}")
	output, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("failed to check container status: %w", err)
	}

	if strings.TrimSpace(string(output)) != d.containerName {
		return fmt.Errorf("container %s is not running. Please start it with: docker-compose up -d linters", d.containerName)
	}

	return nil
}

