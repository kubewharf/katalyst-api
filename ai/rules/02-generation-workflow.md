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
- `zz_generated*.go`
- `*.pb.go`
- `config/crd/bases/*.yaml`

## Commands

- `make generate` regenerates all supported outputs.
- `make crd` regenerates CRD manifests, Go clients, and deepcopy output for CRD type or marker changes.
- `make generate-go` regenerates deepcopy and client code.
- `make generate-manifests` regenerates CRD YAML.
- `make generate-pb` regenerates protobuf outputs.

## Keep In Mind

- CodeRabbit intentionally keeps its default `*.pb.go` exclusion. Review the owning `.proto` source, run `make generate-pb`, and inspect the generated diff locally instead of relying on AI review of generator output.
- Protobuf generation requires external `protoc` and `protoc-gen-gogo` executables on `PATH`; the Makefile does not install or pin them.
- The Makefile resolves and exports `GOPATH` to its generation scripts, so callers do not need to set it separately.
- After Go generation, verify that the expected outputs changed in the current checkout. If they did not, fix the generator invocation before proceeding instead of copying or hand-editing generated files.
- If a generated file changes unexpectedly, fix the source-of-truth edit or generator invocation instead of patching the output by hand.
