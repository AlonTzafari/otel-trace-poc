package server

import (
	"net"

	"github.com/alontzafari/otel-trace-poc/be/worker/db"
	"github.com/alontzafari/otel-trace-poc/proto/hello"
	"github.com/alontzafari/otel-trace-poc/proto/test"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
)

type Srv struct {
	grpcSrv  *grpc.Server
	dbClient *db.DB
}

func New(dbClient *db.DB) *Srv {
	srv := grpc.NewServer(
		grpc.StatsHandler(otelgrpc.NewServerHandler()),
	)
	hello.RegisterHelloServer(srv, &helloServer{dbClient: dbClient})
	test.RegisterTestServer(srv, &TestServer{})
	return &Srv{
		grpcSrv:  srv,
		dbClient: dbClient,
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
