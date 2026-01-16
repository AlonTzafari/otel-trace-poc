package server

import (
	"context"
	"fmt"
	"math/rand/v2"

	"github.com/alontzafari/otel-trace-poc/be/worker/db"
	"github.com/alontzafari/otel-trace-poc/be/worker/telemetry"
	"github.com/alontzafari/otel-trace-poc/proto/hello"
	"go.opentelemetry.io/otel/attribute"
)

type helloServer struct {
	hello.UnimplementedHelloServer
	dbClient *db.DB
}

// Hello implements [hello.HelloServer].
func (h *helloServer) Hello(ctx context.Context, req *hello.HelloReq) (*hello.HelloRes, error) {
	if req.Msg == "error" {
		return nil, fmt.Errorf("errorrrrr")
	}

	num := diceRoll(ctx)

	err := h.dbClient.SaveDiceRoll(ctx, db.DiceRoll{Value: num})
	if err != nil {
		return nil, err
	}

	return &hello.HelloRes{Msg: fmt.Sprintf("hello %s %d", req.Msg, num)}, nil
}

var _ hello.HelloServer = (*helloServer)(nil)

func diceRoll(ctx context.Context) int {
	_, span := telemetry.StartSpan(ctx, "diceRoll")
	defer span.End()

	val := rand.IntN(100)
	span.SetAttributes(attribute.Int("diceResult", val))

	return val
}
