package cmd

import (
	"fmt"
	"os"

	"fsc/pkg/editor"
	"fsc/pkg/ui"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
)

var (
	showHidden bool
	dirsOnly   bool
	filesOnly  bool
)

var rootCmd = &cobra.Command{
	Use:   "fsc [directory]",
	Short: "Finder-like CLI tool for interactive file browsing",
	Long: `fsc is an interactive CLI tool that provides a Finder-like browsing experience 
in the terminal. It allows you to navigate directories, filter files, and open items 
in your editor.

  Examples:
  fsc /path/to/dir         Start browsing in specified directory
  fsc .                    Start in current directory
  fsc ~                    Start in home directory
  fsc /tmp --hidden        Show hidden files
  fsc /docs --dirs-only    Show only directories`,
	Args: cobra.MaximumNArgs(1),
	PreRunE: func(cmd *cobra.Command, args []string) error {
		path := "."
		if len(args) > 0 {
			path = args[0]
		}
		info, err := os.Stat(path)
		if err != nil {
			return fmt.Errorf("path %s does not exist: %w", path, err)
		}
		if !info.IsDir() {
			return fmt.Errorf("path %s is not a directory", path)
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		directory := "."
		if len(args) > 0 {
			directory = args[0]
		}

		model := ui.InitialModel(directory, showHidden, dirsOnly, filesOnly)

		p := tea.NewProgram(
			model,
			tea.WithAltScreen(),
		)

		finalModel, err := p.Run()
		if err != nil {
			fmt.Printf("Error running program: %v\n", err)
			os.Exit(1)
		}

		if finalModel != nil {
			if m, ok := finalModel.(ui.Model); ok {
				if len(m.FilteredEntries()) > 0 {
					selectedPath := m.FilteredEntries()[m.SelectedIndex()].Path
					editor.OpenInEditor(selectedPath)
				}
			}
		}
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().BoolVarP(&showHidden, "hidden", "H", false, "show hidden files (default: false)")
	rootCmd.Flags().BoolVarP(&dirsOnly, "dirs-only", "d", false, "show only directories")
	rootCmd.Flags().BoolVarP(&filesOnly, "files-only", "f", false, "show only files")
}
