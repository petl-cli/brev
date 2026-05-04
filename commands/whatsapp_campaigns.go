package commands

import "github.com/spf13/cobra"

var whatsappCampaignsCmd = &cobra.Command{
	Use:   "whatsapp-campaigns",
	Short: "",
}

func init() {
	rootCmd.AddCommand(whatsappCampaignsCmd)
}
