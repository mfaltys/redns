package main

import (
	"fmt"
	"gopkg.in/gcfg.v1"
)

type Config struct {
	Doic struct {
		Loglevel string
	}

	Redis struct {
		Host     string
		Password string
	}
}

var config = Config{}

func main() {
	readConf()
	println("hello world")
}

func readConf() {
	// init the config
	err := gcfg.ReadFileInto(&config, "config.gcfg")
	if err != nil {
		panic(fmt.Sprintf("could not load config.gcfg, error: %s\n", err))
	}
}
