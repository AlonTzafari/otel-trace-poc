#!/bin/bash
protoc --go_out=. \
    --go-grpc_out=. \
    --go_opt=module=github.com/alontzafari/otel-trace-poc \
    --go-grpc_opt=module=github.com/alontzafari/otel-trace-poc \
     --proto_path=proto $1