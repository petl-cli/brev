package commands

import "github.com/spf13/cobra"

var eventCmd = &cobra.Command{
	Use:   "event",
	Short: "",
}

func init() {
	rootCmd.AddCommand(eventCmd)
}
