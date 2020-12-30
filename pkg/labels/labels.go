package labels

const (
	Nitro = "com.craftcms.nitro"

	// DatabaseEngine is used to identify the engine that is being used for a database container (e.g. mysql, postgres)
	DatabaseEngine = "com.craftcms.nitro.database-engine"

	// DatabaseVersion is the version of the datbase the container is running (e.g. 11, 12, 5.7)
	DatabaseVersion = "com.craftcms.nitro.database-version"

	// DatabaseCompatability is the compatability of the database (e.g. mariadb and mysql are compatible)
	DatabaseCompatability = "com.craftcms.nitro.database-compatability"

	// Host is used to identify a web application by the hostname of the site (e.g demo.nitro)
	Host = "com.craftcms.nitro.host"

	// Network is used to label a network for an environment
	Network = "com.craftcms.nitro.network"

	// Volume is used to identify a volume for an environment
	Volume = "com.craftcms.nitro.volume"

	// Proxy is the label used to idenitfy the proxy container
	Proxy = "com.craftcms.nitro.proxy"

	// ProxyVersion is used to label a proxy container with a specific version
	ProxyVersion = "com.craftcms.nitro.proxy-version"

	// Type is used to idenity the type of container
	Type = "com.craftcms.nitro.type"
)
