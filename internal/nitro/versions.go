package nitro

var (
	PHPVersions = []string{"7.4", "7.3", "7.2", "7.1", "7.0"}
	DBEngines   = []string{"mysql", "postgres"}
	DBVersions  = map[string][]string{
		"mysql":    {"8.0", "5.7", "5.6", "5"},
		"postgres": {"12", "12.2", "11.7", "11", "10.12", "10", "9.6", "9.6", "9"},
	}
)
