package code

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetPathSize_File(t *testing.T) {
	tempDir := t.TempDir()
	filePath := filepath.Join(tempDir, "test.txt")
	content := []byte("test content")
	err := os.WriteFile(filePath, content, 0644)
	require.NoError(t, err)

	result, err := GetPathSize(filePath, false, false, false)
	require.NoError(t, err)
	expected := fmt.Sprintf("%dB", len(content))
	require.Equal(t, expected, result)
}

func TestGetPathSize_Directory(t *testing.T) {
	tempDir := t.TempDir()

	file1 := filepath.Join(tempDir, "file1.txt")
	file2 := filepath.Join(tempDir, "file2.txt")
	err := os.WriteFile(file1, []byte("content1"), 0644)
	require.NoError(t, err)
	err = os.WriteFile(file2, []byte("content2"), 0644)
	require.NoError(t, err)

	result, err := GetPathSize(tempDir, false, false, false)
	require.NoError(t, err)
	require.Equal(t, "16B", result)
}

func TestGetPathSize_Recursive(t *testing.T) {
	tempDir := t.TempDir()

	subDir := filepath.Join(tempDir, "sub")
	err := os.Mkdir(subDir, 0755)
	require.NoError(t, err)

	file1 := filepath.Join(tempDir, "file1.txt")
	file2 := filepath.Join(subDir, "file2.txt")
	err = os.WriteFile(file1, []byte("12345"), 0644)
	require.NoError(t, err)
	err = os.WriteFile(file2, []byte("67890"), 0644)
	require.NoError(t, err)

	result1, err := GetPathSize(tempDir, false, false, false)
	require.NoError(t, err)
	require.Equal(t, "5B", result1)

	result2, err := GetPathSize(tempDir, true, false, false)
	require.NoError(t, err)
	require.Equal(t, "10B", result2)
}

func TestGetPathSize_HiddenFiles(t *testing.T) {
	tempDir := t.TempDir()

	normalFile := filepath.Join(tempDir, "normal.txt")
	hiddenFile := filepath.Join(tempDir, ".hidden.txt")
	err := os.WriteFile(normalFile, []byte("normal"), 0644)
	require.NoError(t, err)
	err = os.WriteFile(hiddenFile, []byte("hidden"), 0644)
	require.NoError(t, err)

	result1, err := GetPathSize(tempDir, false, false, false)
	require.NoError(t, err)
	require.Equal(t, "6B", result1)

	result2, err := GetPathSize(tempDir, false, false, true)
	require.NoError(t, err)
	require.Equal(t, "12B", result2)
}

func TestGetPathSize_HumanReadable(t *testing.T) {
	tempDir := t.TempDir()
	filePath := filepath.Join(tempDir, "test.txt")

	content := make([]byte, 1536)
	for i := range content {
		content[i] = 'A'
	}
	err := os.WriteFile(filePath, content, 0644)
	require.NoError(t, err)

	result1, err := GetPathSize(filePath, false, false, false)
	require.NoError(t, err)
	require.Equal(t, "1536B", result1)

	result2, err := GetPathSize(filePath, false, true, false)
	require.NoError(t, err)
	require.Equal(t, "1.5KB", result2)
}

func TestGetPathSize_HumanReadable_VariousSizes(t *testing.T) {
	testCases := []struct {
		name     string
		size     int64
		expected string
	}{
		{"0 bytes", 0, "0B"},
		{"100 bytes", 100, "100B"},
		{"999 bytes", 999, "999B"},
		{"1023 bytes", 1023, "1023B"},
		{"1KB", 1024, "1.0KB"},
		{"1.5KB", 1536, "1.5KB"},
		{"2KB", 2048, "2.0KB"},
		{"1MB", 1048576, "1.0MB"},
		{"1.2MB", 1258291, "1.2MB"}, // 1.2 * 1024 * 1024
		{"1GB", 1073741824, "1.0GB"},
		{"1.5GB", 1610612736, "1.5GB"}, // 1.5 * 1024 * 1024 * 1024
		{"999KB", 1022976, "999.0KB"},  // 999 * 1024
		{"1024KB (edge)", 1048575, "1024.0KB"},
		{"1024MB (edge)", 1073741823, "1024.0MB"},
	}

	tempDir := t.TempDir()

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			filePath := filepath.Join(tempDir, "test_"+tc.name+".txt")
			content := make([]byte, tc.size)
			err := os.WriteFile(filePath, content, 0644)
			require.NoError(t, err)

			result, err := GetPathSize(filePath, false, true, false)
			require.NoError(t, err)
			require.Equal(t, tc.expected, result, "For size %d", tc.size)
		})
	}
}

func TestGetPathSize_ComplexStructure(t *testing.T) {
	tempDir := t.TempDir()

	err := os.WriteFile(filepath.Join(tempDir, "file1.txt"), make([]byte, 10), 0644)
	require.NoError(t, err)
	err = os.WriteFile(filepath.Join(tempDir, ".hidden.txt"), make([]byte, 5), 0644)
	require.NoError(t, err)

	subDir := filepath.Join(tempDir, "subdir")
	deepDir := filepath.Join(subDir, "deep")
	err = os.MkdirAll(deepDir, 0755)
	require.NoError(t, err)

	err = os.WriteFile(filepath.Join(subDir, "file2.txt"), make([]byte, 15), 0644)
	require.NoError(t, err)
	err = os.WriteFile(filepath.Join(subDir, ".hidden2.txt"), make([]byte, 7), 0644)
	require.NoError(t, err)
	err = os.WriteFile(filepath.Join(deepDir, "file3.txt"), make([]byte, 20), 0644)
	require.NoError(t, err)

	result1, err := GetPathSize(tempDir, false, false, false)
	require.NoError(t, err)
	require.Equal(t, "10B", result1)

	result2, err := GetPathSize(tempDir, true, false, false)
	require.NoError(t, err)
	require.Equal(t, "45B", result2)

	result3, err := GetPathSize(tempDir, true, false, true)
	require.NoError(t, err)
	require.Equal(t, "57B", result3)

	result4, err := GetPathSize(tempDir, true, true, true)
	require.NoError(t, err)
	require.Equal(t, "57B", result4)
}

func TestGetPathSize_NonExistentPath(t *testing.T) {
	_, err := GetPathSize("/non/existent/path", false, false, false)
	require.Error(t, err)
}

func TestGetPathSize_Symlink(t *testing.T) {
	tempDir := t.TempDir()

	targetFile := filepath.Join(tempDir, "target.txt")
	content := []byte("test content")
	err := os.WriteFile(targetFile, content, 0644)
	require.NoError(t, err)

	linkFile := filepath.Join(tempDir, "link.txt")
	err = os.Symlink(targetFile, linkFile)
	require.NoError(t, err)

	result, err := GetPathSize(linkFile, false, false, false)
	require.NoError(t, err)
	expected := fmt.Sprintf("%dB", len(content))
	require.Equal(t, expected, result)
}
