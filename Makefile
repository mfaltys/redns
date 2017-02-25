GOC=go build
GOFLAGS=-a -ldflags '-s'
CGOR=CGO_ENABLED=0
GIT_HASH=$(shell git rev-parse HEAD | head -c 10)

all: doic

doic:
	$(GOC) doic/*.go

dependencies:
	go get github.com/miekg/dns
	go get github.com/unixvoid/glogger
	go get gopkg.in/gcfg.v1

run:
	go run \
		doic/doic.go \
		doic/aname_resolve.go \
		doic/upstream_query.go

stat:
	rm -rf bin/
	mkdir -p bin/
	$(CGOR) $(GOC) $(GOFLAGS) -o bin/doic-$(GIT_HASH)-linux-amd64 doic/*.go

test:
	@$(MAKE) start_test_server || $(MAKE) kill_test_server

start_test_server: stat
	@bin/doic* -port=8053 & echo $$! > bin/.pid
	go test -v doic/*.go
	@$(MAKE) kill_test_server

kill_test_server:
	@kill `cat bin/.pid`
	@rm -rf bin.pid

install: stat
	cp bin/doic* /usr/bin/doic

clean:
	rm -rf bin/
