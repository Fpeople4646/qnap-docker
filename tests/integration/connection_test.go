package integration

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// TestConnectionToQNAP tests direct connection to a configurable QNAP system
func TestConnectionToQNAP(t *testing.T) {
	// Skip in CI environments (GitHub Actions can't reach local networks)
	if os.Getenv("CI") != "" || os.Getenv("GITHUB_ACTIONS") != "" {
		t.Skip("Skipping local network test in CI environment")
	}
	// Skip in short mode (unit tests only)
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// Use environment variables or flags for configuration
	nasHost := os.Getenv("QNAP_HOST")
	if nasHost == "" {
		nasHost = "astrapi.local" // fallback for local testing
	}
	nasUser := os.Getenv("QNAP_USER")
	if nasUser == "" {
		nasUser = "admin" // default user
	}

	// Test hostname resolution
	t.Run("HostnameResolution", func(t *testing.T) {
		cmd := exec.Command("ping", "-c", "1", nasHost)
		output, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatalf("Failed to ping %s: %v\nOutput: %s", nasHost, err, output)
		}
		t.Logf("✅ Hostname %s resolves successfully", nasHost)
	})

	// Test SSH connectivity
	t.Run("SSHConnectivity", func(t *testing.T) {
		sshTarget := fmt.Sprintf("%s@%s", nasUser, nasHost)
		cmd := exec.Command("ssh", "-o", "ConnectTimeout=10", "-o", "BatchMode=yes",
			sshTarget, "echo 'SSH connection successful'; whoami")
		output, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatalf("SSH connection failed: %v\nOutput: %s", err, output)
		}

		if !strings.Contains(string(output), "SSH connection successful") {
			t.Fatalf("SSH connection test failed. Output: %s", output)
		}

		t.Logf("✅ SSH connection to %s successful", sshTarget)
	})

	// Test Docker availability using dynamic discovery
	t.Run("DockerAvailability", func(t *testing.T) {
		sshTarget := fmt.Sprintf("%s@%s", nasUser, nasHost)

		// Find Docker binary dynamically
		findCmd := exec.Command("ssh", sshTarget, "find /share -name docker -type f 2>/dev/null | grep container-station | head -1")
		findOutput, err := findCmd.CombinedOutput()
		if err != nil || strings.TrimSpace(string(findOutput)) == "" {
			t.Fatalf("Container Station docker binary not found: %v\nOutput: %s", err, findOutput)
		}

		dockerPath := strings.TrimSpace(string(findOutput))
		t.Logf("Found Docker binary at: %s", dockerPath)

		// Test Docker version
		cmd := exec.Command("ssh", sshTarget, fmt.Sprintf("%s version --format 'Server: {{.Server.Version}}'", dockerPath))
		output, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatalf("Docker version command failed: %v\nOutput: %s", err, output)
		}

		if !strings.Contains(string(output), "Server:") {
			t.Fatalf("Invalid Docker version output: %s", output)
		}

		t.Logf("✅ Docker accessible via SSH: %s", strings.TrimSpace(string(output)))
	})

	// Test Docker container listing
	t.Run("DockerContainerList", func(t *testing.T) {
		sshTarget := fmt.Sprintf("%s@%s", nasUser, nasHost)

		// Find Docker binary dynamically
		findCmd := exec.Command("ssh", sshTarget, "find /share -name docker -type f 2>/dev/null | grep container-station | head -1")
		findOutput, _ := findCmd.CombinedOutput()
		dockerPath := strings.TrimSpace(string(findOutput))

		cmd := exec.Command("ssh", sshTarget, fmt.Sprintf("%s ps --format 'table {{.Names}}\t{{.Image}}\t{{.Status}}'", dockerPath))
		output, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatalf("Docker ps command failed: %v\nOutput: %s", err, output)
		}

		lines := strings.Split(strings.TrimSpace(string(output)), "\n")
		if len(lines) < 1 {
			t.Fatal("No output from docker ps command")
		}

		// First line should be header
		if !strings.Contains(lines[0], "NAMES") {
			t.Fatalf("Invalid docker ps output format: %s", lines[0])
		}

		t.Logf("✅ Docker container listing successful")
		t.Logf("Current containers on %s:", nasHost)
		for _, line := range lines {
			t.Logf("  %s", line)
		}
	})

	// Test Container Station status by checking Docker binary
	t.Run("ContainerStationStatus", func(t *testing.T) {
		sshTarget := fmt.Sprintf("%s@%s", nasUser, nasHost)

		// Use dynamic discovery to test Container Station
		cmd := exec.Command("ssh", sshTarget, "find /share -name docker -type f 2>/dev/null | grep container-station >/dev/null && echo 'Container Station active' || echo 'Container Station not found'")
		output, err := cmd.CombinedOutput()

		status := strings.TrimSpace(string(output))
		if err != nil || !strings.Contains(status, "Container Station active") {
			t.Logf("Warning: Container Station status: %s (error: %v)", status, err)
		} else {
			t.Logf("✅ Container Station is active")
		}
	})

	// Test volume path access (use first available volume)
	t.Run("VolumePathAccess", func(t *testing.T) {
		sshTarget := fmt.Sprintf("%s@%s", nasUser, nasHost)

		// Find first available volume
		volumeCmd := exec.Command("ssh", sshTarget, "ls -d /share/CACHEDEV*_DATA /share/ZFS*_DATA 2>/dev/null | head -1")
		volumeOutput, err := volumeCmd.CombinedOutput()
		if err != nil {
			t.Fatalf("Failed to find any volumes: %v", err)
		}

		firstVolume := strings.TrimSpace(string(volumeOutput))
		if firstVolume == "" {
			t.Skip("No volumes found for testing")
		}

		testPath := fmt.Sprintf("%s/qnap-docker-test", firstVolume)

		// Create test directory
		cmd := exec.Command("ssh", sshTarget, fmt.Sprintf("mkdir -p %s && echo 'Directory created successfully'", testPath))
		output, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatalf("Failed to create test directory: %v\nOutput: %s", err, output)
		}

		if !strings.Contains(string(output), "Directory created successfully") {
			t.Fatalf("Unexpected output from directory creation: %s", output)
		}

		// Test write permission
		testFile := filepath.Join(testPath, "test.txt")
		cmd = exec.Command("ssh", sshTarget, fmt.Sprintf("echo 'test content' > %s && cat %s", testFile, testFile))
		output, err = cmd.CombinedOutput()
		if err != nil {
			t.Fatalf("Failed to write test file: %v\nOutput: %s", err, output)
		}

		if !strings.Contains(string(output), "test content") {
			t.Fatalf("File write/read test failed: %s", output)
		}

		// Cleanup test file
		cmd = exec.Command("ssh", sshTarget, fmt.Sprintf("rm -f %s", testFile))
		cmd.Run()

		t.Logf("✅ Volume path %s accessible and writable", testPath)
	})

	// Test CACHEDEV and ZFS volume detection
	t.Run("VolumeDetection", func(t *testing.T) {
		sshTarget := fmt.Sprintf("%s@%s", nasUser, nasHost)
		cmd := exec.Command("ssh", sshTarget, "ls -d /share/CACHEDEV*_DATA /share/ZFS*_DATA 2>/dev/null || echo 'No volumes found'")
		output, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatalf("Failed to detect volumes: %v\nOutput: %s", err, output)
		}

		result := strings.TrimSpace(string(output))
		if strings.Contains(result, "No volumes found") {
			t.Fatalf("No CACHEDEV or ZFS volumes detected on QNAP NAS")
		}

		volumes := strings.Fields(result)
		t.Logf("✅ Detected volumes: %v", volumes)

		// Ensure at least one volume exists
		if len(volumes) == 0 {
			t.Errorf("No volumes found")
		}
	})
}
