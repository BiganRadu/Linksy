# Stage 1: Build with Bazel
FROM l.gcr.io/google/bazel:latest AS builder

# Install newer Bazel (6.5.0)
RUN curl -L https://github.com/bazelbuild/bazel/releases/download/6.5.0/bazel-6.5.0-linux-x86_64 -o /usr/local/bin/bazel \
    && chmod +x /usr/local/bin/bazel

WORKDIR /app
COPY . .

# Build your binary with Bazel
RUN bazel build //:backend

# Stage 2: Minimal runtime image
FROM debian:bookworm-slim

# ✅ Install CA certificates so TLS connections (e.g. MongoDB Atlas) work
RUN apt-get update && apt-get install -y --no-install-recommends ca-certificates && \
    rm -rf /var/lib/apt/lists/*

WORKDIR /app
COPY --from=builder /app/bazel-bin/backend_/backend ./server

RUN mkdir -p /app/pictures

EXPOSE 8080
CMD ["./server"]