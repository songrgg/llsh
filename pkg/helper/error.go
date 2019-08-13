package helper

import (
	"fmt"
	"os"
)

// Errorfq prints the message with format and exits with 1.
func Errorfq(format string, a ...interface{}) {
	fmt.Printf(format, a)
	os.Exit(1)
}

// Errorlq prints the message with new line and exits with 1.
func Errorlq(a ...interface{}) {
	fmt.Println(a)
	os.Exit(1)
}
