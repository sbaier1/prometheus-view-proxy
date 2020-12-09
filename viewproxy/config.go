package viewproxy

type Queries struct {
	Query string `yaml:"query"`
	Name  string `yaml:"name"`
}

// Config for the proxy
type Config struct {
	Routes map[string]struct {
		Queries  []Queries `yaml:"queries"`
		Template string    `yaml:"template"`
	} `yaml:"routes"`
}
