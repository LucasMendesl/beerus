FROM golang:1.23-alpine3.21 AS builder

ARG VERSION="unknown"
ARG COMMIT_SHA="unknown"

WORKDIR ${GOPATH}/src/github.com/lucasmendesl/beerus

COPY . .

RUN go mod download && \
    go build -trimpath -ldflags="-s -w \
    -X 'github.com/lucasmendesl/beerus/version.Version=$VERSION' \
    -X 'github.com/lucasmendesl/beerus/version.ReleaseDate=$(date +%Y-%m-%dT%H:%M:%SZ)' \
    -X 'github.com/lucasmendesl/beerus/version.GitCommit=$COMMIT_SHA' \
    " -o /go/bin/beerus .

FROM alpine:3.21

ARG VERSION="unknown"
ARG COMMIT_SHA="unknown"

LABEL \
    com.github.lucasmendesl.beerus.service="true" \
    org.opencontainers.image.source="https://github.com/lucasmendesl/beerus" \
    org.opencontainers.image.version=$VERSION \
    org.opencontainers.image.revision=$COMMIT_SHA \
    org.opencontainers.image.title="beerus" \
    org.opencontainers.image.description="Beerus - Docker image cleaner" \
    org.opencontainers.image.url="https://github.com/lucasmendesl/beerus" \
    org.opencontainers.image.documentation="https://github.com/lucasmendesl/beerus/blob/main/README.md" \
    org.opencontainers.image.authors="mendes.lucas9498@gmail.com"

RUN apk --no-cache add ca-certificates

COPY --from=builder /go/bin/beerus /usr/bin/beerus

ENTRYPOINT ["/usr/bin/beerus"]
