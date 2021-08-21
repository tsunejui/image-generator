package cmd

import (
	"github.com/spf13/cobra"
)

var (
	rootCmd = &cobra.Command{
		Use:   "generator",
		Short: "A CLI tool for merging images",
		Long:  `This tool provides an easy and extensible way to merge images.`,
		CompletionOptions: cobra.CompletionOptions{
			DisableDefaultCmd: true,
		},
	}
)

func init() {
	rootCmd.AddCommand(MergeCommand)
}

func Execute() error {
	return rootCmd.Execute()
}
