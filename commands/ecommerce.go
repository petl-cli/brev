package commands

import "github.com/spf13/cobra"

var ecommerceCmd = &cobra.Command{
	Use:   "ecommerce",
	Short: "",
}

func init() {
	rootCmd.AddCommand(ecommerceCmd)
}
