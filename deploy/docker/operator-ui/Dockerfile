FROM node:22-alpine AS build

WORKDIR /app
COPY clawbot-trust-lab/web/package.json clawbot-trust-lab/web/package-lock.json ./
RUN npm ci
COPY clawbot-trust-lab/web/ ./
RUN npm run build

FROM nginx:1.27-alpine

COPY clawbot-trust-lab/deploy/docker/operator-ui/nginx.conf /etc/nginx/conf.d/default.conf
COPY --from=build /app/dist /usr/share/nginx/html

EXPOSE 8080
