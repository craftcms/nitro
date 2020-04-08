package nitro

type Command struct {
	Machine   string
	Type      string
	Chainable bool
	Input     string
	Args      []string
}
