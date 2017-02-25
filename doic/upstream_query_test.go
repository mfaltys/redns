package main

import (
	"fmt"
	"testing"

	"github.com/miekg/dns"
)

func TestANameResolveCorrect(t *testing.T) {
	// read in conf up a directory
	readConf("../deps/testing.config.gcfg")

	m1 := new(dns.Msg)
	m1.SetQuestion(dns.Fqdn("google.com"), dns.TypeA)

	addr := fmt.Sprintf("127.0.0.1:%d", config.Doic.DNSPort)
	req, err := dns.Exchange(m1, addr)
	if err != nil {
		t.Error("error sending request. is the server running?")
	}

	// print header
	t.Log(req.MsgHdr.String())

	response := fmt.Sprintf("%s", dns.RcodeToString[req.MsgHdr.Rcode])
	if response == "NOERROR" {
		t.Log(req.Answer[0])
		t.Log("\x1b[31mrecieved NOERROR successfully\x1b[39m")
	} else {
		t.Errorf("expected 'NOERROR', got %v instead", response)
	}
}

func TestANameResolveInorrect(t *testing.T) {
	// read in conf up a directory
	readConf("../deps/testing.config.gcfg")

	m1 := new(dns.Msg)
	m1.SetQuestion(dns.Fqdn("not.a.domain"), dns.TypeA)

	addr := fmt.Sprintf("127.0.0.1:%d", config.Doic.DNSPort)
	req, err := dns.Exchange(m1, addr)
	if err != nil {
		t.Error("error sending request. is the server running?")
	}

	// print header
	t.Log(req.MsgHdr.String())

	response := fmt.Sprintf("%s", dns.RcodeToString[req.MsgHdr.Rcode])
	if response == "NXDOMAIN" {
		t.Log("\x1b[31mrecieved NXDOMAIN successfully\x1b[39m")
	} else {
		t.Errorf("expected 'NXDOMAIN', got %v instead", response)
	}
}
