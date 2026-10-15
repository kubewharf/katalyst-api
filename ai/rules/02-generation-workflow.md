# Generation Workflow

Most `katalyst-api` changes follow the pattern: edit source of truth, regenerate, then commit both.

## When To Use

Use this rule when a change touches CRD types, protobuf contracts, shared constants that feed generated outputs, or any generated surface in this repo.

## Sources Of Truth

- CRD Go types live under `pkg/apis/<group>/<version>/`.
- Protobuf contracts live under `pkg/protocol/**/api.proto`.
- Shared constants live under `pkg/consts/`.

## Generated Outputs

- `pkg/client/**`
- `zz_generated_*.go`
- `*.pb.go`
- `config/crd/bases/*.yaml`

## Commands

- `make generate` regenerates all supported outputs.
- `make generate-go` regenerates deepcopy and client code.
- `make generate-manifests` regenerates CRD YAML.
- `make generate-pb` regenerates protobuf outputs.

## Keep In Mind

If a generated file changes unexpectedly, fix the source-of-truth edit or generator invocation instead of patching the output by hand.
