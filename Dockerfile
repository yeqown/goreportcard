FROM golang:1.14.1-alpine

COPY . .

WORKDIR $GOPATH/src/github.com/gojp/goreportcard

RUN ./scripts/make-install.sh

EXPOSE 8000

CMD ["make", "start"]
