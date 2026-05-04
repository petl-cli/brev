package commands

import "github.com/spf13/cobra"

var paymentsCmd = &cobra.Command{
	Use:   "payments",
	Short: "",
}

func init() {
	rootCmd.AddCommand(paymentsCmd)
}
