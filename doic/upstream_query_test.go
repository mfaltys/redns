package main

import (
	"fmt"
	"testing"

	"github.com/miekg/dns"
)

func TestANameResolve(t *testing.T) {
	// read in conf up a directory
	readConf("../deps/testing.config.gcfg")

	m1 := new(dns.Msg)
	m1.SetQuestion(dns.Fqdn("google.com"), dns.TypeA)

	addr := fmt.Sprintf("127.0.0.1:%d", config.Doic.DNSPort)
	req, err := dns.Exchange(m1, addr)
	if err != nil {
		t.Error("error sending request. is the server running?")
	}

	// print response
	// TODO: logic behind response answer
	//t.Log(req.Question[0])
	t.Log(req.Answer[0])
}
