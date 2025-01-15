package commands

import (
	"github.com/spf13/cobra"
)

func EpicCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "epic",
		Short: "Manage epics in the YOLO project",
		Long: `Create, list, and manage epics in your YOLO project.
An epic is a large body of work that can be broken down into features.`,
		Run: func(cmd *cobra.Command, args []string) {
			// TODO: Implement epic management
			cmd.Help()
		},
	}
	return cmd
} 