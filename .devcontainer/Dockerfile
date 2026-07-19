FROM mcr.microsoft.com/devcontainers/go:1-bookworm

# sqlite3 CLI for poking at abyss.db during development
RUN apt-get update && apt-get install -y --no-install-recommends \
    sqlite3 \
    && rm -rf /var/lib/apt/lists/*

# just (task runner used by the justfile)
RUN curl --proto '=https' --tlsv1.2 -sSf https://just.systems/install.sh \
    | bash -s -- --to /usr/local/bin

# bun (frontend package manager / dev server / build)
ENV BUN_INSTALL=/usr/local/bun
ENV PATH=${BUN_INSTALL}/bin:${PATH}
RUN mkdir -p ${BUN_INSTALL} && chmod 755 ${BUN_INSTALL} \
    && curl -fsSL https://bun.sh/install | bash \
    && chmod -R a+rX ${BUN_INSTALL}

# air (backend hot reload), mockery (mock generation), and golangci-lint v2
# (matches .golangci.yml at the repo root), installed for all users
RUN GOBIN=/usr/local/bin go install github.com/air-verse/air@latest \
    && GOBIN=/usr/local/bin go install github.com/vektra/mockery/v3@latest \
    && GOBIN=/usr/local/bin go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@latest
