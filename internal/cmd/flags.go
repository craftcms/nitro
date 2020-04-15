package cmd

var (
	flagConfigFile      string
	flagMachineName     string
	flagDebug           bool
	flagCPUs            int
	flagMemory          string
	flagDisk            string
	flagPhpVersion      string
	flagDatabase        string
	flagDatabaseVersion string
	flagNginxLogsKind   string
	flagPublicDir       string
	flagPermanent       bool
	flagUpgrade         bool

	// flags for the add command
	flagHostname string
	flagWebroot  string
)
