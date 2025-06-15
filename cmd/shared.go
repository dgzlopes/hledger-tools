package cmd

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/siddhantac/hledger"
)

func getAccounts(journalFilePath string) []string {
	h := hledger.New("hledger", journalFilePath)
	reader, err := h.Accounts(hledger.NewOptions())
	if err != nil {
		fmt.Println("Failed to load accounts:", err)
		os.Exit(1)
	}

	data, err := io.ReadAll(reader)
	if err != nil {
		panic(err)
	}

	content := string(data)

	return strings.Split(strings.TrimSpace(content), "\n")
}

func getBalanceSheet(journalFilePath string) string {
	h := hledger.New("hledger", journalFilePath)
	reader, err := h.BalanceSheet(hledger.NewOptions().WithPercent(true).WithTree(true))
	if err != nil {
		fmt.Println("Failed to load balance sheet:", err)
		os.Exit(1)
	}

	data, err := io.ReadAll(reader)
	if err != nil {
		panic(err)
	}

	return string(data)
}
