# Description: Dockerfile to build the currency-api application
# Author: Gusto (https://github.com/parmaster)
# Usage: 
#	docker build -t currency-api .
#	docker run -it --rm -p 8080:8080 currency-api
#	curl http://localhost:8080/v1/status
#	curl http://localhost:8080/v1/rates
#
# Build stage
FROM golang:1.22-bullseye as base

RUN adduser \
  --disabled-password \
  --gecos "" \
  --home "/nonexistent" \
  --shell "/sbin/nologin" \
  --no-create-home \
  --uid 65532 \
  api-user

ADD . /build
WORKDIR /build

# Build the application with the vendor dependencies and the cache
RUN --mount=type=cache,target="/root/.cache/go-build" CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -mod=vendor -v -o currency-api ./cmd/api

# Production stage
# Not using scratch because of the CGO_ENABLED=1
FROM debian:stable-slim

WORKDIR /app

# Copy the certificates, passwd and group files from the base image
COPY --from=base /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=base /etc/passwd /etc/passwd
COPY --from=base /etc/group /etc/group

# Copy the binary and the configuration file from the base image
COPY --from=base /build/currency-api .
COPY --from=base /build/config.ini .

# Change the ownership of the binary and the configuration file
RUN chmod +x ./currency-api

# Set the user to run the application
USER api-user:api-user

# Expose the port
EXPOSE 8080

CMD ["./currency-api"]