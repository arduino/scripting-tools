// Package main is the entry point for the scripting-tools CLI.
package main

import (
	"os"

	"github.com/arduino/scripting-tools/internal/cli"
)

func main() {
	os.Exit(cli.Run(os.Args[1:], os.Stdout, os.Stderr))
}
