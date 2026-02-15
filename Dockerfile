FROM golang:1.26-bookworm as go_builder

WORKDIR /usr/src/myservice

# pre-copy/cache go.mod for pre-downloading dependencies and only redownloading them in subsequent builds if they change
COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .

RUN make build

FROM ghcr.io/cirruslabs/flutter:3.41.0 as flutter_builder
WORKDIR /app/portal
ARG portal_backend_scheme=http
ARG portal_backend_host=localhost
ARG portal_backend_port=8080
ARG portal_backend_path=/internal

COPY portal .

RUN flutter pub get
RUN flutter build web --release --base-href=$portal_backend_path/portal/ --dart-define=portal_backend_scheme=$portal_backend_scheme --dart-define=portal_backend_host=$portal_backend_host --dart-define=portal_backend_port=$portal_backend_port --dart-define=portal_backend_path=$portal_backend_path

FROM debian:bookworm-slim

ENV SERVER_PORTAL_PATH /var/www

RUN set -x && apt-get update && DEBIAN_FRONTEND=noninteractive apt-get install -y \
    ca-certificates && \
    rm -rf /var/lib/apt/lists/*

# Copy Portal
COPY --from=flutter_builder /app/portal/build/web /var/www

# Copy the binary to the production image from the builder stage.
COPY --from=go_builder /usr/src/myservice/bin/myservice /usr/local/bin/myservice

# Copy configuration file
COPY application.yaml .

# Run the web service on container startup.
CMD ["myservice", "start"]
