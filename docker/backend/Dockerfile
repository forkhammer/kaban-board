FROM golang:1.21 as build

WORKDIR /usr/src/app

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN go build -ldflags='-s -w -extldflags "-static"' -v -o /usr/local/bin/app

FROM debian:12-slim

COPY --from=build /usr/local/bin/app /usr/local/bin/app

CMD ["app"]