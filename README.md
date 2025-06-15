# hledger-tools

This repo contains a collection of small tools for working with hledger journals.

- `add` – A better way to interactively add transactions to your journal  
- `import` – Import transactions from CSV or similar files *(AI-powered)*  
- `ask` – Ask a question about your balance sheet *(AI-powered)*

## AI?!

AI-powered tools (`import`, `ask`) use OpenAI’s API and require the `OPENAI_API_KEY` environment variable.

They are implemented to be as "safe" as possible:
- `ask` sends your **balance sheet** (with percentages) and **account names**.
- `import` sends your **account names** and **source file** (with raw transactions).  
  - You *can* send your **entire journal** with `--include-journal`, but you probably shouldn’t.

Use `--show-prompt` to preview the full LLM prompt.  
Use `--context` to add extra info (e.g. goals or categorization rules).

## Usage examples

```bash
# Interactively add a transaction to your journal
hledger-tools add --journal 2025.journal
╭────────────────────────────────────────────────────────────╮
│                                                            │
│  Date:        > 2025-06-16                                 │
│  Description: > My salary                                  │
│                                                            │
│  Account 1:   assets:cash                                  │
│  Amount  1:   > 100€                                       │
│                                                            │
│  Account 2:   assets:savings:bankB                         │
│  Amount  2:   > A                                          │
│                                                            │
│  [ctrl+a] Generate and append to journal                   │
│  [ctrl+o] Generate and output to console                   │
│  [ctrl+n] Add another account/amount pair                  │
│  [q]      Quit                                             │
│                                                            │
╰────────────────────────────────────────────────────────────╯
```

```bash
# Import transactions from a CSV file
hledger-tools import transactions.csv --journal 2025.journal
2025-06-23 * "Dinner at Restaurant"
    expenses:fun  €60.00
    assets:cash

2025-06-28 * "End of Month Savings"
    assets:savings:bankA  €250.00
    assets:cash
```

```bash
# Ask a question about your balance sheet, with optional context
hledger-tools ask "Am I saving enough?" --journal 2025.journal --context who-am-i.txt
Roughly 18% of your assets are in savings accounts and 22% in investment funds.  
....... more things ....... 
```