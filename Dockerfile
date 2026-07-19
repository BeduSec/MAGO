# Copyright (c) BeduSec. All rights reserved.
FROM golang:1.20-alpine AS builder
RUN apk add --no-cache git ca-certificates
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /mago ./cmd/mago

FROM scratch
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /mago /mago
EXPOSE 8080
ENTRYPOINT ["/mago"]
CMD ["serve"]