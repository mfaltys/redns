package main

import (
	"github.com/miekg/dns"
	"github.com/unixvoid/glogger"
)

func anameresolve(w dns.ResponseWriter, req *dns.Msg) {
	hostname := req.Question[0].Name

	// send request upstream
	glogger.Debug.Printf("sending request for '%s' upstream\n", hostname)
	req = upstreamQuery(w, req)
	w.WriteMsg(req)
}
