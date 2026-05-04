package commands

import "github.com/spf13/cobra"

var filesCmd = &cobra.Command{
	Use:   "files",
	Short: "",
}

func init() {
	rootCmd.AddCommand(filesCmd)
}
