package scripts

// Script is a struct that is used for sending "commands" or "scripts" to the machine to run.
// An example script would be to install a package, configure xdebug, or restart a machine service.
// The name field is used to give the machine a friendly name, it also is used to display output
// to the console. Args are the actual commands to run, the args do not contain the machine name,
// only the commands to be executed inside the machine.
type Script struct {
	Name string
	Args []string
}
