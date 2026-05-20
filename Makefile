.PHONY: build test lint fuzz clean install

build:
	go build -ldflags="-s -w -X main.version=$$(git describe --tags --always)" \
	  -o bin/gitresolve ./cmd/gitresolve

test:
	go test -race -count=1 ./...

lint:
	golangci-lint run ./...

fuzz:
	go test -fuzz=FuzzParser -fuzztime=30s ./internal/conflict/

clean:
	rm -rf bin/

install:
	go install ./cmd/gitresolve
