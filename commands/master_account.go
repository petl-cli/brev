package commands

import "github.com/spf13/cobra"

var masterAccountCmd = &cobra.Command{
	Use:   "master-account",
	Short: "",
}

func init() {
	rootCmd.AddCommand(masterAccountCmd)
}
