FROM alpine:3.22

RUN apk add --no-cache \
  ca-certificates \
  git

ARG TARGETOS TARGETARCH
COPY build/abyss-${TARGETOS}-${TARGETARCH} /usr/bin/abyss

EXPOSE 8000

CMD [ "abyss" ]
