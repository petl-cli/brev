package commands

import "github.com/spf13/cobra"

var conversationsCmd = &cobra.Command{
	Use:   "conversations",
	Short: "",
}

func init() {
	rootCmd.AddCommand(conversationsCmd)
}
