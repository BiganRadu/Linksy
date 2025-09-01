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

WORKDIR /app
COPY --from=builder /app/bazel-bin/backend_/backend ./server

EXPOSE 8080
CMD ["./server"]