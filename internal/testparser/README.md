# internal/testparser

This package is reserved for SQL parsing/validation helpers used in **tests only**.

Planned approach (per spec): use `pingcap/tidb` under a build tag to avoid any
runtime dependency.

Example future layout:

- `parser_tidb_test.go` with:
  - `//go:build testparser && !integration`

By default, this directory contains **no Go files** so it is ignored by
`go test ./...`.
