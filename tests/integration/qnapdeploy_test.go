package integration

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"testing"
)

// TestQNAPDockerEndToEnd tests the full qnap-docker command-line tool
func TestQNAPDockerEndToEnd(t *testing.T) {
	// Skip in CI environments
	if os.Getenv("CI") != "" || os.Getenv("GITHUB_ACTIONS") != "" {
		t.Skip("Skipping local network test in CI environment")
	}
	// Skip in short mode (unit tests only)
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// Test qnap-docker version
	t.Run("VersionCommand", func(t *testing.T) {
		cmd := exec.Command("../../bin/qnap-docker", "--version")
		output, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatalf("qnap-docker --version failed: %v\nOutput: %s", err, output)
		}

		if !strings.Contains(string(output), "qnap-docker version") {
			t.Fatalf("Invalid version output: %s", output)
		}

		t.Logf("✅ qnap-docker version: %s", strings.TrimSpace(string(output)))
	})

	// Test qnap-docker help
	t.Run("HelpCommand", func(t *testing.T) {
		cmd := exec.Command("../../bin/qnap-docker", "--help")
		output, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatalf("qnap-docker --help failed: %v\nOutput: %s", err, output)
		}

		requiredSections := []string{
			"CLI tool that simplifies Docker container deployment",
			"Available Commands:",
			"init",
			"run",
			"deploy",
			"ps",
			"rm",
		}

		for _, section := range requiredSections {
			if !strings.Contains(string(output), section) {
				t.Errorf("Help output missing section: %s", section)
			}
		}

		t.Logf("✅ qnap-docker help output complete")
	})

	// Test qnap-docker init help
	t.Run("InitHelpCommand", func(t *testing.T) {
		cmd := exec.Command("../../bin/qnap-docker", "init", "--help")
		output, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatalf("qnap-docker init --help failed: %v\nOutput: %s", err, output)
		}

		requiredFlags := []string{
			"--user",
			"--port",
			"--key",
			"--volume-path",
		}

		for _, flag := range requiredFlags {
			if !strings.Contains(string(output), flag) {
				t.Errorf("Init help missing flag: %s", flag)
			}
		}

		// Check QNAP-specific defaults
		if !strings.Contains(string(output), "/share/CACHEDEV1_DATA/docker") {
			t.Errorf("Init help missing QNAP default volume path")
		}

		t.Logf("✅ qnap-docker init help output complete")
	})

	// Test qnap-docker run help
	t.Run("RunHelpCommand", func(t *testing.T) {
		cmd := exec.Command("../../bin/qnap-docker", "run", "--help")
		output, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatalf("qnap-docker run --help failed: %v\nOutput: %s", err, output)
		}

		requiredFlags := []string{
			"--name",
			"--port",
			"--volume",
			"--env",
			"--restart",
			"--network",
		}

		for _, flag := range requiredFlags {
			if !strings.Contains(string(output), flag) {
				t.Errorf("Run help missing flag: %s", flag)
			}
		}

		t.Logf("✅ qnap-docker run help output complete")
	})
}

// TestDirectDockerCommands tests Docker commands directly via SSH to QNAP
func TestDirectDockerCommands(t *testing.T) {
	// Skip in CI environments
	if os.Getenv("CI") != "" || os.Getenv("GITHUB_ACTIONS") != "" {
		t.Skip("Skipping local network test in CI environment")
	}
	// Skip in short mode (unit tests only)
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// Test Docker info command
	t.Run("DockerInfo", func(t *testing.T) {
		cmd := exec.Command("ssh", "scttfrdmn@astrapi.local",
			"/share/CACHEDEV1_DATA/.qpkg/container-station/bin/docker info --format '{{.ServerVersion}}'")
		output, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatalf("Docker info failed: %v\nOutput: %s", err, output)
		}

		version := strings.TrimSpace(string(output))
		if version == "" {
			t.Fatal("Docker info returned empty version")
		}

		t.Logf("✅ Docker server version: %s", version)
	})

	// Test image pulling (using a small image)
	t.Run("ImagePull", func(t *testing.T) {
		cmd := exec.Command("ssh", "scttfrdmn@astrapi.local",
			"/share/CACHEDEV1_DATA/.qpkg/container-station/bin/docker pull hello-world:latest")
		output, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatalf("Docker pull failed: %v\nOutput: %s", err, output)
		}

		if !strings.Contains(string(output), "Pull complete") && !strings.Contains(string(output), "Image is up to date") {
			t.Fatalf("Unexpected pull output: %s", output)
		}

		t.Logf("✅ Docker image pull successful")
	})

	// Test container run (quick test)
	t.Run("ContainerRun", func(t *testing.T) {
		containerName := fmt.Sprintf("qnap-docker-test-%d", os.Getpid())

		// Cleanup any existing container with same name
		cleanupCmd := exec.Command("ssh", "scttfrdmn@astrapi.local",
			fmt.Sprintf("/share/CACHEDEV1_DATA/.qpkg/container-station/bin/docker rm -f %s 2>/dev/null || true", containerName))
		cleanupCmd.Run()

		// Run hello-world container
		cmd := exec.Command("ssh", "scttfrdmn@astrapi.local",
			fmt.Sprintf("/share/CACHEDEV1_DATA/.qpkg/container-station/bin/docker run --name %s hello-world:latest", containerName))
		output, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatalf("Docker run failed: %v\nOutput: %s", err, output)
		}

		if !strings.Contains(string(output), "Hello from Docker!") {
			t.Fatalf("Hello-world container didn't run correctly: %s", output)
		}

		// Cleanup
		cleanupCmd = exec.Command("ssh", "scttfrdmn@astrapi.local",
			fmt.Sprintf("/share/CACHEDEV1_DATA/.qpkg/container-station/bin/docker rm -f %s", containerName))
		cleanupCmd.Run()

		t.Logf("✅ Docker container run/remove successful")
	})

	// Test CACHEDEV path validation
	t.Run("CacheDevPathValidation", func(t *testing.T) {
		// Test multiple potential CACHEDEV paths
		potentialPaths := []string{
			"/share/CACHEDEV1_DATA",
			"/share/CACHEDEV2_DATA",
			"/share/CACHEDEV3_DATA",
		}

		foundPaths := []string{}
		for _, path := range potentialPaths {
			cmd := exec.Command("ssh", "scttfrdmn@astrapi.local", fmt.Sprintf("test -d %s && echo 'EXISTS' || echo 'NOT_FOUND'", path))
			output, err := cmd.CombinedOutput()
			if err == nil && strings.Contains(string(output), "EXISTS") {
				foundPaths = append(foundPaths, path)
			}
		}

		if len(foundPaths) == 0 {
			t.Fatalf("No CACHEDEV paths found on QNAP NAS")
		}

		t.Logf("✅ Found CACHEDEV paths: %v", foundPaths)
	})
}