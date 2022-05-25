package main

import "github.com/vrischmann/envconfig"

var config struct {
	HTTP struct {
		Port int `envconfig:"HTTP_PORT,default=8080" json:"Port"`
	}
	Debug struct {
		PrintConfig    bool `envconfig:"DEBUG_PRINTCONFIG,default=false" json:"PrintConfig"`
		PrintRootCause bool `envconfig:"DEBUG_PRINTROOTCAUSE,default=false" json:"PrintRootCause"`
	}
	CommandDatabase struct {
		URL string `envconfig:"COMMAND_DATABASE_URL" json:"CommandDatabaseURL"`
	}
}

/// initialize configuration.
func initConfig() error {
	return envconfig.Init(&config)
}
