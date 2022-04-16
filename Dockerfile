FROM golang:1.17 AS build

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . ./

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o /bin/my-bot main.go

FROM scratch
COPY --from=build /bin/my-bot /bin/my-bot
ENTRYPOINT ["/bin/my-bot"]
CMD []
