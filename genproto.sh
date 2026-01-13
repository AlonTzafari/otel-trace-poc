#!/bin/bash
protoc --go_out=. \
    --go-grpc_out=. \
    --go_opt=module=github.com/alontzafari/otel-trace-poc \
    --go-grpc_opt=module=github.com/alontzafari/otel-trace-poc \
    --proto_path=proto \
    $1

protoc --plugin=protoc-gen-ts_proto=".\\node_modules\\.bin\\protoc-gen-ts_proto.cmd" \
    --ts_proto_out=./fe/proto/ \
    --ts_proto_opt=outputServices=grpc-js \
    --ts_proto_opt=env=node \
    --ts_proto_opt=esModuleInterop=true \
    $1