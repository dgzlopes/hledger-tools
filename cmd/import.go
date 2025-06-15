package cmd

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/openai/openai-go"
	"github.com/spf13/cobra"
)

var contextPath string
var showPrompt bool
var includeJournal bool

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

		sourceData, err := os.ReadFile(sourcePath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: could not read source file: %s\n", sourcePath)
			os.Exit(1)
		}

		var contextContent string
		if contextPath != "" {
			data, err := os.ReadFile(contextPath)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Warning: failed to read context file: %s (ignored)\n", contextPath)
			} else {
				contextContent = string(data)
				fmt.Printf("Using context file: %s\n", contextPath)
			}
		}

		var journalContent string
		if includeJournal {
			journalBytes, err := os.ReadFile(journalFilePath)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Warning: failed to read journal file: %s (ignored)\n", journalFilePath)
			} else {
				journalContent = fmt.Sprintf("\nHere is my full journal for context:\n\n%s", string(journalBytes))
			}
		}

		accounts := getAccounts(journalFilePath)

		systemMsg := "You are a financial assistant that converts raw transactions into hledger journal entries."
		userPrompt := fmt.Sprintf(`Here are the account names available:

%s

Here is the source data (e.g., CSV or JSON):

%s

Context:
%s%s

Generate valid hledger journal transactions using only the accounts above.

IMPORTANT:
- Do NOT explain anything.
- Do NOT include any prose or formatting.
- Output ONLY valid hledger journal entries.
`, strings.Join(accounts, "\n"), string(sourceData), contextContent, journalContent)

		if showPrompt {
			fmt.Println("----- Prompt sent to OpenAI -----")
			fmt.Println(userPrompt)
			fmt.Println("----- End of prompt -----")
		}

		if os.Getenv("OPENAI_API_KEY") == "" {
			fmt.Fprintln(os.Stderr, "Error: OPENAI_API_KEY environment variable is not set.")
			os.Exit(1)
		}

		client := openai.NewClient()
		resp, err := client.Chat.Completions.New(context.Background(), openai.ChatCompletionNewParams{
			Model: openai.ChatModelGPT4o,
			Messages: []openai.ChatCompletionMessageParamUnion{
				openai.SystemMessage(systemMsg),
				openai.UserMessage(userPrompt),
			},
		})
		if err != nil {
			fmt.Fprintf(os.Stderr, "OpenAI request failed: %v\n", err)
			os.Exit(1)
		}

		fmt.Println()
		fmt.Println(resp.Choices[0].Message.Content)
	},
}

func init() {
	ImportCmd.Flags().StringVarP(&contextPath, "context", "c", "", "Path to optional context file. It will be included in the LLM prompt.")
	ImportCmd.Flags().BoolVar(&showPrompt, "show-prompt", false, "Print the full prompt sent to the LLM")
	ImportCmd.Flags().BoolVar(&includeJournal, "include-journal", false, "Include the full journal in the LLM prompt (NOT RECOMMENDED)")
}
