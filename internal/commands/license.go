package commands

import (
	"fmt"
	"time"

	"github.com/baudevs/yolo.baudevs.com/internal/license"
	"github.com/spf13/cobra"
)

// NewLicenseCommand returns a new license command
func NewLicenseCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "license",
		Short: "License management",
		Long:  "Manage your YOLO license and credits",
	}

	cmd.AddCommand(
		newLicenseActivateCommand(),
		newLicenseStatusCommand(),
	)

	return cmd
}

func newLicenseActivateCommand() *cobra.Command {
	var key string
	var planType string
	var credits int64

	cmd := &cobra.Command{
		Use:   "activate",
		Short: "Activate license",
		Long:  "Activate your YOLO license with a key",
		RunE: func(cmd *cobra.Command, args []string) error {
			// Create license manager
			manager, err := license.NewManager()
			if err != nil {
				return fmt.Errorf("failed to create license manager: %w", err)
			}

			// Get key from args if not provided as flag
			if key == "" && len(args) > 0 {
				key = args[0]
			}

			// Default to credits plan with 100 credits
			if planType == "" {
				planType = "credits"
			}
			if credits == 0 {
				credits = 100
			}

			// Activate license
			if err := manager.SaveLicense(&license.License{
				IsActive:     true,
				APIKey:       key,
				PlanType:     planType,
				Credits:      credits,
				LastModified: time.Now(),
			}); err != nil {
				return fmt.Errorf("failed to activate license: %w", err)
			}

			fmt.Println(" License activated successfully!")
			return nil
		},
	}

	cmd.Flags().StringVarP(&key, "key", "k", "", "License key")
	cmd.Flags().StringVarP(&planType, "plan", "p", "", "Plan type (credits or unlimited)")
	cmd.Flags().Int64VarP(&credits, "credits", "c", 0, "Initial credits (for credits plan)")

	return cmd
}

func newLicenseStatusCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "status",
		Short: "Check license status",
		Long:  "Check your YOLO license status and remaining credits",
		RunE: func(cmd *cobra.Command, args []string) error {
			// Create license manager
			manager, err := license.NewManager()
			if err != nil {
				return fmt.Errorf("failed to create license manager: %w", err)
			}

			// Get license
			lic := manager.GetLicense()
			if lic == nil {
				fmt.Println(" No license found")
				fmt.Println("\nTo get started, activate your license:")
				fmt.Println("  yolo license activate --key YOUR_LICENSE_KEY")
				return nil
			}

			if !lic.IsActive {
				fmt.Println(" License is inactive")
				return nil
			}

			fmt.Println(" License is active")
			if lic.PlanType == "unlimited" {
				fmt.Println("Credits:  Unlimited")
			} else {
				fmt.Printf("Credits: %d remaining\n", lic.Credits)
			}

			return nil
		},
	}
}
