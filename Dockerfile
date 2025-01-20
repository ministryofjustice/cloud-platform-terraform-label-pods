# syntax=docker/dockerfile:1
FROM golang:1.23.5-alpine AS builder

RUN addgroup -g 1000 -S appgroup && \
  adduser -u 1000 -S appuser -G appgroup

RUN mkdir -p /app/certs

RUN apk --no-cache add ca-certificates

WORKDIR /app

COPY . ./

# Download all the dependencies
RUN go mod download

RUN CGO_ENABLED=0 go build -o /app/main .

RUN chown -R appuser:appgroup /app

FROM scratch

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# copy user permissions from builder
COPY --from=builder /etc/passwd /etc/passwd

COPY --from=builder /app /app

USER 1000

EXPOSE 3000

ENTRYPOINT [ "/app/main" ]

