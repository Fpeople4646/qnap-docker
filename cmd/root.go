package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	// Version is the current version of qnap-docker
	Version string
	// Commit is the git commit hash
	Commit string
	// Date is the build date
	Date string
)

var rootCmd = &cobra.Command{
	Use:   "qnap-docker",
	Short: "Deploy containers to QNAP Container Station",
	Long: `qnap-docker is a CLI tool that simplifies Docker container deployment
to QNAP NAS devices with Container Station. It handles SSH connection management,
Docker client setup, and path resolution issues specific to QNAP Container Station.`,
	Version: getVersion(),
}

func getVersion() string {
	if Version == "" {
		return "dev"
	}
	return fmt.Sprintf("%s (commit %s, built %s)", Version, Commit, Date)
}

// Execute runs the root command
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.AddCommand(initCmd)
	rootCmd.AddCommand(runCmd)
	rootCmd.AddCommand(deployCmd)
	rootCmd.AddCommand(psCmd)
	rootCmd.AddCommand(rmCmd)
}
