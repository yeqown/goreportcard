# build
FROM golang:1.14.1 as build

WORKDIR /tmp/build

COPY . .

RUN export GOPROXY="https://goproxy.cn,direct" \
    && go mod download \
    && CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build -o cli ./cmd/goreportcard-cli/ \
    && ./scripts/make-install.sh \
    && which golangci-lint

# release
FROM alpine as release

WORKDIR /app/goreportcard

COPY --from=build /tmp/build/cli .
COPY --from=build /tmp/build/templates ./templates
COPY --from=build /tmp/build/assets ./assets
COPY --from=build /go/bin/golangci-lint /usr/local/bin

EXPOSE 8000

 ENTRYPOINT ["./cli", "start-web"]