package main

import (
	"os"
	"path/filepath"
	"testing"

	//nolint:depguard
	"github.com/stretchr/testify/require"
)

var inputFilePath = "testdata/input.txt"

const fileSize = 6617

func TestCopy(t *testing.T) {
	dir := os.TempDir()

	t.Run("empty paths", func(t *testing.T) {
		require.NotNil(t, Copy("", "", 0, 0))
	})
	t.Run("offset exceeds file size ", func(t *testing.T) {
		require.Equal(t, ErrOffsetExceedsFileSize, Copy(inputFilePath, "", 10000, 0))
	})
	t.Run("empty path output", func(t *testing.T) {
		require.NotNil(t, Copy(inputFilePath, "", 0, 0))
	})
	t.Run("full copy", func(t *testing.T) {
		outputFilePath := filepath.Join(dir, "testout1.txt")
		require.Nil(t, Copy(inputFilePath, outputFilePath, 0, 0))
		checkOutputSize(t, outputFilePath, fileSize)
	})
	t.Run("part copy", func(t *testing.T) {
		outputFilePath := filepath.Join(dir, "testout2.txt")
		require.Nil(t, Copy(inputFilePath, outputFilePath, 6000, 0))
		checkOutputSize(t, outputFilePath, fileSize-6000)
	})
	t.Run("part copy limited", func(t *testing.T) {
		outputFilePath := filepath.Join(dir, "testout3.txt")
		require.Nil(t, Copy(inputFilePath, outputFilePath, 6000, 100))
		checkOutputSize(t, outputFilePath, 100)
	})
	t.Run("part copy limited less because of the end of the file", func(t *testing.T) {
		outputFilePath := filepath.Join(dir, "testout3.txt")
		require.Nil(t, Copy(inputFilePath, outputFilePath, 6000, 1000))
		checkOutputSize(t, outputFilePath, fileSize-6000)
	})
}

func checkOutputSize(t *testing.T, outputFilePath string, expectedFileSize int64) {
	t.Helper()
	file, err := os.Open(outputFilePath)
	require.Nil(t, err)
	fi, err := file.Stat()
	require.Nil(t, err)
	require.Equal(t, expectedFileSize, fi.Size())
}
