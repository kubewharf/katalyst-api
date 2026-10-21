# Validation And CI

Use local validation that matches the change, not only the narrow CI minimum.

## When To Use

Use this rule when deciding which checks to run for schema, protobuf, generated-output, or validation-only changes in `katalyst-api`.

## Local Checks

- `make fmt`
- `go run ./cmd/inlinecheck`
- `go test ./...`

`inlinecheck` checks for non-empty JSON tag names only on exported, non-embedded fields that resolve to structs after at most one pointer, slice, array, or map unwrap. It traverses discovered struct types but skips scalar fields and deeper container nesting; review those tags separately. It does not enforce the separate ban on anonymous struct literals.

## Generation Checks

Run the generation target that owns the changed surface:

- `make generate-pb` for protobuf-only changes
- `make generate-go` and `make generate-manifests` (or `make crd`) when a CRD change affects both Go and manifest outputs
- `make generate` when all generated surfaces need to be refreshed

The checked-in CI does not verify generated-output freshness, so inspect and commit the expected generated diff locally.

## Keep In Mind

The checked-in CI workflow mainly enforces:

- `go mod tidy`
- `make fmt`
- `go run cmd/inlinecheck/main.go`

CI does not fully replace local schema, generation, and test verification. Its gofmt diff check stages only `pkg`, so review formatting changes to Go files elsewhere as part of local validation. The Makefile builds `controller-gen` from `hack/tools`; protobuf generation still requires external `protoc` and `protoc-gen-gogo` executables on `PATH`. If generation cannot run because required external tooling is missing, say so explicitly in the final handoff.
