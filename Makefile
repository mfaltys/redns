.PHONY: all redns
GOC=go build
GOFLAGS=-a -ldflags '-s -X main.version=${VERSION}'
CGOR=CGO_ENABLED=0
GIT_HASH=$(shell git rev-parse HEAD | head -c 10)
VERSION := $(shell git describe --always --long --dirty)


all: redns

redns:
	$(GOC) -o bin/redns-$(GIT_HASH)-linux-amd64 redns/*.go

dependencies:
	go get github.com/miekg/dns
	go get github.com/unixvoid/glogger
	go get github.com/gorilla/mux
	go get gopkg.in/gcfg.v1
	go get gopkg.in/redis.v5

run:
	go run \
		redns/redns.go \
		redns/upstream_query.go \
		--dns=8053 \
		--web=8080

run_client:
	go run \
		redns_cli/redns_cli.go

stat:
	rm -rf bin/
	mkdir -p bin/
	$(CGOR) $(GOC) $(GOFLAGS) -o bin/redns-$(GIT_HASH)-linux-amd64 redns/*.go

stat_cli:
	rm -rf bin/
	mkdir -p bin/
	$(CGOR) $(GOC) $(GOFLAGS) -o bin/redns_cli redns_cli/*.go

test:
	@$(MAKE) start_test_server || $(MAKE) kill_test_server

test_cli: clean stat_cli
	@echo "LISTING CLIENTS"
	@echo ""
	bin/redns_cli list
	@echo ""
	@echo "GETTING 127.0.0.1 HISTORY"
	@echo ""
	bin/redns_cli get 127.0.0.1

start_test_server: stat
	@bin/redns* --dns=8053 --web=8080 & echo $$! > bin/.pid
	go test -v redns/*.go
	@$(MAKE) kill_test_server

kill_test_server:
	@kill `cat bin/.pid`
	@rm -rf bin.pid

install: stat
	cp bin/redns* /usr/bin/redns

generate_domain_list:
	cp deps/getdomains.sh .
	chmod +x getdomains.sh
	./getdomains.sh
	rm getdomains.sh

populate_redis: generate_domain_list
	bash domains > /dev/null

prep_aci: stat
	mkdir -p redns-layout/rootfs/deps/
	cp deps/manifest.json redns-layout/manifest
	cp -R deps/static redns-layout/rootfs/deps/
	cp bin/redns-* redns-layout/rootfs/redns
	cp config.gcfg redns-layout/rootfs/

build_aci: prep_aci
	actool build redns-layout redns.aci
	@echo "redns.aci built"

build_travis_aci: prep_aci
	wget https://github.com/appc/spec/releases/download/v0.8.7/appc-v0.8.7.tar.gz
	tar -zxf appc-v0.8.7.tar.gz
	# build image
	appc-v0.8.7/actool build redns-layout redns.aci && \
	rm -rf appc-v0.8.7*
	@echo "redns.aci built"

clean:
	rm -rf bin/
	rm -rf redns-layout/
	rm -f redns.aci
	rm -f domains
