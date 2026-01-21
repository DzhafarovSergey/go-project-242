package code

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetPathSizeFile(t *testing.T) {
	tmpfile, err := os.CreateTemp("", "testfile")
	require.NoError(t, err)
	defer os.Remove(tmpfile.Name())

	content := []byte("Hello, World!")
	_, err = tmpfile.Write(content)
	require.NoError(t, err)
	err = tmpfile.Close()
	require.NoError(t, err)

	result, err := GetPathSize(tmpfile.Name(), false, false, false)
	require.NoError(t, err)

	expected := "13B\t" + tmpfile.Name()
	assert.Equal(t, expected, result)
}

func TestGetPathSizeDirectory(t *testing.T) {
	tmpdir, err := os.MkdirTemp("", "testdir")
	require.NoError(t, err)
	defer os.RemoveAll(tmpdir)

	file1 := filepath.Join(tmpdir, "file1.txt")
	file2 := filepath.Join(tmpdir, "file2.txt")

	err = os.WriteFile(file1, []byte("12345"), 0644)
	require.NoError(t, err)

	err = os.WriteFile(file2, []byte("67890"), 0644)
	require.NoError(t, err)

	result, err := GetPathSize(tmpdir, false, false, false)
	require.NoError(t, err)

	expected := "10B\t" + tmpdir
	assert.Equal(t, expected, result)
}

func TestGetPathSizeNonExistentPath(t *testing.T) {
	result, err := GetPathSize("/non/existent/path", false, false, false)

	assert.Error(t, err)
	assert.Empty(t, result)
}

func TestFormatSizeBytes(t *testing.T) {
	testCases := []struct {
		name     string
		size     int64
		human    bool
		expected string
	}{
		{"0 bytes, no human", 0, false, "0B"},
		{"100 bytes, no human", 100, false, "100B"},
		{"1023 bytes, no human", 1023, false, "1023B"},

		{"0 bytes, human", 0, true, "0B"},
		{"500 bytes, human", 500, true, "500B"},
		{"1023 bytes, human", 1023, true, "1023B"},
		{"1024 bytes, human", 1024, true, "1KB"},
		{"1536 bytes, human", 1536, true, "1.5KB"},
		{"1048576 bytes, human", 1048576, true, "1MB"},
		{"1073741824 bytes, human", 1073741824, true, "1GB"},
		{"1099511627776 bytes, human", 1099511627776, true, "1TB"},

		{"24MB exact", 25165824, true, "24MB"},
		{"24.0MB", 25165824, true, "24MB"},
		{"24.5MB", 25690112, true, "24.5MB"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := FormatSizeBytes(tc.size, tc.human)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestGetPathSizeWithHumanFlag(t *testing.T) {
	tmpfile, err := os.CreateTemp("", "testfile")
	require.NoError(t, err)
	defer os.Remove(tmpfile.Name())

	size := 24 * 1024 * 1024
	data := make([]byte, size)

	_, err = tmpfile.Write(data)
	require.NoError(t, err)
	err = tmpfile.Close()
	require.NoError(t, err)

	result, err := GetPathSize(tmpfile.Name(), false, false, false)
	require.NoError(t, err)
	expected := fmt.Sprintf("%dB\t%s", size, tmpfile.Name())
	assert.Equal(t, expected, result)

	result, err = GetPathSize(tmpfile.Name(), false, true, false)
	require.NoError(t, err)
	expected = fmt.Sprintf("24MB\t%s", tmpfile.Name())
	assert.Equal(t, expected, result)
}

func TestGetPathSizeSmallFileHuman(t *testing.T) {
	tmpfile, err := os.CreateTemp("", "testfile")
	require.NoError(t, err)
	defer os.Remove(tmpfile.Name())

	size := int64(25165824)
	data := make([]byte, size)

	_, err = tmpfile.Write(data)
	require.NoError(t, err)
	err = tmpfile.Close()
	require.NoError(t, err)

	result, err := GetPathSize(tmpfile.Name(), false, false, false)
	require.NoError(t, err)
	assert.Equal(t, fmt.Sprintf("%dB\t%s", size, tmpfile.Name()), result)

	result, err = GetPathSize(tmpfile.Name(), false, true, false)
	require.NoError(t, err)
	assert.Equal(t, fmt.Sprintf("24MB\t%s", tmpfile.Name()), result)
}

func TestGetPathSizeWithHiddenFiles(t *testing.T) {
	tmpdir, err := os.MkdirTemp("", "testdir-hidden")
	require.NoError(t, err)
	defer os.RemoveAll(tmpdir)

	file1 := filepath.Join(tmpdir, "file1.txt")
	file2 := filepath.Join(tmpdir, "file2.txt")

	hiddenFile1 := filepath.Join(tmpdir, ".hidden1.txt")
	hiddenFile2 := filepath.Join(tmpdir, ".hidden2.txt")

	hiddenDir := filepath.Join(tmpdir, ".hidden_dir")
	err = os.Mkdir(hiddenDir, 0755)
	require.NoError(t, err)

	fileInHiddenDir := filepath.Join(hiddenDir, "file.txt")

	err = os.WriteFile(file1, []byte("12345"), 0644)
	require.NoError(t, err)

	err = os.WriteFile(file2, []byte("67890"), 0644)
	require.NoError(t, err)

	err = os.WriteFile(hiddenFile1, []byte("hidden1"), 0644)
	require.NoError(t, err)

	err = os.WriteFile(hiddenFile2, []byte("hidden2"), 0644)
	require.NoError(t, err)

	err = os.WriteFile(fileInHiddenDir, []byte("inhidden"), 0644)
	require.NoError(t, err)

	result, err := GetPathSize(tmpdir, false, false, false)
	require.NoError(t, err)
	expected := "10B\t" + tmpdir
	assert.Equal(t, expected, result)

	result, err = GetPathSize(tmpdir, false, false, true)
	require.NoError(t, err)
	expected = "24B\t" + tmpdir
	assert.Equal(t, expected, result)
}

func TestGetPathSizeHiddenFileDirectly(t *testing.T) {
	tmpfile, err := os.CreateTemp("", ".hiddenfile")
	require.NoError(t, err)
	defer os.Remove(tmpfile.Name())

	content := []byte("Hidden content")
	_, err = tmpfile.Write(content)
	require.NoError(t, err)
	err = tmpfile.Close()
	require.NoError(t, err)

	result, err := GetPathSize(tmpfile.Name(), false, false, false)
	require.NoError(t, err)
	assert.Equal(t, "", result)

	result, err = GetPathSize(tmpfile.Name(), false, false, true)
	require.NoError(t, err)
	expected := fmt.Sprintf("%dB\t%s", len(content), tmpfile.Name())
	assert.Equal(t, expected, result)
}

func testIsHidden(name string) bool {
	base := filepath.Base(name)
	if base == "." || base == ".." {
		return false
	}
	return strings.HasPrefix(base, ".")
}

func TestIsHiddenFunction(t *testing.T) {
	testCases := []struct {
		name     string
		path     string
		expected bool
	}{
		{"hidden file", ".gitignore", true},
		{"hidden file with path", "/path/to/.env", true},
		{"hidden dir", ".git", true},
		{"hidden dir with path", "/home/user/.config", true},
		{"normal file", "file.txt", false},
		{"normal dir", "documents", false},
		{"file starting with dot in middle", "file.name.txt", false},
		{"current directory", ".", false},
		{"parent directory", "..", false},
		{"dot at start", ".hidden", true},
		{"multiple dots", "...", true},
		{"dotdot something", "..file", true},
		{"file with two dots start", "..config", true},
		{"only dot", ".", false},
		{"only dot dot", "..", false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := testIsHidden(tc.path)
			assert.Equal(t, tc.expected, result, "Path: %s", tc.path)
		})
	}
}

func TestGetPathSizeCombinedFlags(t *testing.T) {
	tmpdir, err := os.MkdirTemp("", "testdir-combined")
	require.NoError(t, err)
	defer os.RemoveAll(tmpdir)

	normalFile := filepath.Join(tmpdir, "normal.txt")
	hiddenFile := filepath.Join(tmpdir, ".hidden.txt")

	normalSize := int64(1536)
	hiddenSize := int64(512)

	err = os.WriteFile(normalFile, make([]byte, normalSize), 0644)
	require.NoError(t, err)

	err = os.WriteFile(hiddenFile, make([]byte, hiddenSize), 0644)
	require.NoError(t, err)

	result, err := GetPathSize(tmpdir, false, false, false)
	require.NoError(t, err)
	expected := fmt.Sprintf("%dB\t%s", normalSize, tmpdir)
	assert.Equal(t, expected, result)

	result, err = GetPathSize(tmpdir, false, true, false)
	require.NoError(t, err)
	expected = fmt.Sprintf("1.5KB\t%s", tmpdir)
	assert.Equal(t, expected, result)

	result, err = GetPathSize(tmpdir, false, false, true)
	require.NoError(t, err)
	expected = fmt.Sprintf("%dB\t%s", normalSize+hiddenSize, tmpdir)
	assert.Equal(t, expected, result)

	result, err = GetPathSize(tmpdir, false, true, true)
	require.NoError(t, err)
	expected = fmt.Sprintf("2KB\t%s", tmpdir)
	assert.Equal(t, expected, result)
}

func TestGetPathSizeRecursive(t *testing.T) {
	tmpdir, err := os.MkdirTemp("", "testdir-recursive")
	require.NoError(t, err)
	defer os.RemoveAll(tmpdir)

	subdir1 := filepath.Join(tmpdir, "subdir1")
	subsubdir := filepath.Join(subdir1, "subsubdir")
	subdir2 := filepath.Join(tmpdir, "subdir2")

	err = os.MkdirAll(subsubdir, 0755)
	require.NoError(t, err)
	err = os.Mkdir(subdir2, 0755)
	require.NoError(t, err)

	file1 := filepath.Join(tmpdir, "file1.txt")
	file2 := filepath.Join(subdir1, "file2.txt")
	file3 := filepath.Join(subsubdir, "file3.txt")
	file4 := filepath.Join(subdir2, "file4.txt")

	err = os.WriteFile(file1, []byte("12345"), 0644)
	require.NoError(t, err)
	err = os.WriteFile(file2, []byte("1234567890"), 0644)
	require.NoError(t, err)
	err = os.WriteFile(file3, []byte("123456789012345"), 0644)
	require.NoError(t, err)
	err = os.WriteFile(file4, []byte("12345678901234567890"), 0644)
	require.NoError(t, err)

	result, err := GetPathSize(tmpdir, false, false, false)
	require.NoError(t, err)
	assert.Equal(t, "5B\t"+tmpdir, result)

	result, err = GetPathSize(tmpdir, true, false, false)
	require.NoError(t, err)
	assert.Equal(t, "50B\t"+tmpdir, result)

	result, err = GetPathSize(tmpdir, true, true, false)
	require.NoError(t, err)
	assert.Equal(t, "50B\t"+tmpdir, result)

	largeFile := filepath.Join(tmpdir, "large.txt")
	largeContent := make([]byte, 2048)
	err = os.WriteFile(largeFile, largeContent, 0644)
	require.NoError(t, err)

	result, err = GetPathSize(tmpdir, true, true, false)
	require.NoError(t, err)
	assert.Equal(t, "2KB\t"+tmpdir, result)
}

func TestGetPathSizeRecursiveWithHidden(t *testing.T) {
	tmpdir, err := os.MkdirTemp("", "testdir-recursive-hidden")
	require.NoError(t, err)
	defer os.RemoveAll(tmpdir)

	hiddendir := filepath.Join(tmpdir, ".hiddendir")
	subdir := filepath.Join(tmpdir, "subdir")
	err = os.MkdirAll(hiddendir, 0755)
	require.NoError(t, err)
	err = os.Mkdir(subdir, 0755)
	require.NoError(t, err)

	normalFile := filepath.Join(tmpdir, "normal.txt")
	hiddenFile := filepath.Join(tmpdir, ".hidden.txt")
	fileInHiddenDir := filepath.Join(hiddendir, "file.txt")
	normalInSubdir := filepath.Join(subdir, "normal2.txt")
	hiddenInSubdir := filepath.Join(subdir, ".hidden2.txt")

	err = os.WriteFile(normalFile, []byte("12345"), 0644)
	require.NoError(t, err)
	err = os.WriteFile(hiddenFile, []byte("1234567890"), 0644)
	require.NoError(t, err)
	err = os.WriteFile(fileInHiddenDir, []byte("123456789012345"), 0644)
	require.NoError(t, err)
	err = os.WriteFile(normalInSubdir, []byte("12345678901234567890"), 0644)
	require.NoError(t, err)
	err = os.WriteFile(hiddenInSubdir, make([]byte, 25), 0644)
	require.NoError(t, err)

	result, err := GetPathSize(tmpdir, true, false, false)
	require.NoError(t, err)
	assert.Equal(t, "25B\t"+tmpdir, result)

	result, err = GetPathSize(tmpdir, true, false, true)
	require.NoError(t, err)
	assert.Equal(t, "75B\t"+tmpdir, result)

	result, err = GetPathSize(tmpdir, false, false, true)
	require.NoError(t, err)
	assert.Equal(t, "15B\t"+tmpdir, result)
}

func TestGetPathSizeRecursiveSkipHiddenDirs(t *testing.T) {
	tmpdir, err := os.MkdirTemp("", "testdir-skip-hidden")
	require.NoError(t, err)
	defer os.RemoveAll(tmpdir)

	hiddenDir := filepath.Join(tmpdir, ".hidden_dir")
	hiddenSubdir := filepath.Join(hiddenDir, "subdir")
	normalDir := filepath.Join(tmpdir, "normal_dir")

	err = os.MkdirAll(hiddenSubdir, 0755)
	require.NoError(t, err)
	err = os.Mkdir(normalDir, 0755)
	require.NoError(t, err)

	normalFile := filepath.Join(tmpdir, "normal.txt")
	file1 := filepath.Join(hiddenDir, "file1.txt")
	file2 := filepath.Join(hiddenSubdir, "file2.txt")
	file3 := filepath.Join(normalDir, "file3.txt")

	err = os.WriteFile(normalFile, []byte("12345"), 0644)
	require.NoError(t, err)
	err = os.WriteFile(file1, make([]byte, 10), 0644)
	require.NoError(t, err)
	err = os.WriteFile(file2, make([]byte, 15), 0644)
	require.NoError(t, err)
	err = os.WriteFile(file3, make([]byte, 20), 0644)
	require.NoError(t, err)

	result, err := GetPathSize(tmpdir, true, false, false)
	require.NoError(t, err)
	assert.Equal(t, "25B\t"+tmpdir, result)
}

func TestGetPathSizeEdgeCases(t *testing.T) {
	tmpdir, err := os.MkdirTemp("", "testdir-empty")
	require.NoError(t, err)
	defer os.RemoveAll(tmpdir)

	result, err := GetPathSize(tmpdir, true, false, false)
	require.NoError(t, err)
	assert.Equal(t, "0B\t"+tmpdir, result)

	hiddenFile := filepath.Join(tmpdir, ".hidden")
	err = os.WriteFile(hiddenFile, []byte("test"), 0644)
	require.NoError(t, err)

	result, err = GetPathSize(tmpdir, true, false, false)
	require.NoError(t, err)
	assert.Equal(t, "0B\t"+tmpdir, result)

	result, err = GetPathSize(tmpdir, true, false, true)
	require.NoError(t, err)
	assert.Equal(t, "4B\t"+tmpdir, result)
}
