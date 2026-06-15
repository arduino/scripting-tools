package cli

import (
	"fmt"
	"io"

	"github.com/arduino/arduino-upload-scripter/internal/scripts"
	dfuutil "github.com/arduino/arduino-upload-scripter/internal/scripts/dfu-util"
	"github.com/arduino/arduino-upload-scripter/internal/version"
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
	case "script":
		if err := runUploadScript(args[1:]); err != nil {
			fmt.Fprintf(stderr, "error: %v\n", err)
			return 1
		}
		return 0
	default:
		fmt.Fprintf(stderr, "error: unknown command %q\n\n", args[0])
		printUsage(stderr)
		return 1
	}
}

func printUsage(w io.Writer) {
	fmt.Fprintln(w, "arduino-upload-scripter")
	fmt.Fprintln(w)
	fmt.Fprintln(w, "Usage:")
	fmt.Fprintln(w, "  arduino-upload-scripter <command> [flags]")
	fmt.Fprintln(w)
	fmt.Fprintln(w, "Commands:")
	fmt.Fprintln(w, "  script     Run a script")
	fmt.Fprintln(w, "  version    Print version")
	fmt.Fprintln(w, "  help       Show help")
}

var availableScripts = []scripts.Script{
	dfuutil.Script,
}

func runUploadScript(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("missing required script name")
	}

	script := args[0]
	for _, s := range availableScripts {
		if s.Name() == script {
			return s.Run(args[1:])
		}
	}
	return fmt.Errorf("unknown script %q", script)
}
