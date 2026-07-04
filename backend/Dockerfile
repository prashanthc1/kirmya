# syntax=docker/dockerfile:1

# ---- build stage ----
FROM golang:1.26-alpine AS build
WORKDIR /src
ENV GOSUMDB=off CGO_ENABLED=0 GOOS=linux
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -trimpath -ldflags="-s -w" -o /out/server ./cmd/workspace-app

# ---- dev (hot reload via air; used by docker-compose.override.yml) ----
# Source is bind-mounted; air rebuilds & restarts the server on .go changes.
FROM golang:1.26-alpine AS dev
WORKDIR /src
ENV GOSUMDB=off CGO_ENABLED=0
RUN go install github.com/air-verse/air@latest
EXPOSE 8080
CMD ["air", "-c", ".air.toml"]

# ---- runtime stage ----
FROM alpine:3.20
RUN apk add --no-cache ca-certificates && adduser -D -u 10001 app
WORKDIR /app
COPY --from=build /out/server /app/server
# The server reads ./migrations at boot and serves ./web + ./docs (Swagger UI).
COPY migrations /app/migrations
COPY web /app/web
COPY docs /app/docs
RUN mkdir -p /app/uploads && chown -R app:app /app/uploads
USER app
EXPOSE 8080
ENV PORT=8080
# Container-level health probe is configured by the orchestrator
# (docker-compose / Helm), not baked into the image.
ENTRYPOINT ["/app/server"]
