package client

const (
	// DatabaseEngineLabel is used to identify the engine that is being used for a database container (e.g. mysql, postgres)
	DatabaseEngineLabel = "com.craftcms.nitro.database-engine"

	// DatabaseVersionLabel is the version of the datbase the container is running (e.g. 11, 12, 5.7)
	DatabaseVersionLabel = "com.craftcms.nitro.database-version"

	// EnvironmentLabel is the constant used for the label used for docker images to determine
	// the environment (e.g. nitro-dev) for the container
	EnvironmentLabel = "com.craftcms.nitro.environment"

	// HostLabel is used to identify a web application by the hostname of the site (e.g demo.nitro)
	HostLabel = "com.craftcms.nitro.host"
)
