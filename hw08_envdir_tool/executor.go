package main

import (
	"os"
	"os/exec"
)

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmd []string, env Environment) (returnCode int) {
	for keyEnv, envValue := range env {
		if envValue.NeedRemove {
			os.Unsetenv(keyEnv)
		} else {
			if _, ok := os.LookupEnv(keyEnv); ok {
				os.Unsetenv(keyEnv)
			}
			os.Setenv(keyEnv, envValue.Value)
		}
	}

	//nolint:gosec
	command := exec.Command(cmd[0], cmd[1:]...)

	command.Stdin = os.Stdin
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr

	if err := command.Run(); err != nil {
		println("Error in RunCmd", err)
	}

	return command.ProcessState.ExitCode()
}
