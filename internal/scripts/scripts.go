package scripts

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
