package config

import (
	"log"

	"github.com/BurntSushi/toml"
)

//Config db url
type Config struct {
	AtlasURI string
	Database string
	APIKey   string
}

//Read config
func (c *Config) Read() {
	if _, err := toml.DecodeFile("config.toml", &c); err != nil {
		log.Fatal(err)
	}
}
