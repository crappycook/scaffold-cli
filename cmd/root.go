package cmd

import (
	"github.com/crappycook/scaffold-cli/internal/command/new"
	"github.com/spf13/cobra"
)

var CmdRoot = &cobra.Command{
	Use:     "scaffold-cli",
	Example: "scaffold-cli new project",
	Short:   "Build new project from your layout",
}

func init() {
	CmdRoot.AddCommand(new.CmdNew)
}

// Execute executes the root command.
func Execute() error {
	return CmdRoot.Execute()
}
