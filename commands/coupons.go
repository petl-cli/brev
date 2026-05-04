package commands

import "github.com/spf13/cobra"

var couponsCmd = &cobra.Command{
	Use:   "coupons",
	Short: "",
}

func init() {
	rootCmd.AddCommand(couponsCmd)
}
