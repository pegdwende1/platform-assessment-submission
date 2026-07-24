# ============================================
# Multi-stage build for minimal production image
# ============================================

# Stage 1: Build
FROM golang:1.22-alpine AS builder

ARG VERSION=dev

WORKDIR /app
COPY app/go.mod ./
RUN go mod download

COPY app/ ./
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w -X main.Version=${VERSION}" -o /cicd-demo .

# Stage 2: Runtime (distroless for security)
FROM gcr.io/distroless/static:nonroot

COPY --from=builder /cicd-demo /cicd-demo

USER nonroot:nonroot
EXPOSE 8080

ENTRYPOINT ["/cicd-demo"]
