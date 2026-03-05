FROM golang:1.26.0-alpine3.23 as build

WORKDIR /usr/src/app

COPY go.mod go.sum ./
RUN apk add --no-cache --update git make build-base && go mod download && go mod verify

COPY . .
RUN CGO_ENABLED=1 GOOS=linux go build -ldflags='-s -w -extldflags "-static"' -v -o /usr/local/bin/app

FROM alpine:3.23

RUN apk add --no-cache ca-certificates

COPY --from=build /usr/local/bin/app /usr/local/bin/app

CMD ["app"]