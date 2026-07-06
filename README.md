# scripting-tools

This tool is a wrapper for upload tools used in Arduino Platforms.
It helps to perform some simple action before starting the actual upload.

## Usage

The `scripting-tools` runs recipes hardcoded in the tool itself, called "scripts".

Each script may have a number of options or commands available that can be called through the command line arguments.
Show available commands:

```
$ scripting-tools

Usage:
  scripting-tools <command> [flags]

Commands:
  two-phase-upload Upload firmware using two-phase upload
  if              Run a command under a certain condition
  version         Print version
  help            Show help
```

Run the conditional script when a condition is true:

```
$ scripting-tools if eq a a echo OK
Running command: echo OK
OK
```

Run two-phase upload and skip loader upload when versions match:

```
$ scripting-tools two-phase-upload \
    ::loader-installed 1.2.3 \
    ::loader-required 1.2.3 \
    ::upload-loader echo LOADER \
    ::upload echo FIRMWARE
Installed loader is up-to-date.
Running command: echo FIRMWARE
FIRMWARE
```

Unknown command example:

```
$ scripting-tools does-not-exist
scripting-tools

Usage:
  scripting-tools <command> [flags]

Commands:
  two-phase-upload Upload firmware using two-phase upload
  if              Run a command under a certain condition
  version         Print version
  help            Show help

error: unknown command "does-not-exist"
exit status 1
```

## Build

```bash
go build -ldflags "-X github.com/arduino/scripting-tools/internal/version.Value=v1.0.0"
```
