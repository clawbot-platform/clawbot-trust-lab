FROM golang:1.26-alpine AS build

WORKDIR /src
RUN apk add --no-cache ca-certificates git

COPY clawbot-server /src/clawbot-server

WORKDIR /src/clawbot-server
RUN go build -o /out/clawbot-server ./cmd/clawbot-server

FROM alpine:3.21

ARG OCI_SOURCE="https://github.com/clawbot-platform/clawbot-server"
ARG OCI_DESCRIPTION="Reusable Go-first control-plane foundation for Clawbot Platform services."
ARG OCI_LICENSES="Apache-2.0"

LABEL org.opencontainers.image.source="${OCI_SOURCE}" \
      org.opencontainers.image.description="${OCI_DESCRIPTION}" \
      org.opencontainers.image.licenses="${OCI_LICENSES}"

RUN apk add --no-cache ca-certificates wget && adduser -D -u 10001 clawbot
COPY --from=build /out/clawbot-server /usr/local/bin/clawbot-server

USER clawbot
WORKDIR /app
EXPOSE 8080

ENTRYPOINT ["clawbot-server", "serve"]
