package commands

import "github.com/spf13/cobra"

var contactsCmd = &cobra.Command{
	Use:   "contacts",
	Short: "",
}

func init() {
	rootCmd.AddCommand(contactsCmd)
}
