package commands

import "github.com/spf13/cobra"

var processCmd = &cobra.Command{
	Use:   "process",
	Short: "",
}

func init() {
	rootCmd.AddCommand(processCmd)
}
