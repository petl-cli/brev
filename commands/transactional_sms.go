package commands

import "github.com/spf13/cobra"

var transactionalSmsCmd = &cobra.Command{
	Use:   "transactional-sms",
	Short: "",
}

func init() {
	rootCmd.AddCommand(transactionalSmsCmd)
}
