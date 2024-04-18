package main

import (
	"dockerAddHost/pkg/runner"
	"golang.org/x/term"
	"syscall"
)

func main() {
	if term.IsTerminal(syscall.Stdin) {
		runner.Run()
	} else {
		runner.Command()
	}
}
