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
5. note any likely downstream version-bump follow-up for `katalyst-core` or later repos

## Required Special Cases

- CRD type and marker changes require `make generate` (narrower targets `make generate-go` and `make generate-manifests` are permitted by the canonical workflow)
- protobuf changes require `make generate` (narrower target `make generate-pb` is permitted by the canonical workflow)
- shared constant changes do not always need generation, but they can still require downstream consumer updates
- generated files such as `zz_generated_*.go`, `*.pb.go`, `pkg/client/**`, and `config/crd/bases/*.yaml` should not be hand-edited

## Output Contract

Return:
1. the touched source-of-truth area
2. the recommended generation command or commands
3. the expected generated outputs to review
4. whether downstream version bumps or follow-up repo changes are likely
5. any tool or environment gaps that prevent generation
