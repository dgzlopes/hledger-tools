# hledger-tools

This repo contains a collection of small tools for working with hledger journals.

- `add` – A better `hledger add` (with account selector and more)
- `import` – Import transactions from CSV or similar files *(AI-powered)*
- `review` – Analyze your journal and get financial feedback *(AI-powered)*

AI-powered tools use OpenAI’s APIs and require the `OPENAI_API_KEY` environment variable to be set. You can feed additional context to the LLM via the `--context` flag (e.g. financial goals or categorization rules).

---

## Usage

```bash
# Interactively add a transaction to your journal
hledger-tools add --journal 2025.journal

# Import transactions from a CSV file
hledger-tools import transactions.csv --journal 2025.journal

# Review your journal and get AI-generated feedback
hledger-tools review --journal 2025.journal --context who-am-i.txt
```