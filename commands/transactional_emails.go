package commands

import "github.com/spf13/cobra"

var transactionalEmailsCmd = &cobra.Command{
	Use:   "transactional-emails",
	Short: "",
}

func init() {
	rootCmd.AddCommand(transactionalEmailsCmd)
}
