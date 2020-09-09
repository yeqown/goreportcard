# building stage
FROM golang:1.14-alpine3.11 as build
WORKDIR /tmp/build

COPY . .
RUN export GOPROXY="https://goproxy.cn,direct" \
    && export CGO_ENABLED=0 \
    && export GOARCH=amd64 \
    && export GOOS=linux \
    && go mod download \
    && go build -o app ./cmd/goreportcard-cli/ \
    && go get github.com/golangci/golangci-lint && go install github.com/golangci/golangci-lint/cmd/golangci-lint


# # release stage
# FROM golang:1.14-alpine3.11 as release
# WORKDIR /app/goreportcard
#
# COPY --from=build /tmp/build/app .
# COPY --from=build /tmp/build/tpl ./tpl
# COPY --from=build /tmp/build/assets ./assets
#
# ## 安装git
# RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.ustc.edu.cn/g' /etc/apk/repositories \
#     && apk add git \
#     && apk add openssh \
#     && apk add build-base \
#     && rm -fr /var/cache \
#     && export GOPROXY="https://goproxy.cn,direct"
#
#
# # FIXED: 不能使用golangci-lint
# COPY --from=build /go/bin/golangci-lint /usr/local/bin
#
# EXPOSE 8000
#
# ENTRYPOINT ["./app", "start-web", "&"]


# release stage
FROM alpine as release
WORKDIR /app/goreportcard

COPY --from=build /tmp/build/app .
COPY --from=build /tmp/build/tpl ./tpl
COPY --from=build /tmp/build/assets ./assets
COPY --from=build /usr/local/go/bin/go /usr/local/bin
COPY --from=build /go/bin/golangci-lint /usr/local/bin

## 安装git
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.ustc.edu.cn/g' /etc/apk/repositories \
    && apk add git \
    && apk add openssh \
    && apk add build-base \
    && rm -fr /var/cache \
    && export GOPROXY="https://goproxy.cn,direct"

EXPOSE 8000

ENTRYPOINT ["./app", "start-web", "&"]