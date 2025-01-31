package main

import (
	"testing"

	//nolint:depguard
	"github.com/stretchr/testify/require"
)

func TestRunCmd(t *testing.T) {
	var exitcode int

	t.Run("echo test", func(t *testing.T) {
		exitcode = RunCmd([]string{"echo", "\"Hello, World\""}, Environment{})
		require.Equal(t, 0, exitcode)
	})

	t.Run("full test", func(t *testing.T) {
		env := Environment{}
		env["BAR"] = EnvValue{"bar", false}
		env["UNSET"] = EnvValue{"", true}
		exitcode = RunCmd([]string{"/bin/bash", "testdata/customCmd.sh", "5", "2"}, env)
		require.Equal(t, 5, exitcode)
	})
}
