// Package twophase implements the two-phase upload script.
package twophase

import (
	"fmt"
	"os"
	"strings"

	"github.com/arduino/go-paths-helper"
)

type twoPhaseUploadScript struct {
	loaderInstalled string
	loaderRequired  string
}

// Script is the instance of the twoPhaseUploadScript that will be used by the CLI.
var Script = &twoPhaseUploadScript{}

func (s *twoPhaseUploadScript) Name() string {
	return "two-phase-upload"
}

func (s *twoPhaseUploadScript) Description() string {
	return "Upload firmware using two-phase upload"
}

func (s *twoPhaseUploadScript) Help() string {
	return "This script allows you to upload firmware using two-phase upload. It supports the following commands:\n\n" +
		"  ::loader-installed <version>      Set the installed loader version\n" +
		"  ::loader-required <version>       Set the required loader version\n" +
		"  ::upload-loader <cmd> <args...>   Upload the loader if the installed version is different from the required version\n" +
		"  ::upload <cmd> <args...>          Upload firmware using two-phase upload"
}

func (s *twoPhaseUploadScript) Run(args []string) error {
	for {
		if len(args) == 0 {
			return nil
		}
		switch args[0] {
		case "::loader-installed":
			if len(args) < 2 {
				return fmt.Errorf("expected installed loader version")
			}
			if args[1] == "" {
				return fmt.Errorf("installed loader version cannot be empty")
			}
			s.loaderInstalled = args[1]
			args = args[2:]
		case "::loader-required":
			if len(args) < 2 {
				return fmt.Errorf("expected required loader version")
			}
			if args[1] == "" {
				return fmt.Errorf("required loader version cannot be empty")
			}
			s.loaderRequired = args[1]
			args = args[2:]
		case "::upload-loader":
			var cmd []string
			cmd, args = constructCommand(args[1:])
			if s.loaderInstalled == "" {
				return fmt.Errorf("::loader-installed version is not set")
			}
			if s.loaderRequired == "" {
				return fmt.Errorf("::loader-required version is not set")
			}
			if s.loaderInstalled == s.loaderRequired {
				fmt.Println("Installed loader is up-to-date.")
				continue
			}
			if err := runCmd(cmd); err != nil {
				return fmt.Errorf("failed to upload loader: %w", err)
			}
		case "::upload":
			var cmd []string
			cmd, args = constructCommand(args[1:])
			if err := runCmd(cmd); err != nil {
				return fmt.Errorf("failed to upload firmware: %w", err)
			}
		default:
			return fmt.Errorf("unknown command %q", args[0])
		}
	}
}

// Pick arguments until the next "::" and return them as a command and arguments.
// It returns the remaining arguments after the command.
func constructCommand(args []string) ([]string, []string) {
	var cmdArgs []string
	for i, arg := range args {
		if strings.HasPrefix(arg, "::") {
			return cmdArgs, args[i:]
		}
		cmdArgs = append(cmdArgs, arg)
	}
	return cmdArgs, nil
}

func runCmd(cmd []string) error {
	fmt.Printf("Running command: %s\n", strings.Join(cmd, " "))
	if len(cmd) == 0 {
		return fmt.Errorf("command cannot be empty")
	}
	proc, err := paths.NewProcess(nil, cmd...)
	if err != nil {
		return fmt.Errorf("failed to create process: %w", err)
	}
	proc.RedirectStderrTo(os.Stderr)
	proc.RedirectStdoutTo(os.Stdout)
	if err := proc.Run(); err != nil {
		return fmt.Errorf("command failed: %w", err)
	}
	return nil
}
