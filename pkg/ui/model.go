package ui

import (
	"fsc/pkg/filesystem"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	directory          string
	entries            []filesystem.Entry
	filteredEntries    []filesystem.Entry
	selectedIndex      int
	searchQuery        string
	showHidden         bool
	dirsOnly           bool
	filesOnly          bool
	vimMode            bool
	width              int
	height             int
	keyMap             KeyMap
	pendingKeySequence string
	openEditorOnQuit   bool
	selectedFileOnQuit bool
	selectedFilePath   string
}

type initialModelMsg struct {
	entries []filesystem.Entry
	path    string
}

type clearPendingSequenceMsg struct{}

func InitialModel(directory string, showHidden, dirsOnly, filesOnly bool) Model {
	return Model{
		directory:          directory,
		searchQuery:        "",
		showHidden:         showHidden,
		dirsOnly:           dirsOnly,
		filesOnly:          filesOnly,
		vimMode:            false,
		selectedIndex:      0,
		keyMap:             DefaultKeyMap(),
		openEditorOnQuit:   false,
		selectedFileOnQuit: false,
		selectedFilePath:   "",
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

func ClearPendingSequenceCmd() tea.Cmd {
	return tea.Tick(500*time.Millisecond, func(t time.Time) tea.Msg {
		return clearPendingSequenceMsg{}
	})
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "/" {
			if m.pendingKeySequence == "" {
				m.pendingKeySequence = "/"
				return m, ClearPendingSequenceCmd()
			}
		}

		if m.pendingKeySequence != "" {
			if msg.String() == "e" && m.pendingKeySequence == "/" {
				m.pendingKeySequence = ""
				m.openEditorOnQuit = true
				if len(m.filteredEntries) > 0 {
					selectedEntry := m.filteredEntries[m.selectedIndex]
					return m, func() tea.Msg {
						return OpenEditorMsg(selectedEntry.Path)
					}
				}
				return m, nil
			} else if msg.String() == "v" && m.pendingKeySequence == "/" {
				m.vimMode = !m.vimMode
				m.pendingKeySequence = ""
				return m, nil
			} else if msg.String() == "q" && m.pendingKeySequence == "/" {
				m.pendingKeySequence = ""
				return m, tea.Quit
			} else {
				m.searchQuery += m.pendingKeySequence
				if msg.Type == tea.KeyRunes {
					m.searchQuery += string(msg.Runes)
				}
				m.filteredEntries = m.filterEntries()
				m.selectedIndex = 0
				m.pendingKeySequence = ""
				return m, nil
			}
		}

		switch msg.String() {
		case m.keyMap.Quit, "ctrl+c", "esc":
			return m, tea.Quit
		case m.keyMap.Up:
			if len(m.filteredEntries) > 0 && m.selectedIndex > 0 {
				m.selectedIndex--
			}
		case m.keyMap.Down:
			if len(m.filteredEntries) > 0 && m.selectedIndex < len(m.filteredEntries)-1 {
				m.selectedIndex++
			}
		case m.keyMap.Enter, "right":
			if len(m.filteredEntries) > 0 {
				selectedEntry := m.filteredEntries[m.selectedIndex]
				if selectedEntry.IsDir {
					m.pendingKeySequence = ""
					return m, func() tea.Msg {
						return NavigateMsg(selectedEntry.Path)
					}
				} else {
					m.selectedFileOnQuit = true
					m.selectedFilePath = selectedEntry.Path
					return m, tea.Quit
				}
			}
		case m.keyMap.Left:
			m.pendingKeySequence = ""
			return m, func() tea.Msg {
				_, parentPath, _ := filesystem.NavigateUp(m.directory)
				return NavigateMsg(parentPath)
			}
		case m.keyMap.Backspace:
			m.pendingKeySequence = ""
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

		if m.vimMode {
			switch msg.String() {
			case "k":
				if len(m.filteredEntries) > 0 && m.selectedIndex > 0 {
					m.selectedIndex--
				}
			case "j":
				if len(m.filteredEntries) > 0 && m.selectedIndex < len(m.filteredEntries)-1 {
					m.selectedIndex++
				}
			case "l":
				if len(m.filteredEntries) > 0 {
					selectedEntry := m.filteredEntries[m.selectedIndex]
					if selectedEntry.IsDir {
						m.pendingKeySequence = ""
						return m, func() tea.Msg {
							return NavigateMsg(selectedEntry.Path)
						}
					} else {
						m.selectedFileOnQuit = true
						m.selectedFilePath = selectedEntry.Path
						return m, tea.Quit
					}
				}
			case "h":
				m.pendingKeySequence = ""
				return m, func() tea.Msg {
					_, parentPath, _ := filesystem.NavigateUp(m.directory)
					return NavigateMsg(parentPath)
				}
			}
		}

	case clearPendingSequenceMsg:
		if m.pendingKeySequence != "" {
			m.searchQuery += m.pendingKeySequence
			m.filteredEntries = m.filterEntries()
			m.selectedIndex = 0
			m.pendingKeySequence = ""
		}

	case initialModelMsg:
		m.entries = msg.entries
		m.directory = msg.path
		m.filteredEntries = m.filterEntries()

	case NavigateMsg:
		m.pendingKeySequence = ""
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

func (m Model) ShouldOpenEditor() bool {
	return m.openEditorOnQuit
}

func (m Model) HasSelectedFile() bool {
	return m.selectedFileOnQuit
}

func (m Model) SelectedFilePath() string {
	return m.selectedFilePath
}
