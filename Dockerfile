# Multi-stage build for ZiXiao Git Server

# Stage 1: Build C++ library
FROM ubuntu:22.04 AS cpp-builder

RUN apt-get update && apt-get install -y \
    build-essential \
    cmake \
    libssl-dev \
    zlib1g-dev \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /build

COPY git-core ./git-core
COPY Makefile .

RUN make build-cpp

# Stage 2: Build frontend
FROM node:20-alpine AS frontend-builder

WORKDIR /app

COPY frontend/package*.json ./frontend/
RUN cd frontend && npm ci

COPY frontend ./frontend
RUN cd frontend && npm run build

# Stage 3: Build Go server
FROM golang:1.21-bullseye AS go-builder

RUN apt-get update && apt-get install -y \
    build-essential \
    libssl-dev \
    zlib1g-dev \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /app

# Copy C++ library from previous stage
COPY --from=cpp-builder /build/git-core ./git-core

# Copy Go source
COPY go.mod go.sum ./
RUN go mod download

COPY cmd ./cmd
COPY internal ./internal
COPY pkg ./pkg
COPY Makefile .

# Build Go binary
RUN CGO_ENABLED=1 make build-go

# Stage 4: Final runtime image
FROM ubuntu:22.04

RUN apt-get update && apt-get install -y \
    libssl3 \
    zlib1g \
    ca-certificates \
    && rm -rf /var/lib/apt/lists/*

# Create non-root user
RUN useradd -m -u 1000 -s /bin/bash gitserver

WORKDIR /app

# Copy built binary and C++ library
COPY --from=go-builder /app/bin/zixiao-git-server /app/
COPY --from=go-builder /app/git-core/lib /app/git-core/lib

# Copy frontend build
COPY --from=frontend-builder /app/web/dist /app/web/dist

# Copy configuration
COPY configs /app/configs

# Create necessary directories
RUN mkdir -p /app/data/repositories /app/logs && \
    chown -R gitserver:gitserver /app

USER gitserver

# Set library path
ENV LD_LIBRARY_PATH=/app/git-core/lib:$LD_LIBRARY_PATH

EXPOSE 8080

HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD curl -f http://localhost:8080/api/v1/health || exit 1

CMD ["./zixiao-git-server", "-config", "./configs/server.yaml"]
