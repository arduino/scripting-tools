// This file is part of Arduino scripting-tools.
//
// SPDX-FileCopyrightText: Arduino s.r.l. and/or its affiliated companies
// SPDX-License-Identifier: GPL-3.0-or-later

// Package conditional implements the conditional script that allows running commands based on conditions.
package conditional

import (
	"fmt"
	"os"
	"strings"

	"github.com/arduino/go-paths-helper"
)

type conditionalScript struct {
}

// Script is the instance of the conditionalScript that will be used by the CLI.
var Script = &conditionalScript{}

func (s *conditionalScript) Name() string {
	return "if"
}

func (s *conditionalScript) Description() string {
	return "Run a command under a certain condition"
}

func (s *conditionalScript) Help() string {
	return "This script runs a command when a condition is true.\n\n" +
		"Usage:\n" +
		"  if <condition> <command...>\n\n" +
		"Conditions:\n" +
		"  eq <val1> <val2>        True if val1 equals val2\n" +
		"  not <cond>              True if cond is false\n" +
		"  and <cond1> <cond2>     True if both conditions are true\n" +
		"  or  <cond1> <cond2>     True if either condition is true"
}

func (s *conditionalScript) Run(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("expected condition")
	}
	result, cmd, err := parseCondition(args)
	if err != nil {
		return fmt.Errorf("invalid condition: %w", err)
	}
	if len(cmd) == 0 {
		return fmt.Errorf("expected command after condition")
	}
	if result {
		return runCmd(cmd)
	}
	return nil
}

// parseCondition parses a condition expression from the given tokens using a
// simple recursive descent parser. It returns the boolean result, the remaining
// unconsumed tokens, and any parse error.
//
// Grammar:
//
//	condition := "eq"  <value> <value>
//	           | "not" condition
//	           | "and" condition condition
//	           | "or"  condition condition
func parseCondition(tokens []string) (bool, []string, error) {
	if len(tokens) == 0 {
		return false, nil, fmt.Errorf("expected condition operator")
	}
	switch tokens[0] {
	case "eq":
		if len(tokens) < 3 {
			return false, nil, fmt.Errorf("\"eq\" requires two operands")
		}
		return tokens[1] == tokens[2], tokens[3:], nil
	case "not":
		val, rest, err := parseCondition(tokens[1:])
		if err != nil {
			return false, nil, fmt.Errorf("\"not\": %w", err)
		}
		return !val, rest, nil
	case "and":
		left, rest, err := parseCondition(tokens[1:])
		if err != nil {
			return false, nil, fmt.Errorf("\"and\" left operand: %w", err)
		}
		right, rest, err := parseCondition(rest)
		if err != nil {
			return false, nil, fmt.Errorf("\"and\" right operand: %w", err)
		}
		return left && right, rest, nil
	case "or":
		left, rest, err := parseCondition(tokens[1:])
		if err != nil {
			return false, nil, fmt.Errorf("\"or\" left operand: %w", err)
		}
		right, rest, err := parseCondition(rest)
		if err != nil {
			return false, nil, fmt.Errorf("\"or\" right operand: %w", err)
		}
		return left || right, rest, nil
	default:
		return false, nil, fmt.Errorf("unknown condition operator %q", tokens[0])
	}
}

func runCmd(cmd []string) error {
	if len(cmd) == 0 {
		return fmt.Errorf("command cannot be empty")
	}
	fmt.Printf("Running command: %s\n", strings.Join(cmd, " "))
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
