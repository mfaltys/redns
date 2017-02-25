package main

import (
	"fmt"
	"io/ioutil"
	"net"
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
	readConf()
	initLogger(config.Doic.Loglevel)

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

func upstreamQuery(w dns.ResponseWriter, req *dns.Msg) *dns.Msg {
	transport := "udp"
	if _, ok := w.RemoteAddr().(*net.TCPAddr); ok {
		transport = "tcp"
	}
	c := &dns.Client{Net: transport}
	resp, _, err := c.Exchange(req, config.Doic.UpstreamDNS)

	if err != nil {
		glogger.Debug.Println(err)
		dns.HandleFailed(w, req)
	}
	return resp
}

func anameresolve(w dns.ResponseWriter, req *dns.Msg) {
	hostname := req.Question[0].Name

	// send request upstream
	glogger.Debug.Printf("sending request for '%s' upstream\n", hostname)
	req = upstreamQuery(w, req)
	w.WriteMsg(req)
}
