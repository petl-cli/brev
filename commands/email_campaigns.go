package commands

import "github.com/spf13/cobra"

var emailCampaignsCmd = &cobra.Command{
	Use:   "email-campaigns",
	Short: "",
}

func init() {
	rootCmd.AddCommand(emailCampaignsCmd)
}
