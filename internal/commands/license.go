package commands

import (
	"fmt"

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
		newLicenseCheckoutCommand(),
	)

	return cmd
}

func newLicenseActivateCommand() *cobra.Command {
	var key string

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

			// Check if key is provided
			if key == "" {
				return fmt.Errorf("license key is required")
			}

			// Activate license with backend
			if err := manager.ActivateLicense(key); err != nil {
				return fmt.Errorf("failed to activate license: %w", err)
			}

			fmt.Println("✨ License activated successfully!")
			return nil
		},
	}

	cmd.Flags().StringVarP(&key, "key", "k", "", "License key")
	cmd.MarkFlagRequired("key")

	return cmd
}

func newLicenseStatusCommand() *cobra.Command {
	cmd := &cobra.Command{
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
				fmt.Println("❌ No license found")
				fmt.Println("\nTo get started, activate your license:")
				fmt.Println("  yolo license activate --key YOUR_LICENSE_KEY")
				fmt.Println("\nOr purchase a new license:")
				fmt.Println("  yolo license checkout --email YOUR_EMAIL")
				return nil
			}

			if !lic.IsActive {
				fmt.Println("❌ License is inactive")
				return nil
			}

			fmt.Println("✅ License is active")
			fmt.Printf("Credits remaining: %d\n", lic.Credits)
			fmt.Printf("Last modified: %s\n", lic.LastModified.Format("2006-01-02 15:04:05"))

			return nil
		},
	}

	return cmd
}

func newLicenseCheckoutCommand() *cobra.Command {
	var (
		email string
		pkg   string
	)

	cmd := &cobra.Command{
		Use:   "checkout",
		Short: "Purchase a new license",
		Long: `Purchase a new YOLO license. Available packages:
- basic: 100 credits ($10)
- pro: 500 credits ($40)
- enterprise: 2000 credits ($150)`,
		RunE: func(cmd *cobra.Command, args []string) error {
			// Create license manager
			manager, err := license.NewManager()
			if err != nil {
				return fmt.Errorf("failed to create license manager: %w", err)
			}

			// Create checkout session
			sessionID, err := manager.CreateCheckoutSession(email, pkg)
			if err != nil {
				return fmt.Errorf("failed to create checkout session: %w", err)
			}

			fmt.Println("✨ Checkout session created!")
			fmt.Printf("\nPlease complete your purchase at: %s\n", sessionID)
			fmt.Println("\nAfter purchase, activate your license with:")
			fmt.Println("  yolo license activate --key YOUR_LICENSE_KEY")

			return nil
		},
	}

	cmd.Flags().StringVarP(&email, "email", "e", "", "Email address for license")
	cmd.Flags().StringVarP(&pkg, "package", "p", "basic", "License package (basic, pro, enterprise)")
	cmd.MarkFlagRequired("email")

	return cmd
}
