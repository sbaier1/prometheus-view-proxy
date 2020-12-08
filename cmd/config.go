package main

// Config for the proxy
type Config struct {
	Routes map[string]struct {
		Query    string `yaml:"query"`
		Template string `yaml:"template"`
	} `yaml:"routes"`
}
