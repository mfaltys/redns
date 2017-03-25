package main

import (
	"net"
	"strings"

	"github.com/miekg/dns"
	"github.com/unixvoid/glogger"
	"gopkg.in/redis.v5"
)

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

func anameresolve(w dns.ResponseWriter, req *dns.Msg, redisClient *redis.Client) {
	hostname := req.Question[0].Name

	// send request upstream
	client := strings.Split(w.RemoteAddr().String(), ":")
	glogger.Debug.Printf("client: %s\n", client[0])
	glogger.Debug.Printf("sending request for '%s' upstream\n", hostname)
	req = upstreamQuery(w, req)
	w.WriteMsg(req)
}
