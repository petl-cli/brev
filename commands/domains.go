package commands

import "github.com/spf13/cobra"

var domainsCmd = &cobra.Command{
	Use:   "domains",
	Short: "",
}

func init() {
	rootCmd.AddCommand(domainsCmd)
}
