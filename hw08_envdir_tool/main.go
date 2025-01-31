package main

import (
	"errors"
	"fmt"
	"os"
)

var (
	ErrMissingDirectory  = errors.New("missing argument: directory")
	ErrDirectoryNotExist = errors.New("directory does not exist")
)

func main() {
	args := os.Args
	if len(args) < 2 {
		fmt.Println(ErrMissingDirectory)
		os.Exit(1)
	}

	dir := args[1]
	if stat, err := os.Stat(dir); err != nil || !stat.IsDir() {
		fmt.Println(ErrDirectoryNotExist)
		os.Exit(1)
	}

	env, err := ReadDir(dir)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	os.Exit(RunCmd(args[2:], env))
}
