package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var reviewContextPath string

var ReviewCmd = &cobra.Command{
	Use:     "review",
	Short:   "Analyze your balance sheet and get feedback on your financial health",
	GroupID: "ai",
	Run: func(cmd *cobra.Command, args []string) {
		if reviewContextPath != "" {
			if _, err := os.Stat(reviewContextPath); os.IsNotExist(err) {
				fmt.Printf("Warning: context file not found: %s (ignored)\n", reviewContextPath)
			} else {
				fmt.Printf("Using context file: %s\n", reviewContextPath)
			}
		}

		fmt.Println("TODO")
	},
}

func init() {
	ReviewCmd.Flags().StringVarP(&reviewContextPath, "context", "c", "", "Optional context file. It will be included in the LLM prompt.")
}
