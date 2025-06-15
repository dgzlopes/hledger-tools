# hledger-tools

This repo contains a collection of small tools for working with hledger journals.

- `add` – A better way to interactively add transactions to your journal
- `import` – Import transactions from CSV or similar files *(AI-powered)*
- `ask` – Ask a question about your balance sheet *(AI-powered)*

AI-powered tools use OpenAI’s API and require the `OPENAI_API_KEY` environment variable to be set. You can provide additional context to the LLM using the `--context` flag (e.g. financial goals, categorization rules, etc.).


## Usage

```bash
# Interactively add a transaction to your journal
hledger-tools add --journal 2025.journal

# Import transactions from a CSV file
hledger-tools import transactions.csv --journal 2025.journal

# Ask a question about your balance sheet, with optional context
hledger-tools ask "Am I saving enough?" --journal 2025.journal --context who-am-i.txt
```