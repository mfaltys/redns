package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/miekg/dns"
	"github.com/unixvoid/glogger"
	"gopkg.in/gcfg.v1"
	"gopkg.in/redis.v5"
)

type Config struct {
	Doic struct {
		Loglevel       string
		DNSPort        int
		APIPort        int
		UpstreamDNS    string
		BootstrapDelay time.Duration
		RedirectSite   string
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

	// initialize redis connection
	redisClient, err := initRedisConnection()
	if err != nil {
		glogger.Debug.Println("redis conneciton cannot be made, trying again in 1 second")
		time.Sleep(config.Doic.BootstrapDelay * time.Second)
		redisClient, err = initRedisConnection()
		if err != nil {
			glogger.Error.Println("redis connection cannot be made.")
			os.Exit(1)
		}
	} else {
		glogger.Debug.Println("connection to redis succeeded.")
		glogger.Info.Println("link to redis on", config.Redis.Host)
	}

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

	// grab external ip for debugging
	externalIp := getoutboundIP()
	glogger.Info.Printf("external ip: %s\n", externalIp)

	// set up router for static pages
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// func here shillz(w, r, client)
		statichandler(w, r)
	}).Methods("GET")

	// serve up the web view
	glogger.Info.Printf("web frontend running on %d\n", config.Doic.APIPort)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", config.Doic.APIPort), router))

	dns.HandleFunc(".", func(w dns.ResponseWriter, req *dns.Msg) {
		switch req.Question[0].Qtype {
		case 1:
			glogger.Debug.Println("'A' request recieved, continuing")
			go anameresolve(w, req, redisClient)
			break
		case 5:
			glogger.Debug.Println("'CNAME' request detected: TODO")
			break
		case 28:
			glogger.Debug.Println("'AAAA' request recieved, continuing")
			go aaaanameresolve(w, req, redisClient)
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

func initRedisConnection() (*redis.Client, error) {
	// initialize redis connection
	client := redis.NewClient(&redis.Options{
		Addr:     config.Redis.Host,
		Password: config.Redis.Password,
		DB:       0,
	})
	_, redisErr := client.Ping().Result()
	return client, redisErr
}

// Get preferred outbound ip of this machine
func getoutboundIP() string {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		glogger.Error.Println(err)
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().String()
	idx := strings.LastIndex(localAddr, ":")

	return localAddr[0:idx]
}

func statichandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "hello warld")
}
