# Change Proto Or Constants

## When To Use

Use this playbook for shared constants and plugin protocol updates.

## If You Are Changing Shared Constants

1. Edit the source file under `pkg/consts/`.
2. Avoid duplicating existing constants under new names.
3. Run `make fmt` and `go test ./...`.
4. Document downstream follow-up when `katalyst-core` or downstream repos must consume the new constant.

## If You Are Changing Plugin Protocols

1. Edit the source file under `pkg/protocol/**/api.proto`.
2. Keep protobuf field numbers append-only and optional-compatible.
3. Ensure external `protoc` and `protoc-gen-gogo` executables are available on `PATH`.
4. Run `make generate-pb` when only protobuf outputs need regeneration, or `make generate` when broader generated surfaces also changed.
5. Run `make fmt` and `go test ./...`.
6. Document downstream follow-up when `katalyst-core` or downstream repos must consume the new protocol field.

## Watch For

- Reusing a retired protobuf tag
- Adding a required field to an existing protocol
- Treating a downstream-only constant as if it belonged in this repo
