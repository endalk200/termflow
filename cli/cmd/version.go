package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the termflow-cli version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("termflow-cli 0.0.1")
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
