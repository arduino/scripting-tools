package cli

import (
	"fmt"
	"io"
	"os"
	"strings"

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
		if err := runCmds(args); err != nil {
			if isUnknownCommandError(err) {
				printUsage(stderr)
			}
			fmt.Fprintf(stderr, "error: %v\n", err)
			return 1
		}
		return 0
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

func scriptByName(name string) (scripts.Script, bool) {
	for _, s := range availableScripts {
		if s.Name() == name {
			return s, true
		}
	}
	return nil, false
}

func isUnknownCommandError(err error) bool {
	return strings.HasPrefix(err.Error(), "unknown command ")
}

func runCmds(args []string) error {
	commands, err := splitCommands(args)
	if err != nil {
		return err
	}

	for _, cmdArgs := range commands {
		if err := runCmd(cmdArgs); err != nil {
			return err
		}
	}
	return nil
}

func splitCommands(args []string) ([][]string, error) {
	var commands [][]string
	start := 0

	for i, arg := range args {
		if arg != "::" {
			continue
		}
		if i == start {
			return nil, fmt.Errorf("unexpected separator \"::\"")
		}
		commands = append(commands, args[start:i])
		start = i + 1
	}

	if start == len(args) {
		return nil, fmt.Errorf("unexpected separator \"::\"")
	}
	commands = append(commands, args[start:])

	return commands, nil
}

func runCmd(args []string) error {
	script := args[0]
	s, ok := scriptByName(script)
	if !ok {
		return fmt.Errorf("unknown command %q", script)
	}
	if len(args) == 1 {
		fmt.Fprintf(os.Stderr, "Usage: scripting-tools script %s [flags]\n", s.Name())
		fmt.Fprintln(os.Stderr)
		fmt.Fprintln(os.Stderr, s.Help())
		fmt.Fprintln(os.Stderr)
		return fmt.Errorf("missing required flags for script %q", s.Name())
	}
	return s.Run(args[1:])
}
