package config

type server struct {
	Port     int
	Host     string
	Protocol string
}

type Config struct {
	Server server
}

func New() *Config {
	// Configuring basic server
	// Will be upgrade in future

	return &Config{
		Server: server{
			Port:     8080,
			Host:     "localhost",
			Protocol: "http",
		},
	}
}
