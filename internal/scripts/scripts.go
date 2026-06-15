package scripts

type Script interface {
	// Name returns the name of the script, e.g. "dfu-util".
	Name() string

	// Run executes the script with the given arguments.
	Run(args []string) error
}
