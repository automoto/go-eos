# Linux development/test image for go-eos. Mirrors the CI environment
# (Debian-based Go 1.26) and adds the C toolchain required by cgo and
# `make lint-c`.
#
# The EOS C SDK is NOT baked into this image — it's bind-mounted from
# ./static/ at runtime via docker-compose.yml, since the SDK cannot be
# redistributed and the cgo directives in eos/internal/cbinding/cgo.go
# use ${SRCDIR}-relative paths that must resolve against the live repo.
FROM golang:1.26.2-bookworm

RUN apt-get update \
    && apt-get install -y --no-install-recommends \
        build-essential \
        ca-certificates \
        curl \
        gnupg \
        git \
        make \
    && curl -fsSL https://apt.llvm.org/llvm-snapshot.gpg.key | \
       gpg --dearmor -o /usr/share/keyrings/llvm-archive-keyring.gpg \
    && echo "deb [signed-by=/usr/share/keyrings/llvm-archive-keyring.gpg] http://apt.llvm.org/bookworm/ llvm-toolchain-bookworm-18 main" \
       > /etc/apt/sources.list.d/llvm-18.list \
    && apt-get update \
    && apt-get install -y --no-install-recommends clang-format-18 \
    && ln -sf /usr/bin/clang-format-18 /usr/bin/clang-format \
    && rm -rf /var/lib/apt/lists/*

# golangci-lint v2.x to match `.golangci.yml` (version: "2"). CI uses
# `latest` so we default to that; override with
# --build-arg GOLANGCI_LINT_VERSION=v2.1.6 for a reproducible pin.
ARG GOLANGCI_LINT_VERSION=latest
RUN curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh \
    | sh -s -- -b /usr/local/bin ${GOLANGCI_LINT_VERSION}

WORKDIR /workspace

# Module + build caches live on named volumes (see docker-compose.yml)
# so re-runs are incremental. CGO_ENABLED is on because the cbinding
# package always needs it for non-stub builds.
ENV GOMODCACHE=/go/pkg/mod \
    GOCACHE=/root/.cache/go-build \
    CGO_ENABLED=1

CMD ["bash"]
