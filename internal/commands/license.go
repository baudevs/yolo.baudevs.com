package commands

import (
	"fmt"
	"os"

	"github.com/baudevs/yolo.baudevs.com/internal/license"
	"github.com/spf13/cobra"
)

func NewLicenseCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "license",
		Short: "Manage your YOLO license",
		Long: `Manage your YOLO license and credits.
You can activate a new license, check your status, or purchase more credits.`,
	}

	cmd.AddCommand(
		newLicenseActivateCommand(),
		newLicenseStatusCommand(),
	)

	return cmd
}

func newLicenseActivateCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "activate",
		Short: "Activate a YOLO license",
		Long: `Activate a new YOLO license or purchase more credits.
This will open your browser to complete the purchase.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			// Initialize license manager
			manager, err := license.NewManager(license.Config{
				StripeSecretKey:  os.Getenv("STRIPE_SECRET_KEY"),
				DefaultOpenAIKey: os.Getenv("OPENAI_API_KEY"),
			})
			if err != nil {
				return fmt.Errorf("failed to initialize license manager: %w", err)
			}

			// Create checkout session
			starterPack := license.GetPackageByType(license.PlanStarter)
			if starterPack == nil {
				return fmt.Errorf("failed to get starter package")
			}

			url, err := manager.CreateCheckoutSession(*starterPack)
			if err != nil {
				return fmt.Errorf("failed to create checkout session: %w", err)
			}

			fmt.Println("🎉 Opening checkout in your browser...")
			fmt.Printf("\nIf it doesn't open automatically, visit:\n%s\n", url)

			return nil
		},
	}
}

func newLicenseStatusCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "status",
		Short: "Check license status",
		Long:  "Check your YOLO license status and remaining credits.",
		RunE: func(cmd *cobra.Command, args []string) error {
			// Initialize license manager
			manager, err := license.NewManager(license.Config{
				StripeSecretKey:  os.Getenv("STRIPE_SECRET_KEY"),
				DefaultOpenAIKey: os.Getenv("OPENAI_API_KEY"),
			})
			if err != nil {
				return fmt.Errorf("failed to initialize license manager: %w", err)
			}

			// Get current license
			lic := manager.GetLicense()

			fmt.Println("📝 License Status")
			fmt.Println("--------------")

			if lic == nil || !lic.IsActive {
				fmt.Println("❌ No active license")
				fmt.Println("\nTo get started:")
				fmt.Println("1. Purchase credits:")
				fmt.Println("   yolo license activate")
				fmt.Println("\n2. Or use your own OpenAI key:")
				fmt.Println("   yolo ai configure --api-key YOUR_API_KEY")
				return nil
			}

			// Show license details
			fmt.Println("✅ License active")
			fmt.Printf("Plan: %s\n", lic.PlanType)
			if !lic.ExpiresAt.IsZero() {
				fmt.Printf("Expires: %s\n", lic.ExpiresAt.Format("2006-01-02"))
			}
			if lic.PlanType == license.PlanUnlimited {
				fmt.Println("Credits: ♾️  Unlimited")
			} else {
				fmt.Printf("Credits: %d remaining\n", lic.CreditsLeft)
			}

			return nil
		},
	}
}
