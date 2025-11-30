# syntax=docker/dockerfile:1

FROM golang:1.22 AS builder
WORKDIR /app/backend
COPY backend/go.mod backend/go.sum ./
RUN go mod download
COPY backend/ ./
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /bin/events-api ./cmd/api
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /bin/migrate ./cmd/migrate

FROM gcr.io/distroless/base-debian12
WORKDIR /srv
COPY --from=builder /bin/events-api ./events-api
COPY --from=builder /bin/migrate ./migrate
COPY --from=builder /app/backend/docs ./docs
ENV PORT=8080 \
    DATABASE_PATH=/tmp/data.db \
    UPLOAD_DIR=/tmp/uploads \
    GIN_MODE=release
EXPOSE 8080
CMD ["./events-api"]
