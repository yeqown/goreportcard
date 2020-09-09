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
	# TAG=v1.1.0 @2020-08-12
	docker build -t yeqown/goreportcard:${TAG} .

