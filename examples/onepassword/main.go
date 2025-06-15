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
  - type: onepassword_cli
    onepassword_cli:
      ref: op://{{ $.op_vault }}/my_test_credential/username
`

func init() {
	rootCmd.PersistentFlags().String("op-vault", "Personal", "The vault for using onepassword CLI")
}

var rootCmd = &cobra.Command{
	Use:  "example",
	Args: cobra.NoArgs,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		opVault, err := cmd.Flags().GetString("op-vault")
		if err != nil {
			return fmt.Errorf("getting vault flag: %w", err)
		}

		if err := pakay.LoadSecretsFromBytes([]byte(config), pakay.LoadOptions{
			Variables: map[string]string{
				"op_vault": opVault,
			},
		}); err != nil {
			return fmt.Errorf("loading secrets: %w", err)
		}

		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		val, found := pakay.GetSecret(cmd.Context(), "my_test_credential")
		if found {
			fmt.Fprintf(cmd.ErrOrStderr(), "Credential found: %s\n", val)
		} else {
			fmt.Fprintln(cmd.ErrOrStderr(), "Credential not found")
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
