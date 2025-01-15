package commands

import (
	"github.com/spf13/cobra"
)

func TaskCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "task",
		Short: "Manage tasks in the YOLO project",
		Long: `Create, list, and manage tasks in your YOLO project.
A task is a small, well-defined unit of work that belongs to a feature.`,
		Run: func(cmd *cobra.Command, args []string) {
			// TODO: Implement task management
			cmd.Help()
		},
	}
	return cmd
} 