package qnap

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/docker/docker/client"
	"github.com/pkg/errors"
	"golang.org/x/crypto/ssh"

	"github.com/scttfrdmn/qnap-docker/pkg/config"
)

// Connection represents a connection to a QNAP NAS
type Connection struct {
	config     *config.Config
	sshClient  *ssh.Client
	dockerAPI  *client.Client
	dockerPath string // cached Docker binary path
}

// NewConnection creates a new connection with the given configuration
func NewConnection(cfg *config.Config) *Connection {
	return &Connection{
		config: cfg,
	}
}

// Connect establishes the SSH connection to the QNAP NAS
func (c *Connection) Connect() error {
	if err := c.connectSSH(); err != nil {
		return errors.Wrap(err, "failed to establish SSH connection")
	}

	// For v0.1.0, we use SSH commands directly instead of Docker client
	// This avoids the complex SSH tunneling setup and works reliably
	return nil
}

func (c *Connection) connectSSH() error {
	// Try ssh-agent first if available
	if hasSSHAgent() {
		if err := c.connectSSHWithAgent(); err == nil {
			return nil // Success with agent
		}
		// If agent fails, fall back to key file
		fmt.Printf("Warning: ssh-agent authentication failed, trying key file...\n")
	}

	// Fallback to key file authentication
	return c.connectSSHWithKeyFile()
}

func (c *Connection) connectSSHWithKeyFile() error {
	// Read SSH private key
	keyBytes, err := os.ReadFile(c.config.SSHKeyPath)
	if err != nil {
		return fmt.Errorf("failed to read SSH key: %w", err)
	}

	// Parse private key
	signer, err := ssh.ParsePrivateKey(keyBytes)
	if err != nil {
		return fmt.Errorf("failed to parse SSH key: %w", err)
	}

	// Configure SSH client
	sshConfig := &ssh.ClientConfig{
		User: c.config.User,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(), // TODO: Implement proper host key verification
	}

	// Connect to SSH server
	address := fmt.Sprintf("%s:%d", c.config.Host, c.config.Port)
	client, err := ssh.Dial("tcp", address, sshConfig)
	if err != nil {
		return fmt.Errorf("failed to dial SSH: %w", err)
	}

	c.sshClient = client
	return nil
}

// Close closes the connection and cleans up resources
func (c *Connection) Close() error {
	var errs []string

	if c.dockerAPI != nil {
		if err := c.dockerAPI.Close(); err != nil {
			errs = append(errs, fmt.Sprintf("Docker client: %v", err))
		}
	}

	if c.sshClient != nil {
		if err := c.sshClient.Close(); err != nil {
			errs = append(errs, fmt.Sprintf("SSH client: %v", err))
		}
	}

	if len(errs) > 0 {
		return fmt.Errorf("errors closing connections: %s", strings.Join(errs, ", "))
	}

	return nil
}

// ExecuteCommand executes a command over SSH and returns the output
func (c *Connection) ExecuteCommand(cmd string) (string, error) {
	if c.sshClient == nil {
		return "", fmt.Errorf("SSH client not connected")
	}

	session, err := c.sshClient.NewSession()
	if err != nil {
		return "", fmt.Errorf("failed to create SSH session: %w", err)
	}
	defer session.Close()

	// Set up environment to include Container Station binary paths
	containerStationPath := fmt.Sprintf("%s:%s", ContainerStationBinPath, ContainerStationSbinPath)
	session.Setenv("PATH", containerStationPath+":/usr/local/bin:/usr/bin:/bin")

	output, err := session.CombinedOutput(cmd)
	if err != nil {
		return string(output), fmt.Errorf("command failed: %w", err)
	}

	return string(output), nil
}

// ExecuteDockerCommand executes a Docker command with full path
func (c *Connection) ExecuteDockerCommand(args []string) (string, error) {
	// Use cached Docker path or find it
	if c.dockerPath == "" {
		dockerPath, err := c.findDockerBinary()
		if err != nil {
			return "", fmt.Errorf("docker binary not found: %w", err)
		}
		c.dockerPath = dockerPath // cache for future use
	}

	// Execute command with found binary path
	cmd := c.dockerPath + " " + strings.Join(args, " ")
	return c.ExecuteCommand(cmd)
}

// GetDockerClient returns the Docker client (nil in v0.1.x)
func (c *Connection) GetDockerClient() *client.Client {
	// In v0.1.0, we don't use the Docker client API, just SSH commands
	return nil
}

// TestConnection tests SSH and Docker connectivity
func (c *Connection) TestConnection() error {
	// Test SSH connection
	if _, err := c.ExecuteCommand("echo 'SSH connection test'"); err != nil {
		return fmt.Errorf("ssh connection test failed: %w", err)
	}

	// Test Container Station installation
	if err := c.testContainerStation(); err != nil {
		return fmt.Errorf("container station test failed: %w", err)
	}

	// Test Docker command
	if _, err := c.ExecuteDockerCommand([]string{"version", "--format", "'{{.Server.Version}}'"}); err != nil {
		return fmt.Errorf("docker connection test failed: %w", err)
	}

	return nil
}

// testContainerStation verifies that Container Station is installed and accessible
func (c *Connection) testContainerStation() error {
	// Try to find Docker binary if default path doesn't work
	dockerPath, err := c.findDockerBinary()
	if err != nil {
		return fmt.Errorf("container Station not found: %w", err)
	}

	// Check if binary is executable
	if _, err := c.ExecuteCommand(fmt.Sprintf("test -x %s", dockerPath)); err != nil {
		return fmt.Errorf("docker binary at %s is not executable", dockerPath)
	}

	return nil
}

// findDockerBinary attempts to locate the Docker binary on the QNAP system
func (c *Connection) findDockerBinary() (string, error) {
	// Try default path first
	if _, err := c.ExecuteCommand(fmt.Sprintf("test -f %s", DockerBinary)); err == nil {
		return DockerBinary, nil
	}

	// Search for Container Station Docker binary with better error handling
	output, err := c.ExecuteCommand("find /share -name docker -type f 2>/dev/null | grep container-station | head -1")
	if err != nil {
		return "", fmt.Errorf("failed to search for docker binary: %w", err)
	}

	dockerPath := strings.TrimSpace(output)
	if dockerPath == "" {
		// Try more specific paths if general search fails
		commonPaths := []string{
			"/share/CACHEDEV1_DATA/.qpkg/container-station/bin/docker",
			"/share/ZFS530_DATA/.qpkg/container-station/usr/bin/.libs/docker",
			"/share/CACHEDEV2_DATA/.qpkg/container-station/bin/docker",
			"/share/ZFS1_DATA/.qpkg/container-station/usr/bin/.libs/docker",
		}

		for _, path := range commonPaths {
			if _, err := c.ExecuteCommand(fmt.Sprintf("test -f %s", path)); err == nil {
				return path, nil
			}
		}

		return "", fmt.Errorf("no Container Station docker binary found")
	}

	return dockerPath, nil
}

// StreamCommand executes a command and streams output to writers
func (c *Connection) StreamCommand(cmd string, stdout, stderr io.Writer) error {
	if c.sshClient == nil {
		return fmt.Errorf("SSH client not connected")
	}

	session, err := c.sshClient.NewSession()
	if err != nil {
		return fmt.Errorf("failed to create SSH session: %w", err)
	}
	defer session.Close()

	// Set up environment to include Container Station binary paths
	containerStationPath := fmt.Sprintf("%s:%s", ContainerStationBinPath, ContainerStationSbinPath)
	session.Setenv("PATH", containerStationPath+":/usr/local/bin:/usr/bin:/bin")

	// Connect streams
	session.Stdout = stdout
	session.Stderr = stderr

	// Run command
	if err := session.Run(cmd); err != nil {
		return fmt.Errorf("command failed: %w", err)
	}

	return nil
}

// DetectAvailableVolumes detects available CACHEDEV and ZFS volumes on the QNAP NAS
func (c *Connection) DetectAvailableVolumes() ([]string, error) {
	// List all CACHEDEV and ZFS directories
	output, err := c.ExecuteCommand("ls -d /share/CACHEDEV*_DATA /share/ZFS*_DATA 2>/dev/null || true")
	if err != nil {
		return nil, fmt.Errorf("failed to detect volumes: %w", err)
	}

	if strings.TrimSpace(output) == "" {
		return nil, fmt.Errorf("no CACHEDEV or ZFS volumes found")
	}

	volumes := strings.Fields(output)
	return volumes, nil
}

// GetPrimaryVolumePath returns the primary volume path for the connection
func (c *Connection) GetPrimaryVolumePath() string {
	if c.config.Defaults.PrimaryVolume != "" {
		return fmt.Sprintf("/share/%s", c.config.Defaults.PrimaryVolume)
	}
	return "/share/CACHEDEV1_DATA"
}

// GetDockerPath returns the cached Docker binary path
func (c *Connection) GetDockerPath() (string, error) {
	if c.dockerPath == "" {
		dockerPath, err := c.findDockerBinary()
		if err != nil {
			return "", err
		}
		c.dockerPath = dockerPath
	}
	return c.dockerPath, nil
}
