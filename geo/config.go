package geo

// Config represents the configuration objects necessary to
// use the objects in this package.
type Config struct {
	LocationIQAPIKey string
}

// Init returns a configuration struct that can be used to initialize
// the objects in this package.
func Init(locationIQAPIKey string) *Config {
	return &Config{locationIQAPIKey}
}
