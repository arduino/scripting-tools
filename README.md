# scripting-tools

This tool is a wrapper for upload tools used in Arduino Platforms.
It helps to perform some simple action before starting the actual upload.

## Usage

The `scripting-tools` runs recipes hardcoded in the tool itself, called "scripts".

Each script may have a number of options or commands available that can be called through the command line arguments.

To get help for a particular script just add the script name to the command line:

```
$ scripting-tools dfu-util
Usage: scripting-tools dfu-util [flags]

This script allows you to upload firmware using dfu-util. It supports the following commands:

  ::bootloader-installed <version>      Set the installed bootloader version
  ::bootloader-required <version>       Set the required bootloader version
  ::upload-bootloader <dfu-util args>   Upload the bootloader if the installed version is different from the required version
  ::upload <dfu-util args>              Upload firmware using dfu-util

error: missing required flags for script "dfu-util"
exit status 1
```

To run a particular script, just add the needed command line flags to the command:
```
$ scripting-tools dfu-util \
    ::bootloader-installed 1.2.3 \
    ::bootloader-required 3.4.5 \
    ::upload-bootloader echo Uploding boot \
    ::upload echo UPLOAD SKETCH
Running command: echo Uploding boot
Uploding boot
Running command: echo UPLOAD SKETCH
UPLOAD SKETCH
$
```

```
$ scripting-tools dfu-util \
    ::bootloader-installed 1.2.3 \
    ::bootloader-required 1.2.3 \
    ::upload-bootloader echo Uploding boot \
    ::upload echo UPLOAD SKETCH
Installed bootloader is up-to-date.
Running command: echo UPLOAD SKETCH
UPLOAD SKETCH
$
```

## Build

```bash
go build -ldflags "-X github.com/arduino/scripting-tools/internal/version.Value=v1.0.0"
```
