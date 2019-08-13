package main

import (
	"github.com/songrgg/llsh/pkg/helper"
	"github.com/songrgg/llsh/pkg/llsh/cmd"
)

func main() {
	rootCmd := cmd.NewLLSHCommand()
	if err := rootCmd.Execute(); err != nil {
		helper.Errorlq(err)
	}
}
