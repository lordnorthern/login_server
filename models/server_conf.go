package models

import (
	"io/ioutil"
	"log"

	"github.com/BurntSushi/toml"
)

// Config will hold all the configuration
type Config struct {
	MySQL struct {
		DbHost     string
		DbPort     string
		DbName     string
		DbUsername string
		DbPassword string
	}

	PublicServer   ServConf
	InternalServer ServConf
}

// ServConf holds the info for every server instance
type ServConf struct {
	Port       string
	Definition string
}

// Conf will hold all the configuration data
var Conf Config

// ParseConf will parse the conf file
func ParseConf() {
	confFileContents, err := ioutil.ReadFile("serv.conf")
	if err == nil {
		toml.Decode(string(confFileContents), &Conf)
	} else {
		log.Fatal(err)
	}
}
