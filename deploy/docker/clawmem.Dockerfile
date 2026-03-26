FROM golang:1.26-alpine AS build

WORKDIR /src
RUN apk add --no-cache ca-certificates git

COPY clawmem /src/clawmem

WORKDIR /src/clawmem
RUN go build -o /out/clawmem ./cmd/clawmem

FROM alpine:3.21

RUN apk add --no-cache ca-certificates wget && adduser -D -u 10001 clawmem
COPY --from=build /out/clawmem /usr/local/bin/clawmem
COPY --from=build /src/clawmem/configs /app/configs

RUN mkdir -p /app/var && chown -R clawmem:clawmem /app

USER clawmem
WORKDIR /app
EXPOSE 8088

ENTRYPOINT ["clawmem"]
