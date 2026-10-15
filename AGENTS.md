# Repository Guidelines

## 1. Overview

`katalyst-api` (module `github.com/kubewharf/katalyst-api`) is the schema layer of the Katalyst QoS resource-management stack. It owns shared API surfaces consumed by other repos: CRD Go types, generated clients/listers/informers, plugin gRPC contracts, and shared constants.

It ships no runtime binaries; primary consumers are `katalyst-core` and internal downstream repos.

## How To Use This File

- Use `AGENTS.md` for repo contract, invariants, and the main local validation gate.
- Use `ai/rules/` for durable repo-specific constraints.
- Use `ai/playbooks/` for common change procedures.
- Use `ai/skills/` for routing, generation, and validation-selection tasks.

## Local AI Layer

Use `ai/README.md` as the repo-local AI index. The local `ai/` tree holds stable `katalyst-api`-specific `rules/`, `playbooks/`, and `skills/` that should remain useful when this repository is cloned by itself. Keep shared cross-repo policy in the workspace-level `ai/` layer, and keep volatile notes out of this repo-local layer.

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

- `make generate` — normal full regeneration path.

Use `ai/rules/02-generation-workflow.md` or `ai/skills/katalyst-api-generation-selector/SKILL.md` for the detailed generator selection logic and expected outputs.

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

Use `ai/playbooks/add-crd-field.md` for the normal CRD-change test expectations.

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

## 12. Repo-Local References

Use the repo-local `ai/` layer for detailed workflow guidance instead of duplicating it here:

- `ai/rules/01-schema-compatibility.md` — additive API evolution, ownership, and compatibility boundaries.
- `ai/rules/02-generation-workflow.md` — source-of-truth files, generated outputs, and the right generator commands.
- `ai/rules/03-validation-and-ci.md` — local validation, CI minimums, and tooling caveats.
- `ai/playbooks/add-crd-field.md` — repo-local flow for CRD field changes.
- `ai/playbooks/change-proto-or-constants.md` — repo-local flow for protobuf and shared-constant changes.
- `ai/skills/katalyst-api-generation-selector/SKILL.md` — choose the right generation command and expected generated surfaces.
- `ai/skills/katalyst-api-validation-selector/SKILL.md` — choose the right validation set for schema, protobuf, and generated-file changes.

For stack-wide work, this repo owns the schema-first step. Land shared API changes here before `katalyst-core` or downstream repos consume them.

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
