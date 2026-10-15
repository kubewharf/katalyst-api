# Validation And CI

Use local validation that matches the change, not only the narrow CI minimum.

## When To Use

Use this rule when deciding which checks to run for schema, protobuf, generated-output, or validation-only changes in `katalyst-api`.

## Local Checks

- `make fmt`
- `go run ./cmd/inlinecheck`
- `go test ./...`

## Generation Checks

Run `make generate` whenever CRD types, markers, protobuf contracts, or shared generated surfaces change.

## Keep In Mind

The checked-in CI workflow mainly enforces:

- `go mod tidy`
- `make fmt`
- `go run cmd/inlinecheck/main.go`

CI does not fully replace local schema, generation, and test verification. If generation cannot run because required tools such as `controller-gen` or `protoc` are missing, say so explicitly in the final handoff.
