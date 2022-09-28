FROM golang:1.19.1-alpine as builder

WORKDIR /app

RUN apk update && apk add --no-cache git make ca-certificates tzdata && update-ca-certificates

COPY . /app

RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 make ethmetrics

RUN ls -lah /app/bin/

FROM scratch

WORKDIR /app

COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /app/bin/* /app

ENTRYPOINT ["./ethmetrics"]
