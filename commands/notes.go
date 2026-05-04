package commands

import "github.com/spf13/cobra"

var notesCmd = &cobra.Command{
	Use:   "notes",
	Short: "",
}

func init() {
	rootCmd.AddCommand(notesCmd)
}
