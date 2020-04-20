package nitro

var (
	PHPVersions = []string{"7.4", "7.3", "7.2", "7.1", "7.0"}
	DBEngines   = []string{"mysql", "postgres"}
	DBVersions  = map[string][]string{
		"mysql":    {"8.0", "5.7", "5.6"},
		"postgres": {"12", "11", "10", "9"},
	}
)
