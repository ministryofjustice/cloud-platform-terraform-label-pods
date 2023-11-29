# syntax=docker/dockerfile:1
FROM golang:1.21-alpine AS builder

RUN addgroup -g 1000 -S appgroup && \
  adduser -u 1000 -S appuser -G appgroup

RUN mkdir /artifact

RUN apk --no-cache add ca-certificates

USER 1000

EXPOSE 3000

ENTRYPOINT ["/bin/sh", "-c", "cp", "-r", "/etc/ssl/certs/*", "/artifact/", "&&", "cp", "/etc/ssl/tmp/tls.crt", "/artifact", "&&", "cp", "/etc/ssl/tmp/ca.crt", "/artifact", "&&", "cp", "/etc/ssl/tmp/tls.key", "/artifact"]

