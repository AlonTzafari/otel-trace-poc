package queue

import (
	"context"

	"github.com/alontzafari/otel-trace-poc/be/worker/telemetry"
	"github.com/alontzafari/otel-trace-poc/proto/hello"
	kafka "github.com/segmentio/kafka-go"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"google.golang.org/protobuf/proto"
)

type Producer struct {
	w *kafka.Writer
}

type KafkaCarrier struct {
	headers []kafka.Header
}

// Get implements [propagation.TextMapCarrier].
func (k *KafkaCarrier) Get(key string) string {
	for _, h := range k.headers {
		if h.Key == key {
			return string(h.Value)
		}
	}
	return ""
}

// Keys implements [propagation.TextMapCarrier].
func (k *KafkaCarrier) Keys() []string {
	keys := make([]string, 0, len(k.headers))
	for _, h := range k.headers {
		keys = append(keys, h.Key)
	}
	return keys
}

// Set implements [propagation.TextMapCarrier].
func (k *KafkaCarrier) Set(key string, value string) {
	k.headers = append(k.headers, kafka.Header{Key: key, Value: []byte(value)})
}

var _ propagation.TextMapCarrier = (*KafkaCarrier)(nil)

func (p *Producer) SendDiceRoll(ctx context.Context, event *hello.DiceRollEvent) error {
	ctx, span := telemetry.StartSpan(ctx, "SendDiceRoll")
	defer span.End()

	tmp := otel.GetTextMapPropagator()
	kc := &KafkaCarrier{headers: make([]kafka.Header, 0)}
	tmp.Inject(ctx, kc)

	message, err := proto.Marshal(event)
	if err != nil {
		return err
	}

	return p.w.WriteMessages(ctx, kafka.Message{
		Value:   message,
		Headers: kc.headers,
	})
}

func NewProducer(broker string) *Producer {
	w := kafka.NewWriter(kafka.WriterConfig{
		Brokers: []string{broker},
		Topic:   "diceroll",
	})

	return &Producer{
		w: w,
	}
}
