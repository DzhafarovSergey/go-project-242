package code

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func GetPathSize(path string, recursive, human, all bool) (string, error) {
	size, err := calculateSize(path, recursive, all)
	if err != nil {
		return "", err
	}

	if human {
		return formatSizeHuman(size), nil
	}
	return fmt.Sprintf("%dB", size), nil
}

func calculateSize(path string, recursive, all bool) (int64, error) {
	info, err := os.Stat(path)
	if err != nil {
		return 0, err
	}

	if !info.IsDir() {
		return info.Size(), nil
	}

	var totalSize int64

	entries, err := os.ReadDir(path)
	if err != nil {
		return 0, err
	}

	for _, entry := range entries {
		entryName := entry.Name()
		if !all && strings.HasPrefix(entryName, ".") {
			continue
		}

		fullPath := filepath.Join(path, entryName)

		if entry.IsDir() && recursive {
			subDirSize, err := calculateSize(fullPath, recursive, all)
			if err == nil {
				totalSize += subDirSize
			}
		} else if !entry.IsDir() {
			entryInfo, err := entry.Info()
			if err != nil {
				continue
			}
			totalSize += entryInfo.Size()
		}
	}

	return totalSize, nil
}

func formatSizeHuman(size int64) string {
	if size == 0 {
		return "0B"
	}

	const (
		_  = iota
		KB = 1 << (10 * iota)
		MB
		GB
		TB
		PB
		EB
	)

	switch {
	case size >= EB:
		return fmt.Sprintf("%.1fEB", float64(size)/float64(EB))
	case size >= PB:
		return fmt.Sprintf("%.1fPB", float64(size)/float64(PB))
	case size >= TB:
		return fmt.Sprintf("%.1fTB", float64(size)/float64(TB))
	case size >= GB:
		return fmt.Sprintf("%.1fGB", float64(size)/float64(GB))
	case size >= MB:
		return fmt.Sprintf("%.1fMB", float64(size)/float64(MB))
	case size >= KB:
		return fmt.Sprintf("%.1fKB", float64(size)/float64(KB))
	default:
		return fmt.Sprintf("%dB", size)
	}
}
