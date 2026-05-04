package commands

import "github.com/spf13/cobra"

var tasksCmd = &cobra.Command{
	Use:   "tasks",
	Short: "",
}

func init() {
	rootCmd.AddCommand(tasksCmd)
}
