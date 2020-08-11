all: lint build test

build:
	go build ./cmd/goreportcard-cli

lint:
	golangci-lint run --skip-dirs=repos --disable-all \
		--enable=golint --enable=vet --enable=gofmt --enable=misspell ./...

test: 
	go test -cover ./internal

run: build
	./goreportcard-cli start-web

image:
	docker build -t goreportcard:v1.0.0 .

