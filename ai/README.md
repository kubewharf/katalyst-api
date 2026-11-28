# Katalyst API AI Layer

This directory is the repo-local AI layer for `katalyst-api`.

Use `AGENTS.md` as the top-level entrypoint. Use this directory for stable, repo-specific guidance that should remain useful when this repository is cloned by itself.

## Layout

- `rules/` documents repo-specific constraints that agents should follow repeatedly.
- `playbooks/` documents common change flows such as CRD and protobuf work.
- `skills/` captures repeated generation and validation selection workflows with structured outputs.

## Scope

`katalyst-api` is the schema layer of the Katalyst stack. Prefer guidance here for:

- CRD type changes
- protobuf and plugin protocol changes
- generated client or CRD output regeneration
- API compatibility and downstream impact

## Start Here

- Use `rules/01-schema-compatibility.md` for additive API and ownership boundaries.
- Use `rules/02-generation-workflow.md` when you need to know which generated outputs move with a source change.
- Use `playbooks/add-crd-field.md` for CRD field changes and `playbooks/change-proto-or-constants.md` for shared constants or protobuf changes.
- Use `skills/katalyst-api-generation-selector/SKILL.md` or `skills/katalyst-api-validation-selector/SKILL.md` when the main question is which commands or checks apply.
- Use `skills/README.md` if you need the local skill-spec usage model, especially the live-skill versus local-spec distinction.

## Maintenance

Do not duplicate shared workspace policy here.

Keep files focused on `katalyst-api` specifics, and keep volatile notes out of this local layer.
