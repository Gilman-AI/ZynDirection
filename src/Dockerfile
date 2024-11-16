FROM golang:1.23.3-bookworm as builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -ldflags '-s -w' -o app

FROM scratch
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /app/app /app/app
EXPOSE 8080
USER nonroot:nonroot
ENTRYPOINT ["/app/app"]