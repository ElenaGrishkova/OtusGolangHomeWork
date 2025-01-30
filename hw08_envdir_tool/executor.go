package main

import (
	"os"
	"os/exec"
	"strings"
)

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmd []string, env Environment) (returnCode int) {
	envOS := os.Environ() // слайс строк "key=value"

	cmdEnv := make([]string, 0)
	for _, param := range envOS {
		keyOS := strings.Split(param, "=")[0]
		envValue, ok := env[keyOS]
		if ok {
			if !envValue.NeedRemove {
				cmdEnv = append(cmdEnv, strings.Join([]string{keyOS, envValue.Value}, "="))
				delete(env, keyOS)
			}
		} else {
			cmdEnv = append(cmdEnv, param)
		}
	}
	for keyEnv, envValue := range env {
		if !envValue.NeedRemove {
			cmdEnv = append(cmdEnv, strings.Join([]string{keyEnv, envValue.Value}, "="))
		}
	}
	cmdEnv = append(cmdEnv, cmd[1:]...)

	//nolint:gosec
	command := exec.Command(cmd[0], cmdEnv...)

	command.Stdin = os.Stdin
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr

	if err := command.Run(); err != nil {
		println("Error in RunCmd", err)
	}

	return command.ProcessState.ExitCode()
}
