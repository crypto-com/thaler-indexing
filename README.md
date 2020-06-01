# Crypto.com Chain Indexing Service

## Prerequisite

- [Go](https://golang.org/dl/) the programming language

## Build

```bash
go build -o server ./cmd/server/main.go
```

## Lint

#### Prerequisite

- [golangci-lint](https://github.com/golangci/golangci-lint)

```bash
./lint.sh
```

## Test

```bash
./test.sh [--install-dependency]
```

Providing `--install-dependency` will attempt to install test runner [Ginkgo](https://github.com/onsi/ginkgo) if it is not installed before.

## License

[Apache 2.0](./LICENSE)