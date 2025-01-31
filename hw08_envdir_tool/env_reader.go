package main

import (
	"bufio"
	"log"
	"os"
	"path"
	"strings"
)

type Environment map[string]EnvValue

// EnvValue helps to distinguish between empty files and files with the first empty line.
type EnvValue struct {
	Value      string
	NeedRemove bool
}

// ReadDir reads a specified directory and returns map of env variables.
// Variables represented as files where filename is name of variable, file first line is a value.
func ReadDir(dir string) (Environment, error) {
	if dir == "" {
		return nil, ErrMissingDirectory
	}
	dirEntries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	env := make(Environment)
	for _, entr := range dirEntries {
		if entr.IsDir() {
			continue
		}
		entrName := entr.Name()
		if strings.Contains(entrName, "=") {
			// имя не должно содержать =
			continue
		}

		envVal := EnvValue{}
		entrFile, err := os.Open(path.Join(dir, entrName))
		if err != nil {
			return nil, err
		}
		defer closeFile(entrFile)

		scanner := bufio.NewScanner(entrFile)
		if scanner.Scan() {
			envVal.Value = strings.ReplaceAll(
				strings.TrimRightFunc(
					scanner.Text(), func(r rune) bool {
						return r == ' ' || r == '\t'
					}),
				"\x00", "\n")
		} else {
			envVal.NeedRemove = true
		}
		env[entrName] = envVal
	}

	return env, nil
}

func closeFile(fromFile *os.File) {
	err := fromFile.Close()
	if err != nil {
		log.Panicf("failed to close source file: %v", err)
	}
}
