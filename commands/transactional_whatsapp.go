package commands

import "github.com/spf13/cobra"

var transactionalWhatsappCmd = &cobra.Command{
	Use:   "transactional-whatsapp",
	Short: "",
}

func init() {
	rootCmd.AddCommand(transactionalWhatsappCmd)
}
