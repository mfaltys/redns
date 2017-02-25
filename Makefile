GOC=go build
GOFLAGS=-a -ldflags '-s'
CGOR=CGO_ENABLED=0
GIT_HASH=$(shell git rev-parse HEAD | head -c 10)

all: doic

doic:
	$(GOC) doic/*.go

dependencies:

run:
	cd doic && go run \
	doic.go

stat:
	mkdir -p bin/
	$(CGOR) $(GOC) $(GOFLAGS) -o bin/doic-$(GIT_HASH)-linux-amd64 doic/*.go

install: stat
	cp bin/doic* /usr/bin/doic

clean:
	rm -rf bin/
