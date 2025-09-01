# Stage 1: Build with Bazel
FROM l.gcr.io/google/bazel:latest AS builder

WORKDIR /app
COPY . .

# Build your binary with Bazel
RUN bazel build //:backend

# Stage 2: Minimal runtime image
FROM debian:bookworm-slim

WORKDIR /app
COPY --from=builder /app/bazel-bin/backend ./server

EXPOSE 8080
CMD ["./server"]