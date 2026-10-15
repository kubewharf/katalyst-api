---
name: "katalyst-api-generation-selector"
description: "Selects the right generation commands and expected generated outputs for touched katalyst-api schema, protobuf, or shared constant changes."
---

# Katalyst API Generation Selector

Use this skill to decide which generation flow applies to a `katalyst-api` change and which checked-in outputs should move with it.

Invoke this skill when:
- the user changes CRD Go types, kubebuilder markers, protobuf contracts, or shared constants
- you need to decide between `make generate`, `make generate-go`, `make generate-manifests`, or `make generate-pb`
- you need to explain which generated files should be expected in the diff
- you need to distinguish source-of-truth edits from generated outputs

Do not prefer this skill when:
- the task is outside `katalyst-api`
- the change is read-only and no generation guidance is needed
- the user already knows the exact generator command and only needs a file lookup

## Source Of Truth

Use these repo-local files as the canonical guidance:

- `AGENTS.md`
- `ai/rules/01-schema-compatibility.md`
- `ai/rules/02-generation-workflow.md`
- `ai/playbooks/add-crd-field.md`
- `ai/playbooks/change-proto-or-constants.md`

Use the root `Makefile` for exact generator targets and behavior.

## Selection Rules

When using this skill:
1. identify whether the source-of-truth change is in CRD Go types, protobuf files, or shared constants
2. select the smallest generator command that matches the touched surface
3. list the generated outputs that should be reviewed and committed
4. call out when the user edited a generated file instead of the source-of-truth path
5. verify that the expected generated outputs changed in the current checkout
6. note any likely downstream version-bump follow-up for `katalyst-core` or later repos

## Required Special Cases

- CRD type and marker changes normally require `make crd`, which runs both `make generate-go` and `make generate-manifests`. Use a single narrow target only when the change affects that target's output alone, and reserve `make generate` for changes that also require protobuf regeneration.
- protobuf-only changes require `make generate-pb`; use `make generate` only when broader generated surfaces also need refresh, and verify that external `protoc` and `protoc-gen-gogo` executables are available on `PATH`
- shared constant changes do not always need generation, but they can still require downstream consumer updates
- after Go generation, verify that the expected outputs changed in the current checkout
- generated files such as `zz_generated*.go`, `*.pb.go`, `pkg/client/**`, and `config/crd/bases/*.yaml` should not be hand-edited

## Output Contract

Return:
1. the touched source-of-truth area
2. the recommended generation command or commands
3. the expected generated outputs to review
4. whether the expected generated outputs were produced in the current checkout
5. whether downstream version bumps or follow-up repo changes are likely
6. any tool or environment gaps that prevent generation
