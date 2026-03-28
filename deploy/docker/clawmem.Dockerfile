FROM golang:1.26-alpine AS build

WORKDIR /src

RUN apk add --no-cache ca-certificates git

COPY clawmem /src/clawmem
WORKDIR /src/clawmem

RUN go build -o /out/clawmem ./cmd/clawmem

FROM alpine:3.21 AS stage-1

ARG OCI_SOURCE="https://github.com/clawbot-platform/clawmem"
ARG OCI_DESCRIPTION="Reusable Go-first memory, replay, and historical context service."
ARG OCI_LICENSES="Apache-2.0"

LABEL org.opencontainers.image.source="${OCI_SOURCE}" \
      org.opencontainers.image.description="${OCI_DESCRIPTION}" \
      org.opencontainers.image.licenses="${OCI_LICENSES}"

RUN apk add --no-cache ca-certificates wget && adduser -D -u 10001 clawmem

COPY --from=build /out/clawmem /usr/local/bin/clawmem
COPY --from=build /src/clawmem/configs /app/configs

RUN mkdir -p /app/var /data/clawmem && chown -R clawmem:clawmem /app /data

WORKDIR /app
USER clawmem

EXPOSE 8088
ENTRYPOINT ["/usr/local/bin/clawmem"]
