package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var contextPath string

var ImportCmd = &cobra.Command{
	Use:     "import <source>",
	Short:   "Import transactions from any source (e.g., CSV, JSON, etc.)",
	GroupID: "ai",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return fmt.Errorf("missing source file path (e.g., import transactions.csv)")
		}
		if len(args) > 1 {
			return fmt.Errorf("too many arguments — expected just the source file")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		sourcePath := args[0]

		if _, err := os.Stat(sourcePath); os.IsNotExist(err) {
			fmt.Printf("Error: source file not found: %s\n", sourcePath)
			os.Exit(1)
		}

		fmt.Printf("Importing from source: %s\n", sourcePath)

		if contextPath != "" {
			if _, err := os.Stat(contextPath); os.IsNotExist(err) {
				fmt.Printf("Warning: context file not found: %s (ignored)\n", contextPath)
			} else {
				fmt.Printf("Using context file: %s\n", contextPath)
			}
		}

		fmt.Println("TODO")
	},
}

func init() {
	ImportCmd.Flags().StringVarP(&contextPath, "context", "c", "", "Path to optional context file. It will be included in the LLM prompt.")
}
