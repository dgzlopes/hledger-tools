package cmd

import (
	"github.com/spf13/cobra"
)

var journalFilePath string

var rootCmd = &cobra.Command{
	Use:              "hledger-tools",
	Short:            "A suite of utilities to complement hledger",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {},
}

func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&journalFilePath, "file", "f", "", "Path to journal file")

	rootCmd.AddCommand(AddCmd)
}
