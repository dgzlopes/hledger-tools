package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var journalFilePath string

var rootCmd = &cobra.Command{
	Use:   "hledger-tools",
	Short: "Tools for working with hledger journals",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		if journalFilePath == "" {
			return fmt.Errorf("the --journal (-j) flag is required (path to journal file)")
		}
		return nil
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&journalFilePath, "journal", "j", "", "Path to journal file")

	rootCmd.AddCommand(AddCmd)
	rootCmd.AddCommand(ImportCmd)
	rootCmd.AddCommand(ReviewCmd)

	rootCmd.AddGroup(&cobra.Group{
		ID:    "core",
		Title: "Core Commands",
	})

	rootCmd.AddGroup(&cobra.Group{
		ID:    "ai",
		Title: "AI-Powered Commands",
	})
}
