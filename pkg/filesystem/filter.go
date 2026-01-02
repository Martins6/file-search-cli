package filesystem

import (
	"regexp"
	"sort"
	"strings"
)

func FilterByPattern(entries []Entry, pattern string) []Entry {
	if pattern == "" {
		return entries
	}

	var result []Entry

	if strings.HasPrefix(pattern, "/") {
		regexStr := strings.TrimPrefix(pattern, "/")
		re, err := regexp.Compile(regexStr)
		if err != nil {
			return entries
		}
		for _, entry := range entries {
			if re.MatchString(entry.Name) {
				result = append(result, entry)
			}
		}
	} else {
		lowerPattern := strings.ToLower(pattern)
		for _, entry := range entries {
			if strings.Contains(strings.ToLower(entry.Name), lowerPattern) {
				result = append(result, entry)
			}
		}
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].Name < result[j].Name
	})

	return result
}

func FilterByType(entries []Entry, showDirsOnly, showFilesOnly bool) []Entry {
	if !showDirsOnly && !showFilesOnly {
		return entries
	}

	var result []Entry
	for _, entry := range entries {
		if showDirsOnly && entry.IsDir {
			result = append(result, entry)
		}
		if showFilesOnly && !entry.IsDir {
			result = append(result, entry)
		}
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].Name < result[j].Name
	})

	return result
}

func FilterHidden(entries []Entry) []Entry {
	var result []Entry
	for _, entry := range entries {
		if !strings.HasPrefix(entry.Name, ".") {
			result = append(result, entry)
		}
	}
	return result
}
