FROM node:22 as build

ARG API_URL
ENV API_URL $API_URL

WORKDIR /var/www/

COPY ./ /var/www/

RUN npm ci \
    && npm i -g envsub \
    && envsub /var/www/src/assets/env.template.js /var/www/src/assets/env.js \
    && npm run build

FROM caddy:2.10.0-alpine

COPY --from=build /var/www/dist/board/browser/ /var/www/board/
COPY Caddyfile /etc/caddy/Caddyfile
