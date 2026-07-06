package cli

import (
	"fmt"
	"io"
	"os"

	"github.com/arduino/scripting-tools/internal/scripts"
	"github.com/arduino/scripting-tools/internal/scripts/conditional"
	twophase "github.com/arduino/scripting-tools/internal/scripts/two-phase"
	"github.com/arduino/scripting-tools/internal/version"
)

func Run(args []string, stdout, stderr io.Writer) int {
	if len(args) == 0 {
		printUsage(stdout)
		return 0
	}

	switch args[0] {
	case "-h", "--help", "help":
		printUsage(stdout)
		return 0
	case "version", "--version", "-v":
		fmt.Fprintln(stdout, version.Value)
		return 0
	default:
		for _, s := range availableScripts {
			if s.Name() == args[0] {
				if err := runCmd(args); err != nil {
					fmt.Fprintf(stderr, "error: %v\n", err)
					return 1
				}
				return 0
			}
		}
		printUsage(stderr)
		fmt.Fprintf(stderr, "error: unknown command %q\n", args[0])
		return 1
	}
}

func printUsage(w io.Writer) {
	fmt.Fprintln(w, "scripting-tools")
	fmt.Fprintln(w)
	fmt.Fprintln(w, "Usage:")
	fmt.Fprintln(w, "  scripting-tools <command> [flags]")
	fmt.Fprintln(w)
	fmt.Fprintln(w, "Commands:")

	for _, s := range availableScripts {
		fmt.Fprintf(w, "  %-15s %s\n", s.Name(), s.Description())
	}
	fmt.Fprintln(w, "  version         Print version")
	fmt.Fprintln(w, "  help            Show help")
	fmt.Fprintln(w)
}

var availableScripts = []scripts.Script{
	twophase.Script,
	conditional.Script,
}

func runCmd(args []string) error {
	script := args[0]
	for _, s := range availableScripts {
		if s.Name() == script {
			if len(args) == 1 {
				fmt.Fprintf(os.Stderr, "Usage: scripting-tools script %s [flags]\n", s.Name())
				fmt.Fprintln(os.Stderr)
				fmt.Fprintln(os.Stderr, s.Help())
				fmt.Fprintln(os.Stderr)
				return fmt.Errorf("missing required flags for script %q", s.Name())
			}
			return s.Run(args[1:])
		}
	}
	return fmt.Errorf("unknown script %q", script)
}
