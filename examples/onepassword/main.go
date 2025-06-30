package main

import (
	"fmt"
	"os"

	"github.com/jcchavezs/pakay"
	"github.com/spf13/cobra"
)

var config = `---
- name: my_test_credential
  description: Your account
  sources:
  - type: 1password
    1password:
      ref: op://{{ $.op_vault }}/my_test_credential/username
`

var opVault string

func init() {
	rootCmd.PersistentFlags().StringVar(&opVault, "op-vault", "Personal", "The vault for using 1Password CLI")
}

var rootCmd = &cobra.Command{
	Use:  "example",
	Args: cobra.NoArgs,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		if err := pakay.LoadSecretsConfigWithOptions([]byte(config), pakay.LoadOptions{
			Variables: map[string]string{
				"op_vault": opVault,
			},
		}); err != nil {
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
