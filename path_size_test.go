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

func TestFormatSize(t *testing.T) {
	testCases := []struct {
		size     int64
		human    bool
		expected string
	}{
		{123, false, "123B"},
		{123, true, "123B"},
		{1024, true, "1.0KB"},
		{1536, true, "1.5KB"},
		{2048, true, "2.0KB"},
		{1234567, true, "1.2MB"},
		{1048576, true, "1.0MB"},
		{1073741824, true, "1.0GB"},
		{0, true, "0B"},
		{999, true, "999B"},
		{1000, true, "1000B"},
		{1023, true, "1023B"},
		{1024*1024 - 1, true, "1024KB"},
		{1024*1024*1024 - 1, true, "1024MB"},
		{1500000, true, "1.4MB"},
		{999999, true, "977KB"},
		{1500, true, "1.5KB"},
	}

	for _, tc := range testCases {
		result := FormatSize(tc.size, tc.human)
		require.Equal(t, tc.expected, result,
			"For size %d and human=%v", tc.size, tc.human)
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
	require.Contains(t, err.Error(), "no such file")
}

func TestGetPathSizeWithPath(t *testing.T) {
	tempDir := t.TempDir()
	filePath := filepath.Join(tempDir, "test.txt")
	content := []byte("test")
	err := os.WriteFile(filePath, content, 0644)
	require.NoError(t, err)

	result, err := GetPathSizeWithPath(filePath, false, false, false)
	require.NoError(t, err)
	expected := fmt.Sprintf("%dB\t%s", len(content), filePath)
	require.Equal(t, expected, result)
}
