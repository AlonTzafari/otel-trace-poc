package server

import (
	"net"

	"github.com/alontzafari/otel-trace-poc/proto/hello"
	"github.com/alontzafari/otel-trace-poc/proto/test"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
)

type Srv struct {
	grpcSrv *grpc.Server
}

func New() *Srv {
	srv := grpc.NewServer(
		grpc.StatsHandler(otelgrpc.NewServerHandler()),
	)
	hello.RegisterHelloServer(srv, &helloServer{})
	test.RegisterTestServer(srv, &TestServer{})
	return &Srv{
		grpcSrv: srv,
	}
}

func (s *Srv) Start(addr string) error {
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	defer lis.Close()

	return s.grpcSrv.Serve(lis)
}
