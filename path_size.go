package code

import (
	"fmt"
	"math"
	"os"
	"path/filepath"
)

func GetPathSize(path string, recursive, human, all bool) (string, error) {
	size, err := calculateSize(path, recursive, all)
	if err != nil {
		return "", err
	}

	if human {
		return FormatSize(size, true), nil
	}
	return fmt.Sprintf("%dB", size), nil
}

func GetPathSizeWithPath(path string, recursive, human, all bool) (string, error) {
	sizeStr, err := GetPathSize(path, recursive, human, all)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s\t%s", sizeStr, path), nil
}

func calculateSize(path string, recursive, all bool) (int64, error) {
	info, err := os.Lstat(path)
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

		if !all && len(entryName) > 0 && entryName[0] == '.' {
			continue
		}

		fullPath := filepath.Join(path, entryName)
		entryInfo, err := os.Lstat(fullPath)
		if err != nil {
			continue
		}

		if entryInfo.IsDir() && recursive {
			subDirSize, err := calculateSize(fullPath, recursive, all)
			if err == nil {
				totalSize += subDirSize
			}
		} else if !entryInfo.IsDir() {
			totalSize += entryInfo.Size()
		}
	}

	return totalSize, nil
}

func FormatSize(size int64, human bool) string {
	if !human {
		return fmt.Sprintf("%dB", size)
	}

	units := []string{"B", "KB", "MB", "GB", "TB", "PB", "EB"}
	if size == 0 {
		return "0B"
	}

	base := 1024.0
	exp := int(math.Log(float64(size)) / math.Log(base))
	if exp >= len(units) {
		exp = len(units) - 1
	}

	value := float64(size) / math.Pow(base, float64(exp))

	if exp == 0 {
		return fmt.Sprintf("%d%s", size, units[exp])
	}
	if value == math.Trunc(value) {
		return fmt.Sprintf("%.1f%s", value, units[exp])
	}
	if value < 10 {
		return fmt.Sprintf("%.1f%s", value, units[exp])
	}
	return fmt.Sprintf("%.0f%s", math.Round(value), units[exp])
}
