package main

import (
	"fmt"
	"os"
)

func main() {
	args := os.Args
	l := len(args)
	if l < 2 {
		_, _ = fmt.Fprint(os.Stderr, "error: argument count mismatch")
		os.Exit(1)
	}
	env, err := ReadDir(args[1])
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "error: %s", err)
		os.Exit(1)
	}
	if l < 3 {
		os.Exit(0)
	}
	os.Exit(RunCmd(args[2:], env))
}
