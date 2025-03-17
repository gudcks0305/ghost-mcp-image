package config

import (
	"os"
)

type MyConfig struct {
	GHOST_API_URL       string
	GHOST_STAFF_API_KEY string
}

var Config *MyConfig

func init() {
	Config = &MyConfig{
		GHOST_API_URL:       os.Getenv("GHOST_API_URL"),
		GHOST_STAFF_API_KEY: os.Getenv("GHOST_STAFF_API_KEY"),
	}
}
