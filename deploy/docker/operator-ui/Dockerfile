FROM node:22-alpine AS build

WORKDIR /app
COPY web/package.json web/package-lock.json ./
RUN npm ci
COPY web/ ./
RUN npm run build

FROM nginx:1.27-alpine

ARG OCI_SOURCE="https://github.com/clawbot-platform/clawbot-trust-lab"
ARG OCI_DESCRIPTION="Operator UI for the Clawbot Trust Lab Version 1 DRQ runtime."
ARG OCI_LICENSES="Apache-2.0"

LABEL org.opencontainers.image.source="${OCI_SOURCE}" \
      org.opencontainers.image.description="${OCI_DESCRIPTION}" \
      org.opencontainers.image.licenses="${OCI_LICENSES}"

COPY deploy/docker/operator-ui/nginx.conf /etc/nginx/conf.d/default.conf
COPY --from=build /app/dist /usr/share/nginx/html

EXPOSE 8080
