#!/bin/bash

#Update the bazel repository
bazel run //:gazelle -- update-repos -from_file=go.mod -to_macro=deps.bzl%go_dependencies

#Update the BUILD.bazel files
bazel run //:gazelle

# Build the Go program using Bazel
bazel build //:backend