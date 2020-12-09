package viewproxy

import "time"

// Queries describes a single query object consisting of a descriptive name and the corresponding PromQL query.
type Queries struct {
	Query string `yaml:"query"`
	Name  string `yaml:"name"`
}

// Prometheus describes the configuration structure for connecting to the upstream Prometheus
type Prometheus struct {
	URL string `yaml:"url"`
}

// Config for the proxy
type Config struct {
	Routes map[string]struct {
		Queries  []Queries `yaml:"queries"`
		Template string    `yaml:"template"`
	} `yaml:"routes"`
	ResponseExpiryTime time.Duration `yaml:"responseExpiryTime"`
	Prometheus         Prometheus    `yaml:"prometheus"`
}
