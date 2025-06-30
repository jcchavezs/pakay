package main

import (
	"fmt"
	"os"

	"github.com/jcchavezs/pakay"
	"github.com/spf13/cobra"
)

var config = `---
- name: your_email
  description: The e-mail of your account
  sources:
  - type: env
    env: 
      key: YOUR_EMAIL
  - type: stdin
    stdin: 
      prompt: Please insert your e-mail
`

var rootCmd = &cobra.Command{
	Use:  "example",
	Args: cobra.NoArgs,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if err := pakay.LoadSecretsConfig([]byte(config)); err != nil {
			return fmt.Errorf("loading secrets: %w", err)
		}

		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		email, found := pakay.GetSecret(cmd.Context(), "your_email")
		if found {
			_, _ = fmt.Fprintf(cmd.OutOrStdout(), "Your e-mail is: %s\n", email)
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
