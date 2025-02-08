package commands

import (
	"fmt"

	"github.com/baudevs/yolo.baudevs.com/internal/version"
	"github.com/spf13/cobra"
)

// NewVersionCommand returns a new version command
func NewVersionCommand() *cobra.Command {
	var outputJSON bool

	cmd := &cobra.Command{
		Use:   "version",
		Short: "Print the version information",
		Long: `Display detailed version information about the YOLO CLI.
This includes version number, git commit hash, build date, and platform details.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			info := version.Get()
			
			if outputJSON {
				fmt.Println(info.JSON())
				return nil
			}
			
			fmt.Println(info.String())
			return nil
		},
	}

	cmd.Flags().BoolVar(&outputJSON, "json", false, "Output in JSON format")
	return cmd
}
