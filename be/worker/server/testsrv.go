package server

import (
	"context"

	"github.com/alontzafari/otel-trace-poc/proto/test"
)

type TestServer struct {
	test.UnimplementedTestServer
}

// Test implements [test.TestServer].
func (t *TestServer) Test(context.Context, *test.TestReq) (*test.TestRes, error) {
	panic("unimplemented")
}

var _ test.TestServer = (*TestServer)(nil)
