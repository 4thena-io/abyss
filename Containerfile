FROM alpine:3.22

RUN apk add --no-cache \
  ca-certificates \
  git

ARG TARGETOS TARGETARCH
COPY build/abyss-${TARGETOS}-${TARGETARCH} /usr/bin/abyss

EXPOSE 8000

ENV ABYSS_CONFIG=/etc/abyss/config.yaml
ENV ABYSS_DATA_DIR=/var/lib/abyss

CMD [ "abyss" ]
