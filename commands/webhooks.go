package commands

import "github.com/spf13/cobra"

var webhooksCmd = &cobra.Command{
	Use:   "webhooks",
	Short: "",
}

func init() {
	rootCmd.AddCommand(webhooksCmd)
}
