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
	os.WriteFile(filePath, content, 0644)

	result, err := GetPathSize(filePath, false, false, false)
	require.NoError(t, err)
	expected := fmt.Sprintf("%dB\t%s", len(content), filePath)
	require.Equal(t, expected, result)
}

func TestGetPathSize_Directory(t *testing.T) {
	tempDir := t.TempDir()

	file1 := filepath.Join(tempDir, "file1.txt")
	file2 := filepath.Join(tempDir, "file2.txt")
	os.WriteFile(file1, []byte("content1"), 0644)
	os.WriteFile(file2, []byte("content2"), 0644)

	result, err := GetPathSize(tempDir, false, false, false)
	require.NoError(t, err)
	require.Contains(t, result, "16B")
	require.Contains(t, result, tempDir)
}

func TestGetPathSize_Recursive(t *testing.T) {
	tempDir := t.TempDir()

	subDir := filepath.Join(tempDir, "sub")
	os.Mkdir(subDir, 0755)

	file1 := filepath.Join(tempDir, "file1.txt")
	file2 := filepath.Join(subDir, "file2.txt")
	os.WriteFile(file1, []byte("12345"), 0644)
	os.WriteFile(file2, []byte("67890"), 0644)

	result1, err := GetPathSize(tempDir, false, false, false)
	require.NoError(t, err)
	require.Contains(t, result1, "5B")

	result2, err := GetPathSize(tempDir, true, false, false)
	require.NoError(t, err)
	require.Contains(t, result2, "10B")
}

func TestGetPathSize_HiddenFiles(t *testing.T) {
	tempDir := t.TempDir()

	normalFile := filepath.Join(tempDir, "normal.txt")
	hiddenFile := filepath.Join(tempDir, ".hidden.txt")
	os.WriteFile(normalFile, []byte("normal"), 0644)
	os.WriteFile(hiddenFile, []byte("hidden"), 0644)

	result1, err := GetPathSize(tempDir, false, false, false)
	require.NoError(t, err)
	require.Contains(t, result1, "6B")

	result2, err := GetPathSize(tempDir, false, false, true)
	require.NoError(t, err)
	require.Contains(t, result2, "12B")
}

func TestGetPathSize_HumanReadable(t *testing.T) {
	tempDir := t.TempDir()
	filePath := filepath.Join(tempDir, "test.txt")

	content := make([]byte, 1536)
	for i := range content {
		content[i] = 'A'
	}
	os.WriteFile(filePath, content, 0644)

	result1, err := GetPathSize(filePath, false, false, false)
	require.NoError(t, err)
	require.Contains(t, result1, "1536B")

	result2, err := GetPathSize(filePath, false, true, false)
	require.NoError(t, err)
	require.Contains(t, result2, "1.5KB")
}

func TestFormatSize(t *testing.T) {
	testCases := []struct {
		size     int64
		human    bool
		expected string
	}{
		{123, false, "123B"},
		{123, true, "123B"},
		{1024, true, "1KB"},
		{1536, true, "1.5KB"},
		{2048, true, "2KB"},
		{1234567, true, "1.2MB"},
		{1048576, true, "1MB"},
		{1073741824, true, "1GB"},
		{0, true, "0B"},
		{999, true, "999B"},
		{1000, true, "1000B"},
		{1023, true, "1023B"},
		{1024*1024 - 1, true, "1024KB"},
		{1024*1024*1024 - 1, true, "1024MB"},
	}

	for _, tc := range testCases {
		result := FormatSize(tc.size, tc.human)
		require.Equal(t, tc.expected, result,
			"For size %d and human=%v", tc.size, tc.human)
	}
}

func TestGetPathSize_ComplexStructure(t *testing.T) {
	tempDir := t.TempDir()

	os.WriteFile(filepath.Join(tempDir, "file1.txt"), make([]byte, 10), 0644)
	os.WriteFile(filepath.Join(tempDir, ".hidden.txt"), make([]byte, 5), 0644)

	subDir := filepath.Join(tempDir, "subdir")
	deepDir := filepath.Join(subDir, "deep")
	os.MkdirAll(deepDir, 0755)

	os.WriteFile(filepath.Join(subDir, "file2.txt"), make([]byte, 15), 0644)
	os.WriteFile(filepath.Join(subDir, ".hidden2.txt"), make([]byte, 7), 0644)
	os.WriteFile(filepath.Join(deepDir, "file3.txt"), make([]byte, 20), 0644)

	result1, err := GetPathSize(tempDir, false, false, false)
	require.NoError(t, err)
	require.Contains(t, result1, "10B")

	result2, err := GetPathSize(tempDir, true, false, false)
	require.NoError(t, err)
	require.Contains(t, result2, "45B")

	result3, err := GetPathSize(tempDir, true, false, true)
	require.NoError(t, err)
	require.Contains(t, result3, "57B")

	result4, err := GetPathSize(tempDir, true, true, true)
	require.NoError(t, err)
	require.Contains(t, result4, "57B")
}

func TestGetPathSize_NonExistentPath(t *testing.T) {
	_, err := GetPathSize("/non/existent/path", false, false, false)
	require.Error(t, err)
	require.Contains(t, err.Error(), "no such file")
}
