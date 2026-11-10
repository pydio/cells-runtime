# Protocol definitions and generated code

This subtree contains protobuf contracts, generated Go bindings, and REST/OpenAPI
artifacts. Read the root guide as well as these instructions.

- Edit `.proto` sources and generator configuration, then regenerate derived
  bindings/specifications. Do not hand-edit `.pb.go` or Swagger output as the
  primary fix. Some Go files here are handwritten helpers; inspect their headers.
- Preserve field numbers, wire types, RPC contracts, and REST mappings. Reserve
  removed field names/numbers rather than reusing them. Consider both API versions
  and SDK consumers when changing shared messages.
- Each protocol directory has its own Buf configuration and generation templates.
  From `common/proto`, run `bash buf.sh tree` for the tree protocol, for example;
  replace `tree` with the affected protocol directory. Inspect its templates and
  optional tag/Python generation before running.
- The script expects `buf` and protoc plugins on its configured PATH. It uses
  `common/proto/bin` and only installs plugins when that directory does not exist;
  an existing directory does not guarantee the required plugins are present.
  Custom generator sources are separate Go modules under root `cmd/`.
- REST generation uses `bash buf.sh rest` from `common/proto`. Read
  `rest/README.md` and `rest/swagger.sh` first: OpenAPI Generator 7.12.0 is
  required by the script. `GENERATE_SDKS_V1` and `GENERATE_SDKS_V2` enable writes
  to external SDK trees; the V1 flow also stages files and replaces generated
  source directories. Set these only for an explicitly scoped SDK update.
- Run `buf lint` from the affected protocol directory. When changing an existing
  contract, run `buf breaking --against <baseline>` with an actual compatible
  baseline for that module; do not invent a remote branch or treat lint as a
  compatibility check. Compile/test affected Go consumers after regeneration.
- Inspect the complete generated diff. Include outputs required by the changed
  schema and call out SDK updates separately; avoid unrelated generator-version
  churn. Report missing generation prerequisites rather than claiming outputs
  are synchronized.
