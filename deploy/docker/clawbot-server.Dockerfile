FROM golang:1.26-alpine AS build

WORKDIR /src
RUN apk add --no-cache ca-certificates git

COPY clawbot-server /src/clawbot-server

WORKDIR /src/clawbot-server
RUN go build -o /out/clawbot-server ./cmd/clawbot-server

FROM alpine:3.21

RUN apk add --no-cache ca-certificates wget && adduser -D -u 10001 clawbot
COPY --from=build /out/clawbot-server /usr/local/bin/clawbot-server

USER clawbot
WORKDIR /app
EXPOSE 8080

ENTRYPOINT ["clawbot-server", "serve"]
