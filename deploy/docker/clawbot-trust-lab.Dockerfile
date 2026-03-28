FROM golang:1.26-alpine AS build

WORKDIR /src
RUN apk add --no-cache ca-certificates git

ARG TRUST_LAB_VERSION=dev
ARG TRUST_LAB_COMMIT=unknown
ARG TRUST_LAB_BUILD_DATE=unknown

COPY . /src
RUN go build -ldflags "-X 'clawbot-trust-lab/internal/version.Value=${TRUST_LAB_VERSION}' -X 'clawbot-trust-lab/internal/version.Commit=${TRUST_LAB_COMMIT}' -X 'clawbot-trust-lab/internal/version.BuildDate=${TRUST_LAB_BUILD_DATE}'" -o /out/clawbot-trust-lab ./cmd/trust-lab

FROM alpine:3.21

ARG OCI_SOURCE="https://github.com/clawbot-platform/clawbot-trust-lab"
ARG OCI_DESCRIPTION="Version 1 DRQ trust-lab service for agentic commerce trust and fraud benchmarking."
ARG OCI_LICENSES="Apache-2.0"

LABEL org.opencontainers.image.source="${OCI_SOURCE}" \
      org.opencontainers.image.description="${OCI_DESCRIPTION}" \
      org.opencontainers.image.licenses="${OCI_LICENSES}"

RUN apk add --no-cache ca-certificates wget && adduser -D -u 10001 trustlab
COPY --from=build /out/clawbot-trust-lab /usr/local/bin/clawbot-trust-lab
COPY --from=build /src/configs /app/configs

RUN mkdir -p /app/reports /app/var/replay-archive && chown -R trustlab:trustlab /app

USER trustlab
WORKDIR /app
EXPOSE 8090

ENTRYPOINT ["clawbot-trust-lab"]
