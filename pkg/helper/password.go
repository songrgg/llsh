package helper

import (
	"fmt"
	"os"

	"golang.org/x/crypto/ssh/terminal"
)

func ReadPassword(prompt string) (string, error) {
	fmt.Fprintf(os.Stderr, prompt)

	tty, err := os.Open("/dev/tty")
	if err != nil {
		return "", fmt.Errorf(`TTY is required for reading password, but /dev/tty can't be opened: %v`, err)
	}

	password, err := terminal.ReadPassword(int(tty.Fd()))
	if err != nil {
		return "", fmt.Errorf(`can't read password: %v`, err)
	}

	if prompt != "" {
		fmt.Fprintln(os.Stderr)
	}

	return string(password), nil
}
