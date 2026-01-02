package ui

import (
	"fsc/pkg/filesystem"
	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	directory       string
	entries         []filesystem.Entry
	filteredEntries []filesystem.Entry
	selectedIndex   int
	searchQuery     string
	showHidden      bool
	dirsOnly        bool
	filesOnly       bool
	width           int
	height          int
	keyMap          KeyMap
}

type initialModelMsg struct {
	entries []filesystem.Entry
	path    string
}

func InitialModel(directory string, showHidden, dirsOnly, filesOnly bool) Model {
	return Model{
		directory:     directory,
		searchQuery:   "",
		showHidden:    showHidden,
		dirsOnly:      dirsOnly,
		filesOnly:     filesOnly,
		selectedIndex: 0,
		keyMap:        DefaultKeyMap(),
	}
}

func (m Model) Init() tea.Cmd {
	return func() tea.Msg {
		entries, path, err := filesystem.NavigateTo(m.directory)
		if err != nil {
			return tea.Quit
		}
		return initialModelMsg{entries: entries, path: path}
	}
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case m.keyMap.Quit, "ctrl+c", "esc":
			return m, tea.Quit
		case m.keyMap.Up, "k":
			if len(m.filteredEntries) > 0 && m.selectedIndex > 0 {
				m.selectedIndex--
			}
		case m.keyMap.Down, "j":
			if len(m.filteredEntries) > 0 && m.selectedIndex < len(m.filteredEntries)-1 {
				m.selectedIndex++
			}
		case m.keyMap.Enter, "right", "l":
			if len(m.filteredEntries) > 0 {
				selectedEntry := m.filteredEntries[m.selectedIndex]
				if selectedEntry.IsDir {
					return m, func() tea.Msg {
						return NavigateMsg(selectedEntry.Path)
					}
				}
			}
		case m.keyMap.OpenEditor:
			if len(m.filteredEntries) > 0 {
				selectedEntry := m.filteredEntries[m.selectedIndex]
				return m, func() tea.Msg {
					return OpenEditorMsg(selectedEntry.Path)
				}
			}
		case m.keyMap.Backspace:
			if len(m.searchQuery) > 0 {
				m.searchQuery = m.searchQuery[:len(m.searchQuery)-1]
				m.filteredEntries = m.filterEntries()
				m.selectedIndex = 0
			}
		default:
			if msg.Type == tea.KeyRunes {
				m.searchQuery += string(msg.Runes)
				m.filteredEntries = m.filterEntries()
				m.selectedIndex = 0
			}
		}

	case initialModelMsg:
		m.entries = msg.entries
		m.directory = msg.path
		m.filteredEntries = m.filterEntries()

	case NavigateMsg:
		entries, path, err := filesystem.NavigateTo(string(msg))
		if err == nil {
			m.entries = entries
			m.directory = path
			m.searchQuery = ""
			m.filteredEntries = m.filterEntries()
			m.selectedIndex = 0
		}

	case OpenEditorMsg:
		return m, tea.Quit

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	}

	return m, nil
}

func (m Model) filterEntries() []filesystem.Entry {
	filtered := m.entries

	if !m.showHidden {
		filtered = filesystem.FilterHidden(filtered)
	}

	filtered = filesystem.FilterByType(filtered, m.dirsOnly, m.filesOnly)

	filtered = filesystem.FilterByPattern(filtered, m.searchQuery)

	return filtered
}

func (m Model) FilteredEntries() []filesystem.Entry {
	return m.filteredEntries
}

func (m Model) SelectedIndex() int {
	return m.selectedIndex
}
