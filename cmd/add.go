package cmd

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
)

var AddCmd = &cobra.Command{
	Use:     "add",
	Short:   "Add a new transaction interactively",
	GroupID: "core",
	Run: func(cmd *cobra.Command, args []string) {
		var finalEntry string

		p := tea.NewProgram(initialModel(), tea.WithOutput(os.Stdout))
		finalModel, err := p.Run()
		if err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}

		if m, ok := finalModel.(model); ok && m.submit {
			finalEntry = m.JournalString()

			if m.appendToFile && finalEntry != "" {
				f, err := os.OpenFile(journalFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
				if err != nil {
					fmt.Println("Failed to open journal file:", err)
					os.Exit(1)
				}
				defer f.Close()

				if _, err := f.WriteString("\n" + finalEntry); err != nil {
					fmt.Println("Failed to write to journal file:", err)
					os.Exit(1)
				}

				fmt.Println(successStyle.Render("Entry successfully appended to journal file."))
			}
		}
	},
}

var (
	titleStyle   = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#81A1C1"))
	labelStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("#88C0D0"))
	successStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#A3BE8C")).Bold(true)
	errorStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("#BF616A")).Bold(true)
	boxStyle     = lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color("#4C566A")).Padding(1, 2).Margin(1, 0).Width(60)
)

type accountItem string

func (a accountItem) Title() string       { return string(a) }
func (a accountItem) Description() string { return "" }
func (a accountItem) FilterValue() string { return string(a) }

type entryField struct {
	account textinput.Model
	amount  textinput.Model
	chooser list.Model
	picking bool
	focused bool
	id      int
}

func styledList(items []list.Item) list.Model {
	l := list.New(items, list.NewDefaultDelegate(), 30, 5)
	l.SetShowTitle(false)
	l.Styles.PaginationStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("241"))
	l.Styles.HelpStyle = lipgloss.NewStyle().Faint(true)
	return l
}

func newEntryField(id int) entryField {
	accInput := textinput.New()
	accInput.Placeholder = "Account"
	accInput.Prompt = ""

	amtInput := textinput.New()
	amtInput.Placeholder = fmt.Sprintf("Amount %d (e.g. 1000€)", id+1)

	accounts := getAccounts(journalFilePath)
	items := make([]list.Item, len(accounts))
	for i, acc := range accounts {
		items[i] = accountItem(acc)
	}

	chooser := styledList(items)
	chooser.Title = "Select Account"

	return entryField{
		account: accInput,
		amount:  amtInput,
		chooser: chooser,
		id:      id,
	}
}

type model struct {
	date         textinput.Model
	desc         textinput.Model
	entries      []entryField
	focusIdx     int
	errMsg       string
	submit       bool
	appendToFile bool
}

func initialModel() model {
	date := textinput.New()
	date.Placeholder = "Date (YYYY-MM-DD)"
	date.SetValue(time.Now().Format("2006-01-02"))
	date.Focus()

	desc := textinput.New()
	desc.Placeholder = "Description"

	entries := []entryField{
		newEntryField(0),
		newEntryField(1),
	}

	return model{
		date:    date,
		desc:    desc,
		entries: entries,
	}
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func (m *model) focusCurrent() {
	m.date.Blur()
	m.desc.Blur()
	for i := range m.entries {
		m.entries[i].account.Blur()
		m.entries[i].amount.Blur()
		m.entries[i].focused = false
	}

	switch {
	case m.focusIdx == 0:
		m.date.Focus()
	case m.focusIdx == 1:
		m.desc.Focus()
	default:
		idx := (m.focusIdx - 2) / 2
		if (m.focusIdx-2)%2 == 0 {
			m.entries[idx].account.Focus()
			m.entries[idx].focused = true
		} else {
			m.entries[idx].amount.Focus()
		}
	}
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if m.submit {
		return m, tea.Quit
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		key := msg.String()

		for i := range m.entries {
			if m.entries[i].picking {
				var cmd tea.Cmd
				m.entries[i].chooser, cmd = m.entries[i].chooser.Update(msg)
				if key == "enter" {
					if choice, ok := m.entries[i].chooser.SelectedItem().(accountItem); ok {
						m.entries[i].account.SetValue(string(choice))
					}
					m.entries[i].picking = false
				}
				return m, cmd
			}
		}

		switch key {
		case "ctrl+c", "q":
			return m, tea.Quit

		case "tab", "down":
			m.focusIdx = (m.focusIdx + 1) % (2 + len(m.entries)*2)
			m.focusCurrent()

		case "up":
			m.focusIdx = (m.focusIdx - 1 + (2 + len(m.entries)*2)) % (2 + len(m.entries)*2)
			m.focusCurrent()

		case "ctrl+n":
			newEntry := newEntryField(len(m.entries))
			m.entries = append(m.entries, newEntry)
			m.focusIdx = 2 + (len(m.entries)-1)*2 // Focus the new account input
			m.focusCurrent()
			return m, nil

		case "ctrl+a":
			if err := m.validate(); err != "" {
				m.errMsg = err
				return m, nil
			}
			m.appendToFile = true
			m.submit = true
			m.errMsg = ""
			return m, tea.Quit

		case "ctrl+o":
			if err := m.validate(); err != "" {
				m.errMsg = err
				return m, nil
			}
			m.appendToFile = false
			m.submit = true
			m.errMsg = ""
			return m, tea.Quit

		case "enter":
			for i := range m.entries {
				if m.entries[i].focused {
					m.entries[i].picking = true
					m.entries[i].chooser.Select(0)
					return m, nil
				}
			}
		}

		if m.date.Focused() {
			var cmd tea.Cmd
			m.date, cmd = m.date.Update(msg)
			return m, cmd
		}
		if m.desc.Focused() {
			var cmd tea.Cmd
			m.desc, cmd = m.desc.Update(msg)
			return m, cmd
		}
		for i := range m.entries {
			if m.entries[i].account.Focused() {
				var cmd tea.Cmd
				m.entries[i].account, cmd = m.entries[i].account.Update(msg)
				return m, cmd
			}
			if m.entries[i].amount.Focused() {
				var cmd tea.Cmd
				m.entries[i].amount, cmd = m.entries[i].amount.Update(msg)
				return m, cmd
			}
		}
	}

	return m, nil
}

func (m *model) validate() string {
	if len(m.entries) < 2 {
		return "At least 2 account/amount pairs required."
	}
	if strings.TrimSpace(m.date.Value()) == "" || strings.TrimSpace(m.desc.Value()) == "" {
		return "Date and description are required."
	}

	total := 0.0
	lastIdx := len(m.entries) - 1

	for i, e := range m.entries {
		acc := strings.TrimSpace(e.account.Value())
		amt := strings.TrimSpace(e.amount.Value())

		if acc == "" {
			return "All account fields must be filled."
		}

		if amt == "" {
			if i == lastIdx {
				continue
			}
			return "All amount fields must be filled."
		}

		var val float64
		if _, err := fmt.Sscanf(amt, "%f", &val); err != nil {
			return "Amounts must be numeric."
		}
		total += val
	}

	return ""
}

func (m model) View() string {
	for _, e := range m.entries {
		if e.picking {
			return e.chooser.View()
		}
	}

	if m.submit {
		return m.JournalString()
	}

	var b strings.Builder
	b.WriteString(labelStyle.Render("Date:        ") + m.date.View() + "\n")
	b.WriteString(labelStyle.Render("Description: ") + m.desc.View() + "\n\n")

	for i, e := range m.entries {
		b.WriteString(labelStyle.Render(fmt.Sprintf("Account %d:   ", i+1)) + e.account.View() + "\n")
		b.WriteString(labelStyle.Render(fmt.Sprintf("Amount  %d:   ", i+1)) + e.amount.View() + "\n\n")
	}

	if m.errMsg != "" {
		b.WriteString(errorStyle.Render("Error: " + m.errMsg + "\n\n"))
	}

	footer := strings.Join([]string{
		labelStyle.Render("[ctrl+a] Generate and append to journal"),
		labelStyle.Render("[ctrl+o] Generate and output to console"),
		labelStyle.Render("[ctrl+n] Add another account/amount pair"),
		labelStyle.Render("[q]      Quit"),
	}, "\n")
	b.WriteString(footer)

	return boxStyle.Render(b.String())
}

func (m model) JournalString() string {
	var out strings.Builder
	fmt.Fprintf(&out, "%s %s\n", m.date.Value(), m.desc.Value())
	for _, e := range m.entries {
		fmt.Fprintf(&out, "  %s    %s\n", e.account.Value(), e.amount.Value())
	}
	return out.String()
}
