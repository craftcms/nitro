package cmd

var (
	flagMachineName   string
	flagDebug         bool
	flagCPUs          int
	flagMemory        string
	flagDisk          string
	flagPhpVersion    string
	flagNginxLogsKind string
	flagClean         bool
	flagSkipBackup    bool

	// flags for the add command
	flagHostname string
	flagWebroot  string

	// flags for apply
	flagSkipHosts bool
)
