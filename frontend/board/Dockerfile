FROM node:22 as build

ARG API_URL
ENV API_URL $API_URL

WORKDIR /var/www/

COPY ./ /var/www/

RUN npm ci \
    && npm i -g envsub \
    && envsub /var/www/src/assets/env.template.js /var/www/src/assets/env.js \
    && npm run build

FROM nginx:1.29

COPY --from=build /var/www/dist/board/browser/ /var/www/board/
COPY nginx.conf /etc/nginx/conf.d/default.conf
