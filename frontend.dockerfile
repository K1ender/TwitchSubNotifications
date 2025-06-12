FROM oven/bun:latest AS builder

WORKDIR /app

COPY web/package.json web/bun.lock ./
RUN bun install

COPY web ./

RUN bun run build-only

FROM nginx:alpine

COPY --from=builder /app/dist /usr/share/nginx/html

EXPOSE 80

CMD ["nginx", "-g", "daemon off;"]