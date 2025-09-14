package qnap

import (
	"fmt"
	"net"
	"os"

	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/agent"
)

// hasSSHAgent checks if ssh-agent is available and has keys
func hasSSHAgent() bool {
	agentSock := os.Getenv("SSH_AUTH_SOCK")
	if agentSock == "" {
		return false
	}

	// Try to connect to the agent
	conn, err := net.Dial("unix", agentSock)
	if err != nil {
		return false
	}
	defer conn.Close()

	agentClient := agent.NewClient(conn)
	keys, err := agentClient.List()
	if err != nil {
		return false
	}

	return len(keys) > 0
}

// connectSSHWithAgent connects to SSH using ssh-agent
func (c *Connection) connectSSHWithAgent() error {
	agentSock := os.Getenv("SSH_AUTH_SOCK")
	if agentSock == "" {
		return fmt.Errorf("SSH_AUTH_SOCK not set")
	}

	conn, err := net.Dial("unix", agentSock)
	if err != nil {
		return fmt.Errorf("failed to connect to ssh-agent: %w", err)
	}
	defer conn.Close()

	agentClient := agent.NewClient(conn)

	// Configure SSH client with agent
	sshConfig := &ssh.ClientConfig{
		User: c.config.User,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeysCallback(agentClient.Signers),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(), // TODO: Implement proper host key verification
	}

	// Connect to SSH server
	address := fmt.Sprintf("%s:%d", c.config.Host, c.config.Port)
	client, err := ssh.Dial("tcp", address, sshConfig)
	if err != nil {
		return fmt.Errorf("failed to dial SSH with agent: %w", err)
	}

	c.sshClient = client
	return nil
}
