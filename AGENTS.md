# Repository Guidelines

`katalyst-api` (module `github.com/kubewharf/katalyst-api`) is the schema layer of the Katalyst QoS resource-management stack. It owns shared API surfaces consumed by other repos: CRD Go types, generated clients/listers/informers, plugin gRPC contracts, and shared constants.

It ships no runtime binaries; primary consumers are `katalyst-core` and internal downstream adapters.

## Critical Rules

1. **Never hand-edit generated artifacts.** This includes `pkg/client/**`, `zz_generated_*.go`, `*.pb.go`, and `config/crd/bases/*.yaml`.
2. **Regenerate after schema/proto changes.** Run `make generate` after touching CRD types or protos, and commit generated outputs in the same change.
3. **Keep compatibility additive.** Do not remove or rename exported fields/methods/constants in an existing API version. Introduce a new version instead (for example `v1alpha2`).
4. **Avoid reverse or cyclic dependencies.** Never import `github.com/kubewharf/katalyst-core` or known downstream consumers such as `github.com/kubewharf/katalyst-adapter`.
5. **No CRD cross-group imports.** Do not import one API group directly from another (for example `autoscaling` <-> `workload`). Move shared types to a neutral package.
6. **Respect API struct constraints.** No anonymous inline structs in API types; `cmd/inlinecheck` must pass.

## Project Structure

- `pkg/apis/<group>/<version>/`: CRD Go types and kubebuilder markers.
- `pkg/client/{clientset,informers,listers}/`: generated clients; read-only.
- `pkg/consts/`: shared QoS/annotation/resource constants.
- `pkg/protocol/{evictionplugin,reporterplugin}/v1alpha1/`: protobuf and gRPC contracts.
- `pkg/plugins/{registration,skeleton}/`: plugin registration and skeletons.
- `config/crd/bases/`: generated CRD manifests.
- `cmd/inlinecheck/`: inline-struct policy checker.
- `hack/`: codegen scripts and boilerplate headers.

## Build And Validation

### Local Quick Checks

- `make fmt` - run `go fmt ./...`.
- `go run ./cmd/inlinecheck` - verify no inline anonymous structs in API types.

### CI / Environment-Specific Checks

- `make generate` - run all generators (`generate-manifests`, `generate-go`, `generate-pb`).
- `make generate-go` - regenerate deepcopy/client/informer/lister code.
- `make generate-manifests` - regenerate CRD YAML via `controller-gen`.
- `make generate-pb` - regenerate protobuf outputs.
- `go test ./...` - run unit tests.

Use the environment-dependent commands when you are working on schema/proto changes or preparing a change for CI/review, but do not assume they work in every fresh local environment without additional toolchain setup. Some of these commands may still be required before merge or in CI.

## Coding Rules

- Go baseline: `go 1.17` (follow `go.mod`).
- Use `gofmt` formatting and tabs.
- Keep package names lower-case and short.
- Keep kubebuilder marker style consistent with neighboring code (field-level vs struct-level).
- Add Apache-2.0 headers to new `.go` files using `hack/boilerplate.go.txt`.

## Dependency Rules

- `pkg/apis` must not import generated client packages from `pkg/client/**`.
- `pkg/consts` must not import `pkg/apis`.
- Avoid broad controller-runtime usage; marker-related usage is acceptable where already established.
- `klog` is present in existing plugin code; do not introduce new logging dependencies into schema-only packages unless there is a clear existing pattern in that package.

## Testing

- Place tests next to implementation (`foo.go` and `foo_test.go`).
- Prefer table-driven tests; each `t.Run` should call `t.Parallel()`.
- For every new CRD field, add a deepcopy round-trip test.
- For every new CRD field, add a JSON marshal/unmarshal round-trip test when validation/serialization behavior is relevant.

Example JSON round-trip pattern:

```go
func TestServiceProfileDescriptor_JSON(t *testing.T) {
	t.Parallel()
	raw := []byte(`{"apiVersion":"workload.katalyst.kubewharf.io/v1alpha1","kind":"ServiceProfileDescriptor",...}`)
	var got v1alpha1.ServiceProfileDescriptor
	require.NoError(t, json.Unmarshal(raw, &got))
	out, err := json.Marshal(&got)
	require.NoError(t, err)
	require.JSONEq(t, string(raw), string(out))
}
```

## Error Handling

- Prefer wrapped errors from helpers: `fmt.Errorf("<subject>: <reason>: %w", err)`.
- Plugin-protocol errors should follow Kubernetes `metav1.Status` style:
  - `Reason`: `MixedCase`.
  - `Message`: sentence starting with an upper-case letter.
- Do not `panic` in plugin skeletons; return errors.

## Cross-Repo Change Playbook

`katalyst-api` is the first stop of a cross-repo change. Land shared API, constant, and protocol changes here before wiring them into downstream repositories.

### Add A CRD Field

1. Edit `pkg/apis/<group>/<version>/types.go` and markers.
2. Run `make generate`.
3. Add/update deepcopy and JSON round-trip tests.
4. Keep changes additive within the version.
5. Document downstream bump expectations (`katalyst-core`, then internal downstream adapters) in PR notes when relevant.

### Change Plugin Protocol

1. Edit `pkg/protocol/**/*.proto`.
2. Keep protobuf tags additive; never reuse retired tags.
3. Run `make generate-pb` (or `make generate`).
4. Note downstream consumer impact in PR notes.

## Review Checklist

1. Critical rules above are still satisfied.
2. Kubebuilder markers are correct and stylistically consistent.
3. New/changed fields include required tests.
4. `go run ./cmd/inlinecheck` passes.
5. `go test ./...` passes when the required environment/tooling is available.

## Common Pitfalls

- Cross-importing API groups.
- Reusing protobuf field numbers.
- Duplicating constants instead of reusing canonical keys in `pkg/consts`.
