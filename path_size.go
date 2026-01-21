// path_size.go
package code

import (
	"fmt"
	"math"
	"os"
	"path/filepath"
	"strings"
)

func GetPathSize(path string, recursive, human, all bool) (string, error) {
	return GetSize(path, recursive, human, all)
}

func GetSize(path string, recursive, human, all bool) (string, error) {
	info, err := os.Lstat(path)
	if err != nil {
		return "", err
	}

	var size int64
	if info.IsDir() {
		size, err = getDirSize(path, recursive, all)
		if err != nil {
			return "", err
		}
	} else {
		if !all && isHidden(path) {
			return "", nil
		}
		size = info.Size()
	}

	if !human {
		return fmt.Sprintf("%dB", size), nil
	}
	return FormatSizeBytes(size, true), nil
}

func getDirSize(path string, recursive, all bool) (int64, error) {
	var totalSize int64

	if !recursive {
		entries, err := os.ReadDir(path)
		if err != nil {
			return 0, err
		}

		for _, entry := range entries {
			if !all && isHidden(entry.Name()) {
				continue
			}

			info, err := entry.Info()
			if err != nil {
				return 0, err
			}
			if !info.IsDir() {
				totalSize += info.Size()
			}
		}
	} else {
		err := filepath.Walk(path, func(filePath string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if filePath == path {
				return nil
			}

			relPath, err := filepath.Rel(path, filePath)
			if err != nil {
				return err
			}

			parts := strings.Split(relPath, string(filepath.Separator))
			shouldSkip := false
			for _, part := range parts {
				if !all && isHidden(part) {
					shouldSkip = true
					break
				}
			}

			if shouldSkip {
				if info.IsDir() {
					return filepath.SkipDir
				}
				return nil
			}
			if !info.IsDir() {
				totalSize += info.Size()
			}

			return nil
		})
		if err != nil {
			return 0, err
		}
	}

	return totalSize, nil
}

func isHidden(name string) bool {
	base := filepath.Base(name)
	if base == "." || base == ".." {
		return false
	}
	return strings.HasPrefix(base, ".")
}

func FormatSize(size int64, path string, human bool) string {
	if !human {
		return fmt.Sprintf("%dB\t%s", size, path)
	}

	return fmt.Sprintf("%s\t%s", FormatSizeBytes(size, true), path)
}

func FormatSizeBytes(size int64, human bool) string {
	if !human {
		return fmt.Sprintf("%dB", size)
	}

	if size < 1024 {
		return fmt.Sprintf("%dB", size)
	}

	units := []string{"KB", "MB", "GB", "TB", "PB", "EB"}

	exp := int(math.Log(float64(size)) / math.Log(1024))
	if exp > len(units) {
		exp = len(units)
	}
	if exp > 0 {
		exp = exp - 1
	}

	value := float64(size) / math.Pow(1024, float64(exp+1))
	return fmt.Sprintf("%.1f%s", value, units[exp])
}
