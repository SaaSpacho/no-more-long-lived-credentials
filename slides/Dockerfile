FROM node:18-alpine AS builder

WORKDIR /app

RUN npm install reveal.js

FROM nginx:alpine

COPY --from=builder /app/node_modules/ /usr/share/nginx/html/node_modules/
COPY assets/ /usr/share/nginx/html/assets/
COPY index.html /usr/share/nginx/html/index.html

EXPOSE 80
