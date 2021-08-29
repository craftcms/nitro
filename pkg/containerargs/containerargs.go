package containerargs

// Command is the struct that contains the name of the container and the args to pass
// into the container. It does not know about the image to use and instead only focuses
// on breaking the args from a command (e.g `nitro composer sitename.nitro -- install`)
// would result in a Command with the Container `sitename.nitro` and the Args []string{"install"}
type Command struct {
	Container string
	Args      []string
}

// Parse is designed to take args from a command and determine if the container name is the first
// argument.
func Parse(args []string) (*Command, error) {
	// check if the args contains the delimiter
	delimiter := false
	var pos int
	for k, v := range args {
		if v == "--" {
			pos = k
			delimiter = true
		}
	}

	// if the delimiter was provided get the name of the container and get the args
	if delimiter {
		return &Command{
			Container: args[0],
			Args:      args[pos+1:],
		}, nil
	}

	return &Command{Args: args}, nil
}
