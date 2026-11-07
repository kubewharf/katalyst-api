# Katalyst API Skill Specs

Use the skill specs in this directory for repeated `katalyst-api` workflows that benefit from structured generation, compatibility, and validation answers.

These files are repo-local skill specs. Treat them as local source-of-truth documents unless the current harness explicitly exposes them as invocable skills.

## Scope

- Use `rules/` for durable schema, generation, and validation constraints.
- Use `playbooks/` for common CRD and protobuf change flows.
- Use `skills/` for repeated selection tasks where an agent should explain the correct commands, source-of-truth edits, or downstream impact.

## How To Use

- If the current harness exposes a matching skill, invoke it directly.
- If not, read the relevant `SKILL.md` as a structured decision guide and return the same output contract manually.
- Prefer the generation selector when the main question is which outputs regenerate.
- Prefer the validation selector when the main question is which checks to run.

## Current Skill Specs

- `katalyst-api-generation-selector`
  Choose the right generation commands and expected generated outputs for a schema or protobuf change.
- `katalyst-api-validation-selector`
  Choose the right validation checks for touched API, protobuf, or generated files.

## Maintenance

Keep these files focused on trigger conditions, selection logic, and output contracts. Do not duplicate detailed schema rules or implementation steps that already live in `rules/` or `playbooks/`.
