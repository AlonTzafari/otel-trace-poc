package main

import (
	"context"
	"fmt"

	"github.com/alontzafari/otel-trace-poc/be/worker/config"
	"github.com/alontzafari/otel-trace-poc/be/worker/db"
	"github.com/alontzafari/otel-trace-poc/be/worker/queue"
	"github.com/alontzafari/otel-trace-poc/be/worker/server"
	"github.com/alontzafari/otel-trace-poc/be/worker/telemetry"
)

func main() {
	workerConf, err := config.FromEnv()
	if err != nil {
		panic(err)
	}

	ctx := context.Background()

	tp, err := telemetry.InitTracer(ctx, workerConf.COLLECTOR, "worker")
	if err != nil {
		panic(err)
	}
	defer telemetry.Shutdown(tp)

	dbClient := db.New(workerConf.MONGO_URL)
	err = dbClient.Connect(ctx)
	if err != nil {
		panic(err)
	}
	defer dbClient.Disconnect(ctx)

	fmt.Println("db connected")

	producer := queue.NewProducer(workerConf.KAFKA_BROKER)

	srv := server.New(dbClient, producer)

	addr := fmt.Sprintf("%s:%d", workerConf.IFACE, workerConf.PORT)

	fmt.Printf("listening on %s\n", addr)
	err = srv.Start(addr)
	fmt.Printf("SERVER ERR: %v", err)

}
