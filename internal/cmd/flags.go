package cmd

var (
	flagConfigFile    string
	flagMachineName   string
	flagDebug         bool
	flagCPUs          int
	flagMemory        string
	flagDisk          string
	flagPhpVersion    string
	flagNginxLogsKind string

	// flags for the add command
	flagHostname string
	flagWebroot  string
)
