package config

type Config struct {
	Host     string
	Port     string
	Username string
	Password string
	Database string
}

func LoadConfig() Config {
	return Config{
		Host:     "localhost",
		Port:     "5400",
		Username: "postgres",
		Password: "docker",
		Database: "postgres",
	}
}
