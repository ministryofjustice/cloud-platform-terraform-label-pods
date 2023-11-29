# syntax=docker/dockerfile:1
FROM golang:1.21-alpine

RUN mkdir /artifact

RUN apk --no-cache add ca-certificates

ENTRYPOINT [ "echo", "completed" ]

