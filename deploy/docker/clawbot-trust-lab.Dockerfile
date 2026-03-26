FROM golang:1.26-alpine AS build

WORKDIR /src
RUN apk add --no-cache ca-certificates git

COPY clawbot-trust-lab /src/clawbot-trust-lab

WORKDIR /src/clawbot-trust-lab
RUN go build -o /out/clawbot-trust-lab ./cmd/trust-lab

FROM alpine:3.21

RUN apk add --no-cache ca-certificates wget && adduser -D -u 10001 trustlab
COPY --from=build /out/clawbot-trust-lab /usr/local/bin/clawbot-trust-lab
COPY --from=build /src/clawbot-trust-lab/configs /app/configs

RUN mkdir -p /app/reports /app/var/replay-archive && chown -R trustlab:trustlab /app

USER trustlab
WORKDIR /app
EXPOSE 8090

ENTRYPOINT ["clawbot-trust-lab"]
