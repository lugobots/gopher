VERSION 0.6
FROM golang:1.17
WORKDIR /app
ARG tag='latest'

deps:
    COPY go.mod go.sum ./
    RUN go mod download
    SAVE ARTIFACT go.mod AS LOCAL go.mod
    SAVE ARTIFACT go.sum AS LOCAL go.sum

build:
    FROM +deps
    COPY . ./
    RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o /bin/my-bot main.go
    SAVE ARTIFACT /bin/my-bot /my-bot

final:
    FROM scratch
    COPY +build/my-bot /bin/my-bot
    ENTRYPOINT ["/bin/my-bot"]
    SAVE IMAGE lugobots/the-dummies:$tag

test-match:
    FROM earthly/dind:alpine
    COPY docker-compose.yml ./

    ENV SERVER_VERSION='v1.0.0-beta'
    ENV HOME_TEAM='lugobots/the-dummies:$tag'
    ENV AWAY_TEAM='lugobots/the-dummies-go:v0.0.0-alpha'

    WITH DOCKER --compose docker-compose.yml --load lugobots/the-dummies:test=+final
        RUN docker run lugobots/the-dummies:test=+final
    END
