package commands

import "github.com/spf13/cobra"

var inboundParsingCmd = &cobra.Command{
	Use:   "inbound-parsing",
	Short: "",
}

func init() {
	rootCmd.AddCommand(inboundParsingCmd)
}
