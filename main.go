// This file is part of Arduino scripting-tools.
//
// SPDX-FileCopyrightText: Arduino s.r.l. and/or its affiliated companies
// SPDX-License-Identifier: GPL-3.0-or-later

// Package main is the entry point for the scripting-tools CLI.
package main

import (
	"os"

	"github.com/arduino/scripting-tools/internal/cli"
)

func main() {
	os.Exit(cli.Run(os.Args[1:], os.Stdout, os.Stderr))
}
