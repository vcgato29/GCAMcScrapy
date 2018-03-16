package main

import (
	"github.com/spf13/cobra"
)

var (
	rootCmd = &cobra.Command{
		Use:   "mcscrapy",
		Short: "GCA Wordpress website scraper.",
		Run:   run,
	}
)

func main() {
	rootCmd.AddCommand(
		scrapeCmd,
		previewCmd,
	)
	rootCmd.Execute()
}

func run(cmd *cobra.Command, args []string) {
	cmd.Help()
}
