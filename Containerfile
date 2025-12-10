FROM --platform=$BUILDPLATFORM docker.io/library/alpine:3.22

RUN apk add --no-cache \
  ca-certificates \
  git \
  sqlite

WORKDIR /opt/app

ARG TARGETARCH
COPY build/abyss-linux-${TARGETARCH} ./abyss

EXPOSE 8000

CMD [ "./abyss" ]
