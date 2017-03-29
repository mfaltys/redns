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

func upstreamQuery(w dns.ResponseWriter, req *dns.Msg) *dns.Msg {
	// TODO
	//   - check hostname against malicious domain list
	//   - check ip against malicious domain list
	//   - handler when malicious address/domain is detected

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
	hostname := req.Question[0].Name
	client := strings.Split(w.RemoteAddr().String(), ":")
	t := time.Now()
	timestamp := fmt.Sprintf("%s", t.Format("2006/01/02 15:04:05"))
	//timestamp := fmt.Sprintf("%s", t.Format(time.RFC1123))

	// add client to redis client:list
	err := redisClient.SAdd("client:list", client[0]).Err()
	if err != nil {
		glogger.Error.Printf("error adding client: '%s' to 'client:list'", client[0])
	}

	err = redisClient.RPush(fmt.Sprintf("client:%s", client[0]), fmt.Sprintf("%s :: %s", timestamp, hostname)).Err()
	if err != nil {
		glogger.Error.Printf("error adding hostname: '%s' for client: 'client:%s'\n", hostname, client)
		glogger.Error.Printf("%s", err)
	}

	// send request upstream
	glogger.Debug.Printf("client: %s\n", client[0])
	glogger.Debug.Printf("sending request for '%s' upstream\n", hostname)
	req = upstreamQuery(w, req)
	w.WriteMsg(req)
}
