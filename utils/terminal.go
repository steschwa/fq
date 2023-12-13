package utils

import (
	"os"

	"golang.org/x/term"
)

func IsInteractiveTTY() bool {
	return term.IsTerminal(int(os.Stdout.Fd()))
}
