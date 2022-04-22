package config


type Config struct {
	MaxDepth   int
}

func NewConfig(maxDepth int) *Config {
	return &Config{maxDepth}
}
