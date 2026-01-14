package main

import (
	"context"
	"fmt"
	"os"

	"github.com/alontzafari/otel-trace-poc/be/worker/server"
	"github.com/alontzafari/otel-trace-poc/be/worker/telemetry"
)

const service = "worker"

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8090"
	}
	iface := os.Getenv("IFACE")
	if iface == "" {
		iface = "127.0.0.1"
	}
	collector := os.Getenv("COLLECTOR")
	if collector == "" {
		collector = "localhost:4317"
	}

	ctx := context.Background()

	tp, err := telemetry.InitTracer(ctx, collector, service)
	if err != nil {
		panic(err)
	}
	defer telemetry.Shutdown(tp)

	srv := server.New()

	addr := fmt.Sprintf("%s:%s", iface, port)

	fmt.Printf("listening on %s\n", addr)
	srv.Start(addr)

}
