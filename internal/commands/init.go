package commands

import (
	"github.com/baudevs/yolo-cli/internal/core"
	"github.com/spf13/cobra"
)

func InitCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "init",
		Short: "Initialize a new project with YOLO methodology",
		Long: `Creates the YOLO directory structure and initial files required 
for following the YOLO methodology in your project.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return core.InitializeProject()
		},
	}

	return cmd
} 