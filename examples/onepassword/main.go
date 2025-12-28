package main

import (
	"fmt"
	"os"

	"github.com/jcchavezs/pakay"
	"github.com/spf13/cobra"
)

func createConfig(opVault string) pakay.SecretsConfig {
	return pakay.SecretsConfig{
		{
			Name:        "my_test_credential",
			Description: "Your account",
			Sources: []pakay.SecretSource{
				{
					TypedConfig: &pakay.OnePasswordConfig{
						Ref: fmt.Sprintf("op://%s/my_test_credential/username", opVault),
					},
				},
			},
		},
	}

}

var opVault string

func init() {
	rootCmd.PersistentFlags().StringVar(&opVault, "op-vault", "Personal", "The vault for using 1Password CLI")
}

var rootCmd = &cobra.Command{
	Use:  "example",
	Args: cobra.NoArgs,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		if err := pakay.LoadSecrets(createConfig(opVault)); err != nil {
			return fmt.Errorf("loading secrets config: %w", err)
		}

		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		cmd.Printf("Reading your credential from \"op://%s/my_test_credential/username\"...\n", opVault)

		val, found := pakay.GetSecret(cmd.Context(), "my_test_credential")
		if found {
			cmd.PrintErrf("✅ Credential found: %s\n", val)
		} else {
			cmd.PrintErrln("🚫 Credential not found")
		}

		return nil
	},
}

func main() {
	err := rootCmd.Execute()
	if err != nil {
		fmt.Printf("ERROR: %v\n", err)
		os.Exit(1)
	}
}
