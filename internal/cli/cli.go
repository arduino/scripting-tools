// This file is part of Arduino scripting-tools.
//
// SPDX-FileCopyrightText: Arduino s.r.l. and/or its affiliated companies
// SPDX-License-Identifier: GPL-3.0-or-later

// Package cli implements the command-line interface for the scripting-tools CLI.
package cli

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/arduino/scripting-tools/internal/scripts"
	"github.com/arduino/scripting-tools/internal/scripts/conditional"
	runscript "github.com/arduino/scripting-tools/internal/scripts/run"
	twophase "github.com/arduino/scripting-tools/internal/scripts/two-phase"
	"github.com/arduino/scripting-tools/internal/version"
)

// Run is the entry point for the scripting-tools CLI. It takes the command-line arguments and executes the appropriate script.
func Run(args []string, stdout, stderr io.Writer) int {
	separator, args, err := parseGlobalArgs(args)
	if err != nil {
		_, _ = fmt.Fprintf(stderr, "error: %v\n", err)
		return 1
	}

	if len(args) == 0 {
		printUsage(stdout)
		return 0
	}

	switch args[0] {
	case "-h", "--help", "help":
		printUsage(stdout)
		return 0
	case "version", "--version", "-v":
		_, _ = fmt.Fprintln(stdout, version.Version, "commit", version.Commit, "timestamp", version.Timestamp)
		return 0
	default:
		if err := runCmds(args, separator); err != nil {
			if isUnknownCommandError(err) {
				printUsage(stderr)
			}
			_, _ = fmt.Fprintf(stderr, "error: %v\n", err)
			return 1
		}
		return 0
	}
}

func parseGlobalArgs(args []string) (string, []string, error) {
	separator := "::"

	if len(args) > 0 && args[0] == "--sep" {
		if len(args) < 2 {
			return "", nil, fmt.Errorf("missing value for --sep")
		}
		if args[1] == "" {
			return "", nil, fmt.Errorf("separator cannot be empty")
		}
		separator = args[1]
		args = args[2:]
	}

	return separator, args, nil
}

func printUsage(w io.Writer) {
	_, _ = fmt.Fprintln(w, "scripting-tools")
	_, _ = fmt.Fprintln(w)
	_, _ = fmt.Fprintln(w, "Usage:")
	_, _ = fmt.Fprintln(w, "  scripting-tools [--sep <token>] <command> [flags]")
	_, _ = fmt.Fprintln(w, "  scripting-tools [--sep <token>] <command> [flags] [:: <command> [flags] ...]")
	_, _ = fmt.Fprintln(w)
	_, _ = fmt.Fprintln(w, "Options:")
	_, _ = fmt.Fprintln(w, "  --sep <token>   Top-level command separator token (default ::)")
	_, _ = fmt.Fprintln(w)
	_, _ = fmt.Fprintln(w, "Commands:")

	for _, s := range availableScripts {
		_, _ = fmt.Fprintf(w, "  %-15s %s\n", s.Name(), s.Description())
	}
	_, _ = fmt.Fprintln(w, "  version         Print version")
	_, _ = fmt.Fprintln(w, "  help            Show help")
	_, _ = fmt.Fprintln(w)
}

var availableScripts = []scripts.Script{
	twophase.Script,
	conditional.Script,
	runscript.Script,
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

func runCmds(args []string, separator string) error {
	commands, err := splitCommands(args, separator)
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

func splitCommands(args []string, separator string) ([][]string, error) {
	var commands [][]string
	start := 0

	for i, arg := range args {
		if arg != separator {
			continue
		}
		if i == start {
			return nil, fmt.Errorf("unexpected separator %q", separator)
		}
		commands = append(commands, args[start:i])
		start = i + 1
	}

	if start == len(args) {
		return nil, fmt.Errorf("unexpected separator %q", separator)
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
