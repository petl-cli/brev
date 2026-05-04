package commands

import "github.com/spf13/cobra"

var resellerCmd = &cobra.Command{
	Use:   "reseller",
	Short: "",
}

func init() {
	rootCmd.AddCommand(resellerCmd)
}
