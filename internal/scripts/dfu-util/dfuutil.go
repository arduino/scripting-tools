package dfuutil

import (
	"fmt"
	"os"
	"strings"

	"github.com/arduino/go-paths-helper"
)

type dfuUtilScript struct {
	bootloaderInstalled string
	bootloaderRequired  string
}

var Script = &dfuUtilScript{}

func (s *dfuUtilScript) Name() string {
	return "dfu-util"
}

func (s *dfuUtilScript) Run(args []string) error {
	for {
		if len(args) == 0 {
			return nil
		}
		switch args[0] {
		case "::bootloader-installed":
			if len(args) < 2 {
				return fmt.Errorf("expected installed bootloader version")
			}
			if args[1] == "" {
				return fmt.Errorf("installed bootloader version cannot be empty")
			}
			s.bootloaderInstalled = args[1]
			args = args[2:]
		case "::bootloader-required":
			if len(args) < 2 {
				return fmt.Errorf("expected required bootloader version")
			}
			if args[1] == "" {
				return fmt.Errorf("required bootloader version cannot be empty")
			}
			s.bootloaderRequired = args[1]
			args = args[2:]
		case "::upload-bootloader":
			var cmd []string
			cmd, args = constructCommand(args[1:])
			if s.bootloaderInstalled == "" {
				return fmt.Errorf("::bootloader-installed version is not set")
			}
			if s.bootloaderRequired == "" {
				return fmt.Errorf("::bootloader-required version is not set")
			}
			if s.bootloaderInstalled == s.bootloaderRequired {
				fmt.Println("Installed bootloader is up-to-date.")
				continue
			}
			if err := runCmd(cmd); err != nil {
				return fmt.Errorf("failed to upload bootloader: %w", err)
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
