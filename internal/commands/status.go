package commands

import (
	"github.com/spf13/cobra"
)

func StatusCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "status",
		Short: "Show the status of the YOLO project",
		Long: `Display the current status of your YOLO project,
including active epics, features, and tasks.`,
		Run: func(cmd *cobra.Command, args []string) {
			// TODO: Implement status display
			cmd.Help()
		},
	}
	return cmd
} 