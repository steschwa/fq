package utils

import (
	"os"
)

func IsStdinEmpty() bool {
	fi, err := os.Stdin.Stat()
	if err != nil {
		return true
	}
	if fi.Size() <= 0 {
		return true
	}

	return false
}
