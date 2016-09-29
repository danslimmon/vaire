package main

import (
	"github.com/kelseyhightower/envconfig"
)

var Config ConfigStruct

type ConfigStruct struct {
	Debug  bool
	Listen string `required:"true"`
	Token  string `required:"true"`
}

// LoadConfig() loads Vairë's config into the `Config` variable.
func LoadConfig() error {
	return envconfig.Process("vaire", &Config)
}
