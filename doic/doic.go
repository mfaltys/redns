package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/miekg/dns"
	"github.com/unixvoid/glogger"
	"gopkg.in/gcfg.v1"
)

type Config struct {
	Doic struct {
		Loglevel    string
		DNSPort     int
		UpstreamDNS string
	}

	Redis struct {
		Host     string
		Password string
	}
}

var config = Config{}

func main() {
	readConf("config.gcfg")
	initLogger(config.Doic.Loglevel)

	// parse override flags
	overrideDNSPort := flag.Int("port", config.Doic.DNSPort, "DNS port to bind to.")
	flag.Parse()

	if *overrideDNSPort != config.Doic.DNSPort {
		config.Doic.DNSPort = *overrideDNSPort
	}

	// format the string to be :port
	fPort := fmt.Sprint(":", config.Doic.DNSPort)

	udpServer := &dns.Server{Addr: fPort, Net: "udp"}
	tcpServer := &dns.Server{Addr: fPort, Net: "tcp"}
	glogger.Info.Println("started server on", config.Doic.DNSPort)

	dns.HandleFunc(".", func(w dns.ResponseWriter, req *dns.Msg) {
		switch req.Question[0].Qtype {
		case 1:
			glogger.Debug.Println("'A' request recieved, continuing")
			go anameresolve(w, req)
			break
		case 5:
			glogger.Debug.Println("'CNAME' request detected: TODO")
			break
		case 28:
			glogger.Debug.Println("'AAAA' request detected: TODO")
			break
		default:
			glogger.Debug.Printf("non supported '%d' request detected. skipping", req.Question[0].Qtype)
			break
		}
	})

	go func() {
		glogger.Error.Println(udpServer.ListenAndServe())
	}()
	glogger.Error.Println(tcpServer.ListenAndServe())

}

func readConf(location string) {
	// init the config
	err := gcfg.ReadFileInto(&config, location)
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
