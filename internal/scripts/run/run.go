// This file is part of Arduino scripting-tools.
//
// SPDX-FileCopyrightText: Arduino s.r.l. and/or its affiliated companies
// SPDX-License-Identifier: GPL-3.0-or-later

// Package run implements the run script that allows running arbitrary commands.
package run

import (
	"fmt"
	"os"
	"strings"

	"github.com/arduino/go-paths-helper"
)

type runScript struct {
}

// Script is the instance of the runScript that will be used by the CLI.
var Script = &runScript{}

func (s *runScript) Name() string {
	return "run"
}

func (s *runScript) Description() string {
	return "Run a command"
}

func (s *runScript) Help() string {
	return "This script runs the given command.\n\n" +
		"Usage:\n" +
		"  run <command...>"
}

func (s *runScript) Run(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("expected command")
	}
	fmt.Printf("Running command: %s\n", strings.Join(args, " "))
	proc, err := paths.NewProcess(nil, args...)
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
