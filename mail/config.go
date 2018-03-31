package mail

type Config struct {
	From     string
	Host     string
	Port     int
	Password string
	User     string
}

func New(from string, host string, port int, password string, user string) *Config {
	return &Config{from, host, port, password, user}
}
