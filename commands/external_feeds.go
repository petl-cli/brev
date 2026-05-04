package commands

import "github.com/spf13/cobra"

var externalFeedsCmd = &cobra.Command{
	Use:   "external-feeds",
	Short: "",
}

func init() {
	rootCmd.AddCommand(externalFeedsCmd)
}
