# Schema Compatibility

`katalyst-api` is the source of truth for shared API surfaces. Favor additive evolution.

## When To Use

Use this rule when changing CRD types, shared constants, protobuf contracts, or any other exported API surface in `katalyst-api`.

## Guidance

- Do not remove, rename, or repurpose exported fields or constants in an existing API version.
- Introduce a new version when a change is not backward-compatible.
- Keep JSON tags stable within a released version.
- Keep protobuf tags append-only; never reuse retired tags and never add required fields.
- Avoid cyclic API-group dependencies. Reuse established cross-group types when appropriate, and move broadly shared types to a neutral package.
- Do not hand-edit generated outputs such as `zz_generated*.go`, `*.pb.go`, or `config/crd/bases/*.yaml`.

## Review Questions

- Is the change additive for current consumers?
- Does a downstream repo need a version bump after this change?
- Is the source-of-truth file edited rather than a generated artifact?
