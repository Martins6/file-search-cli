package filesystem

import (
	"os"
	"sort"
)

func ScanDirectory(path string) ([]Entry, error) {
	entries, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}

	result := make([]Entry, 0, len(entries))
	for _, dirEntry := range entries {
		info, err := dirEntry.Info()
		if err != nil {
			continue
		}

		entry := Entry{
			Name:    dirEntry.Name(),
			Path:    path + string(os.PathSeparator) + dirEntry.Name(),
			IsDir:   dirEntry.IsDir(),
			Size:    info.Size(),
			ModTime: info.ModTime(),
		}
		result = append(result, entry)
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].Name < result[j].Name
	})

	return result, nil
}
