# Dockerfile
FROM golang:1.24.2 as builder

WORKDIR /app
COPY . .

# Build a static Linux binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -ldflags="-s -w" -o bin/updates_exporter main.go

# Minimal runtime image
FROM ubuntu:22.04
COPY --from=builder /app/updates_exporter /usr/local/bin/updates_exporter

# Needed to run apt-get inside container (if testing)
#RUN apt-get update && apt-get install -y apt-utils

EXPOSE 9101
CMD ["/usr/local/bin/updates_exporter"]
