package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "bot",
	Short: "What is bot is your friendly onboarding helper for Slack",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Done!")
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
