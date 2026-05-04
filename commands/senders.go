package commands

import "github.com/spf13/cobra"

var sendersCmd = &cobra.Command{
	Use:   "senders",
	Short: "",
}

func init() {
	rootCmd.AddCommand(sendersCmd)
}
