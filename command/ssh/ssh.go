package ssh

var (
	// RootUser is used to tell the container to run as root and not the default user www-data
	RootUser bool

	// ProxyContainer is used to ssh into the proxy container and is mostly used for troubleshooting
	ProxyContainer bool
)

const exampleText = `  # ssh into a container - assuming its the current working directory
  nitro ssh

  # ssh into the container as root - changes may not persist after "nitro apply"
  nitro ssh --root

  # ssh into the proxy container
  nitro ssh --proxy`
