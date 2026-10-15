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

1. **Never hand-edit generated artifacts.** This includes `pkg/client/**`, `zz_generated*.go`, `*.pb.go`, and `config/crd/bases/*.yaml`.
2. **Regenerate after schema/proto changes.** Run `make crd` for CRD type or marker changes and `make generate-pb` for protobuf-only changes. Reserve `make generate` for changes that require all generated surfaces. Commit generated outputs in the same change.
3. **Keep compatibility additive.** Do not remove or rename exported fields/methods/constants in an existing API version. Introduce a new version instead (for example `v1alpha2`).
4. **Avoid reverse or cyclic dependencies.** Never import `github.com/kubewharf/katalyst-core` or any internal downstream repo.
5. **Avoid cyclic API-group dependencies.** Reuse established cross-group types when appropriate, and move broadly shared types to a neutral package.
6. **Respect API struct constraints.** Use named types instead of anonymous struct literals in API fields; embedded named types are allowed. Give non-embedded API fields proper JSON tags. Separately, `cmd/inlinecheck` must pass; its tag checks cover only exported, non-embedded fields that resolve to structs after at most one pointer or collection unwrap.

## 3. Project Structure

- `pkg/apis/<group>/<version>/`: CRD Go types and kubebuilder markers.
- `pkg/client/{clientset,informers,listers}/`: generated clients; read-only.
- `pkg/consts/`: shared QoS/annotation/resource constants.
- `pkg/protocol/{evictionplugin,reporterplugin}/v1alpha1/`: protobuf and gRPC contracts.
- `pkg/plugins/{registration,skeleton}/`: plugin registration and skeletons.
- `config/crd/bases/`: generated CRD manifests.
- `cmd/inlinecheck/`: JSON-tag checker for exported, non-embedded fields that resolve to structs after at most one pointer or collection unwrap.
- `hack/`: codegen scripts and boilerplate headers.

## 4. Build & Validation

### Local Gate

Before pushing, run:

1. `make fmt` — run `go fmt ./...`.
2. `go run ./cmd/inlinecheck` — check for non-empty JSON tag names on exported, non-embedded fields that resolve to structs while traversing registered API types. The checker unwraps at most one pointer or collection layer per field; scalar fields and deeper container nesting are not checked.
3. `go test ./...` — run unit tests.

### How to Build

This repo ships no runtime binaries. When you change source schemas, regenerate the outputs they own:

- `make generate-pb` — protobuf-only changes.
- `make crd` — CRD type or marker changes; regenerates manifests plus Go clients/deepcopy output.
- `make generate` — full regeneration when CRD, client, and protobuf surfaces must all be refreshed.

Use `ai/rules/02-generation-workflow.md` or `ai/skills/katalyst-api-generation-selector/SKILL.md` for the detailed generator selection logic and expected outputs.

### CI Notes

CI (`.github/workflows/ci.yml`) enforces `go mod tidy`, runs `make fmt`, checks the resulting formatting diff under `pkg`, and runs `go run cmd/inlinecheck/main.go`. Go files outside `pkg` are formatted by the command but are not included in that CI diff check. Generate targets, generated-output freshness, and `go test ./...` are not currently checked by CI, so verify them locally before pushing. The Makefile builds its pinned `controller-gen` tool and exports its resolved `GOPATH` to generation scripts; protobuf generation still requires external `protoc` and `protoc-gen-gogo` executables on `PATH`. Follow `ai/rules/02-generation-workflow.md` so generated Go files land in the intended checkout. If a step fails only because required external tooling is missing, note it in the PR.

## 5. Coding & Editing Rules

- Go toolchains: the root module targets Go 1.17 through `go.mod`. The nested `hack/tools` module requires Go 1.22 with toolchain 1.22.3 to build the pinned `controller-gen`; follow the module that owns the command being run.
- Use `gofmt` formatting and tabs.
- `make fmt` formats the root module through `go fmt ./...`; it does not traverse the nested `hack/tools` module.
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

- Generated files owned by this repo (do not hand-edit — see section 2, rule 1): `pkg/client/**`, `zz_generated*.go`, `*.pb.go`, `config/crd/bases/*.yaml`. Regenerate via the applicable target described in section 4.
- CodeRabbit intentionally keeps its default `*.pb.go` exclusion. Review the owning `.proto` source and verify protobuf output by running the applicable generator and inspecting the generated diff; the exclusion does not relax the requirement to commit regenerated output.
- This repo does not vendor code from downstream repos. Its API schemas, shared constants, and public external plugin protocols are the source of truth for the surfaces it owns; downstream repos consume them via module pin. Protocols used only inside `katalyst-core` remain core-owned.
- Downstream CRD copies must follow their documented provenance policy for the intended `katalyst-api` release: exact-copy entries remain byte-identical, while documented semantic copies preserve parsed-schema equality.

## 8. Testing

- Place tests next to implementation (`foo.go` and `foo_test.go`).
- Prefer table-driven tests and parallel subtests where they are independent.
- Add focused deepcopy coverage when a CRD change introduces reference-bearing fields or other non-trivial copy behavior.
- For every new CRD field, add a JSON marshal/unmarshal round-trip test when validation/serialization behavior is relevant.

Use `ai/playbooks/add-crd-field.md` for the normal CRD-change test expectations.

## 9. Logging, Errors & Metrics

- This repo is schema-only. Do not introduce new logging or metrics in schema packages.
- Existing plugin skeletons may use `klog`; follow the pattern already in that package.
- Prefer wrapped errors from helpers: `fmt.Errorf("<subject>: <reason>: %w", err)`.
- Keep plugin-protocol errors backward-compatible and follow the conventions already established by the protocol package being changed.
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
- Editing `zz_generated*.go`, `*.pb.go`, or `config/crd/bases/*.yaml` by hand to work around a codegen problem (violates section 2, rule 1).

## 15. Appendix

### Quick Reference

| I want to... | Edit here |
| --- | --- |
| Add a CRD field | `pkg/apis/<group>/<version>/types.go` |
| Change a plugin protocol | `pkg/protocol/**/*.proto` |
| Add or update a shared constant | `pkg/consts/<domain>.go` |
| Regenerate everything | `make generate` |
| Check JSON tags while traversing registered API structs | `go run ./cmd/inlinecheck` |
