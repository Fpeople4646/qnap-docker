package integration

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// TestConnectionToAstrapi tests direct connection to your local QNAP system
func TestConnectionToAstrapi(t *testing.T) {
	// Skip in CI environments (GitHub Actions can't reach local networks)
	if os.Getenv("CI") != "" || os.Getenv("GITHUB_ACTIONS") != "" {
		t.Skip("Skipping local network test in CI environment")
	}
	// Skip in short mode (unit tests only)
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}
	// Test hostname resolution
	t.Run("HostnameResolution", func(t *testing.T) {
		cmd := exec.Command("ping", "-c", "1", "astrapi.local")
		output, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatalf("Failed to ping astrapi.local: %v\nOutput: %s", err, output)
		}
		t.Logf("✅ Hostname astrapi.local resolves successfully")
	})

	// Test SSH connectivity
	t.Run("SSHConnectivity", func(t *testing.T) {
		cmd := exec.Command("ssh", "-o", "ConnectTimeout=10", "-o", "BatchMode=yes",
			"scttfrdmn@astrapi.local", "echo 'SSH connection successful'; whoami")
		output, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatalf("SSH connection failed: %v\nOutput: %s", err, output)
		}

		if !strings.Contains(string(output), "SSH connection successful") {
			t.Fatalf("SSH connection test failed. Output: %s", output)
		}

		if !strings.Contains(string(output), "scttfrdmn") {
			t.Fatalf("Wrong user returned. Output: %s", output)
		}

		t.Logf("✅ SSH connection to scttfrdmn@astrapi.local successful")
	})

	// Test Docker availability
	t.Run("DockerAvailability", func(t *testing.T) {
		cmd := exec.Command("ssh", "scttfrdmn@astrapi.local", "/share/CACHEDEV1_DATA/.qpkg/container-station/bin/docker version --format 'Server: {{.Server.Version}}'")
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
		cmd := exec.Command("ssh", "scttfrdmn@astrapi.local", "/share/CACHEDEV1_DATA/.qpkg/container-station/bin/docker ps --format 'table {{.Names}}\t{{.Image}}\t{{.Status}}'")
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
		t.Logf("Current containers on astrapi.local:")
		for _, line := range lines {
			t.Logf("  %s", line)
		}
	})

	// Test Container Station status by checking Docker binary
	t.Run("ContainerStationStatus", func(t *testing.T) {
		cmd := exec.Command("ssh", "scttfrdmn@astrapi.local", "test -x /share/CACHEDEV1_DATA/.qpkg/container-station/bin/docker && echo 'Container Station active' || echo 'Container Station not found'")
		output, err := cmd.CombinedOutput()

		status := strings.TrimSpace(string(output))
		if err != nil || !strings.Contains(status, "Container Station active") {
			t.Logf("Warning: Container Station status: %s (error: %v)", status, err)
			// Don't fail - might be different CACHEDEV or path
		} else {
			t.Logf("✅ Container Station is active")
		}
	})

	// Test CACHEDEV volume path access
	t.Run("CacheDevVolumePathAccess", func(t *testing.T) {
		testPath := "/share/CACHEDEV1_DATA/qnap-docker-test"

		// Create test directory
		cmd := exec.Command("ssh", "scttfrdmn@astrapi.local", fmt.Sprintf("mkdir -p %s && echo 'Directory created successfully'", testPath))
		output, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatalf("Failed to create test directory: %v\nOutput: %s", err, output)
		}

		if !strings.Contains(string(output), "Directory created successfully") {
			t.Fatalf("Unexpected output from directory creation: %s", output)
		}

		// Test write permission
		testFile := filepath.Join(testPath, "test.txt")
		cmd = exec.Command("ssh", "scttfrdmn@astrapi.local", fmt.Sprintf("echo 'test content' > %s && cat %s", testFile, testFile))
		output, err = cmd.CombinedOutput()
		if err != nil {
			t.Fatalf("Failed to write test file: %v\nOutput: %s", err, output)
		}

		if !strings.Contains(string(output), "test content") {
			t.Fatalf("File write/read test failed: %s", output)
		}

		// Cleanup test file
		cmd = exec.Command("ssh", "scttfrdmn@astrapi.local", fmt.Sprintf("rm -f %s", testFile))
		cmd.Run()

		t.Logf("✅ CACHEDEV volume path %s accessible and writable", testPath)
	})

	// Test CACHEDEV volume detection
	t.Run("CacheDevVolumeDetection", func(t *testing.T) {
		cmd := exec.Command("ssh", "scttfrdmn@astrapi.local", "ls -d /share/CACHEDEV*_DATA 2>/dev/null || echo 'No CACHEDEV volumes found'")
		output, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatalf("Failed to detect CACHEDEV volumes: %v\nOutput: %s", err, output)
		}

		result := strings.TrimSpace(string(output))
		if strings.Contains(result, "No CACHEDEV volumes found") {
			t.Fatalf("No CACHEDEV volumes detected on QNAP NAS")
		}

		volumes := strings.Fields(result)
		t.Logf("✅ Detected CACHEDEV volumes: %v", volumes)

		// Ensure at least CACHEDEV1_DATA exists
		foundCacheDev1 := false
		for _, volume := range volumes {
			if strings.Contains(volume, "CACHEDEV1_DATA") {
				foundCacheDev1 = true
				break
			}
		}

		if !foundCacheDev1 {
			t.Errorf("CACHEDEV1_DATA not found in detected volumes: %v", volumes)
		}
	})
}
