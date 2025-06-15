package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/openai/openai-go"
	"github.com/spf13/cobra"
)

var (
	askContextPath string
	askShowPrompt  bool
)

var AskCmd = &cobra.Command{
	Use:     "ask <question>",
	Short:   "Ask a question about your balance sheet",
	GroupID: "ai",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return fmt.Errorf("missing question (e.g., ask \"Am I saving enough?\")")
		}
		if len(args) > 1 {
			return fmt.Errorf("too many arguments — expected exactly one question")
		}
		return nil
	},

	Run: func(cmd *cobra.Command, args []string) {
		question := args[0]

		var contextContent string
		if askContextPath != "" {
			data, err := os.ReadFile(askContextPath)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Warning: failed to read context file: %s (ignored)\n", askContextPath)
			} else {
				contextContent = string(data)
				fmt.Printf("Using context file: %s\n", askContextPath)
			}
		}

		balance := getBalanceSheet(journalFilePath)

		systemMsg := "You are a helpful financial assistant reviewing a user's personal finances."
		userPrompt := fmt.Sprintf(`Here is my balance sheet:

%s

Context:
%s

Question:
%s`, balance, contextContent, question)

		if askShowPrompt {
			fmt.Println("----- Prompt sent to OpenAI -----")
			fmt.Println(userPrompt)
			fmt.Println("----- End of prompt -----")
			fmt.Println()
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
	AskCmd.Flags().StringVarP(&askContextPath, "context", "c", "", "Optional context file. It will be included in the LLM prompt.")
	AskCmd.Flags().BoolVar(&askShowPrompt, "show-prompt", false, "Print the full prompt sent to the LLM")
}
