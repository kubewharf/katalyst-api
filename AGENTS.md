# Repository Guidelines

## 1. Overview

`katalyst-api` (module `github.com/kubewharf/katalyst-api`) is the schema layer of the Katalyst QoS resource-management stack. It owns shared API surfaces consumed by other repos: CRD Go types, generated clients/listers/informers, plugin gRPC contracts, and shared constants.

It ships no runtime binaries; primary consumers are `katalyst-core` and internal downstream repos.

## 2. Critical Rules

1. **Never hand-edit generated artifacts.** This includes `pkg/client/**`, `zz_generated_*.go`, `*.pb.go`, and `config/crd/bases/*.yaml`.
2. **Regenerate after schema/proto changes.** Run `make generate` after touching CRD types or protos, and commit generated outputs in the same change.
3. **Keep compatibility additive.** Do not remove or rename exported fields/methods/constants in an existing API version. Introduce a new version instead (for example `v1alpha2`).
4. **Avoid reverse or cyclic dependencies.** Never import `github.com/kubewharf/katalyst-core` or any internal downstream repo.
5. **No CRD cross-group imports.** Do not import one API group directly from another (for example `autoscaling` <-> `workload`). Move shared types to a neutral package.
6. **Respect API struct constraints.** No anonymous inline structs in API types; `cmd/inlinecheck` must pass.

## 3. Project Structure

- `pkg/apis/<group>/<version>/`: CRD Go types and kubebuilder markers.
- `pkg/client/{clientset,informers,listers}/`: generated clients; read-only.
- `pkg/consts/`: shared QoS/annotation/resource constants.
- `pkg/protocol/{evictionplugin,reporterplugin}/v1alpha1/`: protobuf and gRPC contracts.
- `pkg/plugins/{registration,skeleton}/`: plugin registration and skeletons.
- `config/crd/bases/`: generated CRD manifests.
- `cmd/inlinecheck/`: inline-struct policy checker.
- `hack/`: codegen scripts and boilerplate headers.

## 4. Build & Validation

### Local Gate

Before pushing, run:

1. `make fmt` — run `go fmt ./...`.
2. `go run ./cmd/inlinecheck` — verify no inline anonymous structs in API types.
3. `go test ./...` — run unit tests.

### How to Build

This repo ships no runtime binaries. When you change CRD types or protos, regenerate outputs:

- `make generate` — run all generators (`generate-manifests`, `generate-go`, `generate-pb`).
- `make generate-go` — regenerate deepcopy/client/informer/lister code.
- `make generate-manifests` — regenerate CRD YAML via `controller-gen`.
- `make generate-pb` — regenerate protobuf outputs.

### CI Notes

CI (`.github/workflows/ci.yml`) enforces `go mod tidy`, `make fmt` (gofmt check), and `go run cmd/inlinecheck/main.go`. Generate targets and `go test ./...` are not currently run by CI, so run them locally before pushing. Some generation steps need extra toolchain (`controller-gen`, `protoc`) which may not be present in every fresh local environment; if a step fails only because of missing tooling, note it in the PR.

## 5. Coding & Editing Rules

- Go toolchain: `go.mod` is the single source of truth. This module currently targets `go 1.17`; always follow whatever `go.mod` declares.
- Use `gofmt` formatting and tabs.
- Keep package names lower-case and short.
- Keep kubebuilder marker style consistent with neighboring code (field-level vs struct-level).
- Add Apache-2.0 headers to new `.go` files using `hack/boilerplate.go.txt`.

## 6. Import & Dependency Rules

- `pkg/apis` must not import generated client packages from `pkg/client/**`.
- `pkg/consts` must not import `pkg/apis`.
- Avoid broad controller-runtime usage; marker-related usage is acceptable where already established.
- `klog` is present in existing plugin code; do not introduce new logging dependencies into schema-only packages unless there is a clear existing pattern in that package.
- Reverse imports into `katalyst-core` or any internal downstream repo are forbidden (see section 2, rule 4).

## 7. Generated & Vendored Artifacts

- Generated files owned by this repo (do not hand-edit — see section 2, rule 1): `pkg/client/**`, `zz_generated_*.go`, `*.pb.go`, `config/crd/bases/*.yaml`. Regenerate via `make generate`.
- This repo does not vendor code from downstream repos. Its schema, constants, and plugin protocols are the source of truth for the entire stack; downstream repos consume them via module pin.
- Downstream copies of CRDs in internal deployment repos must remain byte-identical to `config/crd/bases/*.yaml` here.

## 8. Testing

- Place tests next to implementation (`foo.go` and `foo_test.go`).
- Prefer table-driven tests; each `t.Run` must call `t.Parallel()`.
- For every new CRD field, add a deepcopy round-trip test.
- For every new CRD field, add a JSON marshal/unmarshal round-trip test when validation/serialization behavior is relevant.

Example JSON round-trip pattern:

```go
func TestServiceProfileDescriptor_JSON(t *testing.T) {
    t.Parallel()
    raw := []byte(`{"apiVersion":"workload.katalyst.kubewharf.io/v1alpha1","kind":"ServiceProfileDescriptor",...}`)
    var got v1alpha1.ServiceProfileDescriptor
    require.NoError(t, json.Unmarshal(raw, &got))
    out, err := json.Marshal(&got)
    require.NoError(t, err)
    require.JSONEq(t, string(raw), string(out))
}
```

## 9. Logging, Errors & Metrics

- This repo is schema-only. Do not introduce new logging or metrics in schema packages.
- Existing plugin skeletons may use `klog`; follow the pattern already in that package.
- Prefer wrapped errors from helpers: `fmt.Errorf("<subject>: <reason>: %w", err)`.
- Plugin-protocol errors should follow Kubernetes `metav1.Status` style: `Reason` is `MixedCase`; `Message` is a sentence starting with an upper-case letter.
- Do not `panic` in plugin skeletons; return errors.

## 10. Versioning & Compatibility

Compatibility of exported API surface is stated in section 2, rule 3. Additional operational rules for this repo:

- Introduce a new API version (for example `v1alpha2`) for incompatible changes.
- Protobuf tags are append-only; never reuse retired tags and never add required fields.
- Document downstream bump expectations (`katalyst-core`, then internal downstream repos) in PR notes when relevant.

## 11. Commit & Release

- Conventional Commits are not enforced by tooling in this repo; follow upstream `kubewharf` project conventions.
- Schema/proto changes should note downstream consumer impact (`katalyst-core` and internal downstream repos) in the PR description. The "same commit" rule for generated outputs is stated in section 2, rule 2.

## 12. Cross-Repo Change Playbooks

### Stack Layering & Dependency Direction

The Katalyst stack is layered strictly from schema to rollout:

1. `katalyst-api` — shared schema, constants, and plugin protocols.
2. `katalyst-core` — upstream runtime binaries (agent, controller, scheduler, webhook, metric).
3. Internal downstream repos — internal-only repos vendor the upstream releases and produce deployment artifacts. Their contents and workflows are out of scope for this file.

Dependency direction is one-way: internal downstream <- `katalyst-core` <- `katalyst-api`. Never introduce reverse or cyclic imports. Cross-repo changes always land in `katalyst-api` first, then `katalyst-core`, then internal downstream.

### Playbook Skeleton (shared step IDs)

Cross-repo changes follow a stable skeleton with numbered step IDs. Each repo below only lists its own delta.

- **CRD-1** — Design and land the CRD field in `katalyst-api`: edit `pkg/apis/<group>/<version>/types.go`, run `make generate`, add deepcopy + JSON round-trip tests. Keep the change additive.
- **CRD-2** — Bump the `katalyst-api` pin in `katalyst-core` and wire the field into the relevant agent/controller/scheduler/webhook logic; add focused unit tests.
- **CRD-3 / CRD-4** — Downstream integration is handled by internal repos; details are out of scope here.
- **QRM-1** — Land any new resource constants or plugin-protocol changes in `katalyst-api` (`pkg/consts/`, `pkg/protocol/`). Keep protobuf tags additive.
- **QRM-2** — Implement the upstream plugin in `katalyst-core/pkg/agent/qrm-plugins/<resource>/`, register it in `cmd/katalyst-agent/app/enableagents.go`, and add config wiring under `pkg/config/agent/qrm/`.
- **QRM-3 / QRM-4** — Downstream integration is handled by internal repos; details are out of scope here.
- **CTRL-1** — Ensure the CRD is released from `katalyst-api`.
- **CTRL-2** — Add the reconciler under `katalyst-core/pkg/controller/<name>/` and register it in `cmd/katalyst-controller/app/enablecontrollers.go`.
- **CTRL-3 / CTRL-4** — Downstream integration is handled by internal repos; details are out of scope here.

### katalyst-api delta

`katalyst-api` owns **CRD-1**, **QRM-1**, and **CTRL-1**. Concretely:

- **CRD-1**: edit `pkg/apis/<group>/<version>/types.go` and markers; run `make generate`; add deepcopy and JSON round-trip tests; keep changes additive within the version.
- **QRM-1**: if the new plugin needs new resource constants or plugin-protocol changes, land them in `pkg/consts/` and `pkg/protocol/**/*.proto` first. Keep protobuf tags additive; never reuse retired tags.
- **CTRL-1**: land the CRD types under `pkg/apis/<group>/<version>/`, run `make generate`, then let downstream bump.

## 13. Code Review Checklist

1. Critical rules in section 2 are still satisfied.
2. Kubebuilder markers are correct and stylistically consistent.
3. New/changed fields include deepcopy and JSON round-trip tests where relevant.
4. `go run ./cmd/inlinecheck` passes.
5. `go test ./...` passes when the required environment/tooling is available.
6. Generated files are regenerated and committed alongside source changes.
7. No reverse or cyclic imports were introduced.

## 14. Common Pitfalls

Concrete symptoms that signal a critical-rule violation:

- Reusing a protobuf field number after it was retired (violates section 2, rule 3).
- Renaming a `json` tag or adding a `required` proto field inside an existing version (violates section 2, rule 3).
- Duplicating constants locally instead of reusing canonical keys in `pkg/consts` (rules of thumb, not a rule 2 violation, but discovered frequently in review).
- Editing `zz_generated_*.go`, `*.pb.go`, or `config/crd/bases/*.yaml` by hand to work around a codegen problem (violates section 2, rule 1).

## 15. Appendix

### Quick Reference

| I want to... | Edit here |
| --- | --- |
| Add a CRD field | `pkg/apis/<group>/<version>/types.go` |
| Change a plugin protocol | `pkg/protocol/**/*.proto` |
| Add or update a shared constant | `pkg/consts/<domain>.go` |
| Regenerate everything | `make generate` |
| Verify no inline anonymous structs | `go run ./cmd/inlinecheck` |
