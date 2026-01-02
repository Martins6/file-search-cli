package filesystem

import (
	"path/filepath"
)

func NavigateTo(path string) ([]Entry, string, error) {
	entries, err := ScanDirectory(path)
	if err != nil {
		return nil, "", err
	}
	absPath, err := filepath.Abs(path)
	if err != nil {
		return nil, "", err
	}
	return entries, absPath, nil
}

func NavigateUp(currentPath string) ([]Entry, string, error) {
	parentPath := filepath.Dir(currentPath)
	if parentPath == currentPath {
		return NavigateTo(currentPath)
	}
	return NavigateTo(parentPath)
}
