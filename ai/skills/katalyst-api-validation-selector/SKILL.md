---
name: "katalyst-api-validation-selector"
description: "Selects the right validation steps for touched katalyst-api schema, protobuf, generation, or API-compatibility changes."
---

# Katalyst API Validation Selector

Use this skill to choose the smallest correct validation set for a `katalyst-api` change.

Invoke this skill when:
- the user asks what to run after editing `katalyst-api`
- touched files span source-of-truth and generated outputs
- the change affects CRD schemas, protobuf contracts, or compatibility-sensitive API surfaces
- you need to explain whether generation is part of validation for this edit

Do not prefer this skill when:
- the task is read-only and no validation guidance is needed
- the change is outside `katalyst-api`
- a single obvious command already answers the question

## Source Of Truth

Use these repo-local files as the canonical guidance:

- `AGENTS.md`
- `ai/rules/01-schema-compatibility.md`
- `ai/rules/02-generation-workflow.md`
- `ai/rules/03-validation-and-ci.md`

Use the root `Makefile` for exact command behavior.

## Selection Rules

When using this skill:
1. identify whether the change affects schema, protobuf, shared constants, or generated files
2. decide whether regeneration is required before validation
3. choose the smallest local check set that still matches repo expectations
4. verify that generated Go output landed in the current checkout
5. call out CI-equivalent checks separately from stronger local verification
6. state any missing generator tooling or environment limits

## Required Special Cases

- CRD type, kubebuilder marker, and protobuf changes require regeneration before review or test validation. Use `make generate-pb` for protobuf-only changes; use both `make generate-go` and `make generate-manifests` (or `make crd`) when a CRD change affects Go output and CRD YAML; use `make generate` when all generated surfaces need refresh (see `02-generation-workflow.md`)
- protobuf generation requires external `protoc` and `protoc-gen-gogo` executables on `PATH`
- after Go generation, verify that the expected outputs changed in the current checkout
- `make fmt`, `go run ./cmd/inlinecheck`, and `go test ./...` are the normal local validation set
- checked-in CI mainly enforces `go mod tidy`, `make fmt`, and `go run cmd/inlinecheck/main.go`
- compatibility-sensitive exported API changes should explicitly call out downstream consumer impact even when tests pass

## Output Contract

Return:
1. the touched `katalyst-api` area
2. the recommended validation commands
3. whether regeneration is required first
4. whether generated output landed in the intended checkout
5. why each check matters for this change
6. any CI, tooling, or downstream follow-up gaps
