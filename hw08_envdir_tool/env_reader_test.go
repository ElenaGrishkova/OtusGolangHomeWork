package main

import (
	"errors"
	"testing"

	//nolint:depguard
	"github.com/stretchr/testify/require"
)

func TestReadDir(t *testing.T) {
	expected := Environment{
		"BAR": EnvValue{
			Value:      "bar",
			NeedRemove: false,
		},
		"EMPTY": EnvValue{
			Value:      "",
			NeedRemove: false,
		},
		"FOO": EnvValue{
			Value:      "   foo\nwith new line",
			NeedRemove: false,
		},
		"HELLO": EnvValue{
			Value:      "\"hello\"",
			NeedRemove: false,
		},
		"UNSET": EnvValue{
			NeedRemove: true,
		},
	}

	t.Run("full test", func(t *testing.T) {
		result, err := ReadDir("testdata/env")
		require.Equal(t, expected, result)
		require.NoError(t, err)
	})

	t.Run("missing directory argument", func(t *testing.T) {
		_, err := ReadDir("")
		require.Truef(t, errors.Is(err, ErrMissingDirectory), "actual err - %v", err)
	})

	t.Run("directory is not exists", func(t *testing.T) {
		nonexistentdir := "testdata/nonexistentdirectory"
		require.NoDirExists(t, nonexistentdir)
		_, err := ReadDir(nonexistentdir)
		require.NotNil(t, err)
	})
}
