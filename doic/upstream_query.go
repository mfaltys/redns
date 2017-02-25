package main

import (
	"net"

	"github.com/miekg/dns"
	"github.com/unixvoid/glogger"
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
