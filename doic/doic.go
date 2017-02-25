package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/unixvoid/glogger"
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
	initLogger(config.Doic.Loglevel)
	glogger.Debug.Println("hello world")
}

func readConf() {
	// init the config
	err := gcfg.ReadFileInto(&config, "config.gcfg")
	if err != nil {
		panic(fmt.Sprintf("could not load config.gcfg, error: %s\n", err))
	}
}

func initLogger(logLevel string) {
	// init logger
	if logLevel == "debug" {
		glogger.LogInit(os.Stdout, os.Stdout, os.Stdout, os.Stderr)
	} else if logLevel == "cluster" {
		glogger.LogInit(os.Stdout, os.Stdout, ioutil.Discard, os.Stderr)
	} else if logLevel == "info" {
		glogger.LogInit(os.Stdout, ioutil.Discard, ioutil.Discard, os.Stderr)
	} else {
		glogger.LogInit(ioutil.Discard, ioutil.Discard, ioutil.Discard, os.Stderr)
	}
}
