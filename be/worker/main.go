package main

import (
	"context"
	"fmt"
	"net"

	"github.com/alontzafari/otel-trace-poc/proto/hello"
	"github.com/alontzafari/otel-trace-poc/proto/test"
	"google.golang.org/grpc"
)

type helloServer struct {
	hello.UnimplementedHelloServer
}

// Hello implements [hello.HelloServer].
func (h *helloServer) Hello(ctx context.Context, req *hello.HelloReq) (*hello.HelloRes, error) {
	if req.Msg == "error" {
		return nil, fmt.Errorf("errorrrrr")
	}
	return &hello.HelloRes{Msg: "helllooo go"}, nil
}

var _ hello.HelloServer = (*helloServer)(nil)

type TestServer struct {
	test.UnimplementedTestServer
}

// Test implements [test.TestServer].
func (t *TestServer) Test(context.Context, *test.TestReq) (*test.TestRes, error) {
	panic("unimplemented")
}

var _ test.TestServer = (*TestServer)(nil)

const addr = "0.0.0.0:90"

func main() {
	srv := grpc.NewServer()
	hello.RegisterHelloServer(srv, &helloServer{})
	test.RegisterTestServer(srv, &TestServer{})
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		panic(err)
	}
	defer lis.Close()
	fmt.Printf("listening on %s\n", addr)
	srv.Serve(lis)
}
