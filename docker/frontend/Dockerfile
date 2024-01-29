FROM node:21 as build

ARG API_URL
ENV API_URL $API_URL

WORKDIR /var/www/

COPY ./ /var/www/

RUN npm i -g --unsafe-perm @angular/cli@16.2.12 \
    && npm i -g \
    && npm i -g envsub \
    && envsub /var/www/src/assets/env.template.js /var/www/src/assets/env.js \
    && npm run build

FROM caddy:2.7.6-alpine

COPY --from=build /var/www/dist/board/ /var/www/board/
COPY Caddyfile /etc/caddy/Caddyfile
