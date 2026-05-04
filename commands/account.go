package commands

import "github.com/spf13/cobra"

var accountCmd = &cobra.Command{
	Use:   "account",
	Short: "",
}

func init() {
	rootCmd.AddCommand(accountCmd)
}
