FROM alpine:3.14

RUN apk add --no-cache git

git clone --branch v1.10.6 https://github.com/pion/ion-sfu.git

WORKDIR "/ion-sfu"

COPY . .

# build
#go build ./cmd/signal/json-rpc/main.go && ./main -c config.toml
#
#docker run -p 7000:7000 -p 5000-5020:5000-5020/udp pionwebrtc/ion-sfu:v1.10.6-jsonrpc