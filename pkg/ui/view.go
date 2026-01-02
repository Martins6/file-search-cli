package ui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
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

	modalStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#7D56F4")).
			Padding(1, 2)
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
	var filterLabel string
	if m.regexMode {
		filterLabel = "Filter (Regex): "
	} else {
		filterLabel = "Filter: "
	}
	b.WriteString(filterStyle.Render(filterLabel + m.searchQuery + "_"))
	b.WriteString("\n")

	var helpText string
	if m.vimMode {
		helpText = "j k h l | /h: help"
	} else {
		helpText = "↑ ↓ ← → | /h: help"
	}
	b.WriteString(helpStyle.Render(helpText))
	b.WriteString("\n")

	if m.showHelp {
		return m.renderHelpModal()
	}

	return b.String()
}

func (m Model) renderHelpModal() string {
	var content strings.Builder

	content.WriteString(titleStyle.Render("Help"))
	content.WriteString("\n\n")

	content.WriteString("Navigation:\n")
	if m.vimMode {
		content.WriteString("• j - move down\n")
		content.WriteString("• k - move up\n")
		content.WriteString("• h - parent directory\n")
		content.WriteString("• l - enter directory / select file\n")
	} else {
		content.WriteString("• ↑ - move up\n")
		content.WriteString("• ↓ - move down\n")
		content.WriteString("• ← - parent directory\n")
		content.WriteString("• → - enter directory / select file\n")
	}
	content.WriteString("\n")

	content.WriteString("Actions:\n")
	content.WriteString("• enter - navigate to directory / select file\n")
	content.WriteString("\n")

	content.WriteString("Commands:\n")
	content.WriteString("• /e - open editor\n")
	content.WriteString("• /p - print path\n")
	content.WriteString("• /r - toggle regex mode\n")
	content.WriteString("• /v - toggle vim mode\n")
	content.WriteString("• /q - quit\n")
	content.WriteString("• /h - toggle help\n")
	content.WriteString("\n")

	content.WriteString(dimStyle.Render("Press esc to close"))

	return modalStyle.Width(60).Align(lipgloss.Left, lipgloss.Top).Render(content.String())
}
