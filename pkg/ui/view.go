package ui

import (
	"github.com/charmbracelet/lipgloss"
	"strings"
)

var (
	titleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FAFAFA")).
			Background(lipgloss.Color("#7D56F4")).
			Padding(0, 1)

	pathStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FAFAFA")).
			Background(lipgloss.Color("#1E1E2E")).
			Padding(0, 1)

	separatorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#45475A")).
			Width(50)

	selectedStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#1E1E2E")).
			Background(lipgloss.Color("#89B4FA")).
			Bold(true)

	dimStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#6C7086"))

	filterStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#F9E2AF")).
			Background(lipgloss.Color("#1E1E2E")).
			Padding(0, 1)

	helpStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#6C7086")).
			Background(lipgloss.Color("#1E1E2E")).
			Padding(0, 1)
)

func (m Model) View() string {
	var b strings.Builder

	b.WriteString(titleStyle.Render("fsc - File Search CLI"))
	b.WriteString("\n")

	b.WriteString(pathStyle.Render(m.directory))
	b.WriteString("\n")

	b.WriteString(separatorStyle.Render(strings.Repeat("─", 50)))
	b.WriteString("\n")

	if len(m.filteredEntries) == 0 {
		b.WriteString(dimStyle.Render("  No entries match filter"))
		b.WriteString("\n")
	} else {
		for i, entry := range m.filteredEntries {
			line := "  "
			if entry.IsDir {
				line += "📁 "
			} else {
				line += "📄 "
			}
			line += entry.Name

			if i == m.selectedIndex {
				line = selectedStyle.Render(line)
			}
			b.WriteString(line)
			b.WriteString("\n")
		}
	}

	b.WriteString("\n")
	b.WriteString(filterStyle.Render("Filter: " + m.searchQuery + "_"))
	b.WriteString("\n")

	helpText := "↑/j: up | ↓/k: down | enter: navigate | e: open editor | q: quit"
	b.WriteString(helpStyle.Render(helpText))
	b.WriteString("\n")

	return b.String()
}
