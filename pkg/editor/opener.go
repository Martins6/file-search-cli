package editor

import (
	"fmt"
	"os"
	"os/exec"
)

func OpenInEditor(path string) error {
	editor := os.Getenv("EDITOR")
	if editor == "" {
		editor = os.Getenv("VISUAL")
	}

	editors := []string{editor, "vim", "nano", "code", "vi", "edit"}

	for _, e := range editors {
		if e == "" {
			continue
		}
		if _, err := exec.LookPath(e); err == nil {
			cmd := exec.Command(e, path)
			cmd.Stdin = os.Stdin
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			return cmd.Run()
		}
	}

	return fmt.Errorf("no editor found (checked: EDITOR, VISUAL, vim, nano, code, vi, edit)")
}
