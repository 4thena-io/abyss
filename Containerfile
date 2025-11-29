FROM alpine:3.22

RUN apk add --no-cache git ca-certificates

WORKDIR /opt/build

ARG TARGETARCH
COPY build/app-${TARGETARCH}.out ./app.out

EXPOSE 8000

CMD [ "./app.out" ]
