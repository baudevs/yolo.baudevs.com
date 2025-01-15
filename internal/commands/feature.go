package commands

import (
	"github.com/spf13/cobra"
)

func FeatureCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "feature",
		Short: "Manage features in the YOLO project",
		Long: `Create, list, and manage features in your YOLO project.
A feature is a collection of related tasks that deliver value to users.`,
		Run: func(cmd *cobra.Command, args []string) {
			// TODO: Implement feature management
			cmd.Help()
		},
	}
	return cmd
} 