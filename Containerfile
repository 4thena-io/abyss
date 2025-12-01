FROM alpine:3.22

RUN apk add --no-cache ca-certificates git libc6-compat sqlite-libs

WORKDIR /opt/app

ARG TARGETARCH
COPY build/app-${TARGETARCH} ./app

EXPOSE 8000

CMD [ "./app.out" ]
