# build stage
FROM node:22.21.1-alpine AS build-stage

WORKDIR /app

COPY package*.json ./
RUN npm ci

COPY . .
RUN npm run build

# production stage
FROM caddy:2-alpine AS production-stage
COPY Caddyfile /etc/caddy/Caddyfile
COPY --from=build-stage /app/dist /usr/share/caddy

EXPOSE 80
EXPOSE 443
CMD ["caddy", "run", "--config", "/etc/caddy/Caddyfile", "--adapter", "caddyfile"]
