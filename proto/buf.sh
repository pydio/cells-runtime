#!/usr/bin/env bash
set -euo pipefail

root="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
target="${1:?usage: proto/buf.sh <broker|config|log|options/orm|options/setter|registry|server|service>}"
bin="$root/bin"
mkdir -p "$bin"

go -C "$root" build -o "$bin/protoc-gen-go-client-stub" ./tools/protoc-gen-go-client-stub
go -C "$root" build -o "$bin/protoc-gen-go-setter" ./tools/protoc-gen-go-setter
go -C "$root" build -o "$bin/protoc-gen-go-tags" ./tools/protoc-gen-go-tags

[[ -x "$bin/protoc-gen-go" ]] || GOBIN="$bin" go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.36.11
[[ -x "$bin/protoc-gen-go-grpc" ]] || GOBIN="$bin" go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.5.1

export PATH="$bin:$PATH"
cd "$root/proto/$target"
if ! buf lint; then
	echo "buf lint reported existing schema-style violations; generation will continue" >&2
fi
buf generate --output .
if [[ -f buf.gen.tag.yaml ]]; then
	buf generate --template buf.gen.tag.yaml --output .
fi
