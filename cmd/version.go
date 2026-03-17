package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

var (
	Version = "dev"
	Commit  = "none"
	Date    = "unknown"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print build version information",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Fprintf(cmd.OutOrStdout(), "version: %s\n", Version)
		fmt.Fprintf(cmd.OutOrStdout(), "commit: %s\n", Commit)
		fmt.Fprintf(cmd.OutOrStdout(), "date: %s\n", Date)
	},
}

func versionString() string {
	parts := []string{"0ximg", Version}
	if Commit != "" && Commit != "none" {
		parts = append(parts, fmt.Sprintf("(%s)", Commit))
	}

	return strings.Join(parts, " ")
}
