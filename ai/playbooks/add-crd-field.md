# Add CRD Field

## When To Use

Use this playbook when adding or changing a CRD field in `katalyst-api`.

## Steps

1. Edit the source type in `pkg/apis/<group>/<version>/`.
2. Keep the change additive within the current version. If not possible, create a new version.
3. Keep marker style consistent with neighboring fields.
4. Run `make generate`.
5. Confirm the expected generated surfaces moved with the change, especially `config/crd/bases/*.yaml` and any affected files under `pkg/client/`.
6. Add or update focused tests near the type package, typically next to the touched API version package.
7. Verify `go run ./cmd/inlinecheck` and `go test ./...`.
8. Note downstream bump expectations in PR context when the new field must be consumed by `katalyst-core` or later repos.

## Useful Anchors

- Source types usually live under `pkg/apis/<group>/<version>/types*.go`.
- Generated CRDs land under `config/crd/bases/`.
- Generated clients and deep-copy outputs usually move under `pkg/client/` and `zz_generated*.go`.

## Watch For

- Editing `config/crd/bases/*.yaml` directly
- Renaming a JSON tag in place
- Forgetting generated outputs in the same commit
