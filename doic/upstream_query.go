package main

import (
	"fmt"
	"net"
	"strings"
	"time"

	"github.com/miekg/dns"
	"github.com/unixvoid/glogger"
	"gopkg.in/redis.v5"
)

func upstreamQuery(w dns.ResponseWriter, req *dns.Msg, redisClient *redis.Client) *dns.Msg {
	// parse hostname
	hostname := req.Question[0].Name

	// parse client ip
	client := strings.Split(w.RemoteAddr().String(), ":")

	// un-fully qualify domain if its qualified
	hostLen := len(hostname)
	if hostLen > 0 && hostname[hostLen-1] == '.' {
		hostname = hostname[:hostLen-1]
	}

	// query redis to see if entry exists in 'blacklist:domain' O(1)
	exists, err := redisClient.SIsMember("blacklist:domain", hostname).Result()
	if err != nil {
		glogger.Error.Println("error getting result from blacklist:domain")
		glogger.Error.Println(err)
	}

	// handle blacklisted domain case
	if exists {
		// TODO add 'nonexistent' response to client
		glogger.Debug.Printf("intercepted blacklisted domain '%s' on client '%s'\n", hostname, client[0])
	}

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

func anameresolve(w dns.ResponseWriter, req *dns.Msg, redisClient *redis.Client) {
	// parse hostname
	hostname := req.Question[0].Name

	// parse client ip
	client := strings.Split(w.RemoteAddr().String(), ":")

	// generate timestamp
	t := time.Now()
	timestamp := fmt.Sprintf("%s", t.Format("2006/01/02 15:04:05"))
	//timestamp := fmt.Sprintf("%s", t.Format(time.RFC1123))

	// add client to redis client:list
	err := redisClient.SAdd("client:list", client[0]).Err()
	if err != nil {
		glogger.Error.Printf("error adding client: '%s' to 'client:list'", client[0])
	}

	// add client history entry to redis
	err = redisClient.RPush(fmt.Sprintf("client:%s", client[0]), fmt.Sprintf("%s :: %s", timestamp, hostname)).Err()
	if err != nil {
		glogger.Error.Printf("error adding hostname: '%s' for client: 'client:%s'\n", hostname, client)
		glogger.Error.Printf("%s", err)
	}

	// send request upstream
	glogger.Debug.Printf("client: %s\n", client[0])
	glogger.Debug.Printf("sending request for '%s' upstream\n", hostname)
	req = upstreamQuery(w, req, redisClient)
	// write response back from client
	w.WriteMsg(req)
}
