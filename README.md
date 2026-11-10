# Cells Runtime

`cells-runtime` provides the reusable service runtime extracted from Pydio Cells.
It contains lifecycle orchestration, registries, brokers, configuration stores,
storage drivers, HTTP and gRPC servers, telemetry, and the infrastructure
protobuf contracts needed by those components.

Applications supply product policy through explicit integration hooks. The
runtime does not import the Cells application module.

## Development

The module uses the Go version and toolchain declared in `go.mod`.

```sh
GOWORK=off go build ./...
GOWORK=off go test ./...
```

Some filesystem notification tests require native host event access, and the
NTP test requires network access. Run those checks on a host with the required
facilities.

Protocol definitions live under `proto/`. Generated Go bindings are committed;
edit schemas and generation configuration before regenerating them.

The runtime owns its custom protobuf generators under `tools/`. Generate one
contract group with `./proto/buf.sh service` (or another directory listed in
`proto/buf.work.yaml`). Generator and protobuf dependency versions are pinned
by this module; the script does not generate external SDKs.
