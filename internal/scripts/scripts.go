// This file is part of Arduino scripting-tools.
//
// SPDX-FileCopyrightText: Arduino s.r.l. and/or its affiliated companies
// SPDX-License-Identifier: GPL-3.0-or-later

// Package scripts defines the interface for scripts that can be run by the scripting-tools CLI.
package scripts

// Script defines the interface that all scripts must implement.
type Script interface {
	// Name returns the name of the script, e.g. "dfu-util".
	Name() string

	// Description returns a short description of the script.
	Description() string

	// Help returns a detailed help message for the script.
	Help() string

	// Run executes the script with the given arguments.
	Run(args []string) error
}
