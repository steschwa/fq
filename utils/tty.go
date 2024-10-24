package utils

import "fmt"

func ClearLine() {
	fmt.Print("\033[2K\r")
}
