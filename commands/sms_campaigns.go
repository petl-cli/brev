package commands

import "github.com/spf13/cobra"

var smsCampaignsCmd = &cobra.Command{
	Use:   "sms-campaigns",
	Short: "",
}

func init() {
	rootCmd.AddCommand(smsCampaignsCmd)
}
