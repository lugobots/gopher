FROM golang:1.17

WORKDIR /app

ENTRYPOINT ["/bin/sh", "-c" , "go mod vendor && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags=\"-w -s\" -o my_bot_binary maidn.go"]
CMD []
