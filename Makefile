.PHONY: all doic
GOC=go build
GOFLAGS=-a -ldflags '-s'
CGOR=CGO_ENABLED=0
GIT_HASH=$(shell git rev-parse HEAD | head -c 10)

all: doic

doic:
	$(GOC) -o bin/doic-$(GIT_HASH)-linux-amd64 doic/*.go

dependencies:
	go get github.com/miekg/dns
	go get github.com/unixvoid/glogger
	go get github.com/gorilla/mux
	go get gopkg.in/gcfg.v1
	go get gopkg.in/redis.v5

run:
	go run \
		doic/doic.go \
		doic/upstream_query.go \
		--port=8053

run_client:
	go run \
		doic_cli/doic_cli.go

stat:
	rm -rf bin/
	mkdir -p bin/
	$(CGOR) $(GOC) $(GOFLAGS) -o bin/doic-$(GIT_HASH)-linux-amd64 doic/*.go

stat_cli:
	rm -rf bin/
	mkdir -p bin/
	$(CGOR) $(GOC) $(GOFLAGS) -o bin/doic_cli doic_cli/*.go

test:
	@$(MAKE) start_test_server || $(MAKE) kill_test_server

test_cli: clean stat_cli
	@echo "LISTING CLIENTS"
	@echo ""
	bin/doic_cli list
	@echo ""
	@echo "GETTING 127.0.0.1 HISTORY"
	@echo ""
	bin/doic_cli get 127.0.0.1

start_test_server: stat
	@bin/doic* -port=8053 & echo $$! > bin/.pid
	go test -v doic/*.go
	@$(MAKE) kill_test_server

kill_test_server:
	@kill `cat bin/.pid`
	@rm -rf bin.pid

install: stat
	cp bin/doic* /usr/bin/doic

generate_domain_list:
	cp deps/getdomains.sh .
	chmod +x getdomains.sh
	./getdomains.sh
	rm getdomains.sh

populate_redis: generate_domain_list
	bash domains > /dev/null

prep_aci: stat
	mkdir -p doic-layout/rootfs/
	cp deps/manifest.json doic-layout/manifest
	cp bin/doic-* doic-layout/rootfs/doic
	cp config.gcfg doic-layout/rootfs/

build_aci: prep_aci
	actool build doic-layout doic.aci
	@echo "doic.aci built"

build_travis_aci: prep_aci
	wget https://github.com/appc/spec/releases/download/v0.8.7/appc-v0.8.7.tar.gz
	tar -zxf appc-v0.8.7.tar.gz
	# build image
	appc-v0.8.7/actool build doic-layout doic.aci && \
	rm -rf appc-v0.8.7*
	@echo "doic.aci built"

clean:
	rm -rf bin/
	rm -rf doic-layout/
	rm -f doic.aci
	rm -f domains
