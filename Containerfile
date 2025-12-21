FROM --platform=$BUILDPLATFORM oven/bun:1.3.3-alpine AS frontend

WORKDIR /opt/abyss

COPY ./frontend . 

RUN bun install

RUN bun run build

FROM --platform=$BUILDPLATFORM ghcr.io/techknowlogick/xgo:go-1.24.x AS binary

WORKDIR /opt/abyss

COPY . .
COPY --from=frontend /opt/abyss/dist ./dist

ARG TARGETOS TARGETARCH
RUN xgo \
    -targets '${TARGETOS}/${TARGETARCH}' \
    -dest build -tags 'sqlite' -out abyss -pkg cmd/api .

RUN mkdir -p build
RUN mv /build/* ./build

FROM docker.io/library/alpine:3.22 AS runner

RUN apk add --no-cache \
  ca-certificates \
  git \
  sqlite

COPY --from=binary /opt/abyss/build/abyss-* /usr/bin/abyss

EXPOSE 8000

CMD [ "/usr/bin/abyss" ]
