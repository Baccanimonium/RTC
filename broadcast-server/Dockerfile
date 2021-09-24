# syntax=docker/dockerfile:1
FROM golang:1.16-alpine

RUN apk add build-base

RUN apk add --no-cache v4l-utils
RUN apk add --no-cache libvpx-dev

WORKDIR /app

COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY *.go ./

RUN go build -o /broadcast-server

EXPOSE 8000

CMD [ "/broadcast-server" ]
# build
#go build ./cmd/signal/json-rpc/main.go && ./main -c config.toml
#
#docker run -p 7000:7000 -p 5000-5020:5000-5020/udp pionwebrtc/ion-sfu:v1.10.6-jsonrpc
