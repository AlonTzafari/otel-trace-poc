package config

import (
	"os"
)

type Conf struct {
	PORT         int
	IFACE        string
	COLLECTOR    string
	MONGO_URL    string
	KAFKA_BROKER string
}

func getDefaultConf() *Conf {
	return &Conf{
		PORT:         8090,
		IFACE:        "127.0.0.1",
		COLLECTOR:    "localhost:4317",
		MONGO_URL:    "mongodb://localhost:27017",
		KAFKA_BROKER: "kafka:9092",
	}
}

type EnvGetter struct{}

var _ IGetter = (*EnvGetter)(nil)

func (e *EnvGetter) Get(key string) string {
	return os.Getenv(key)
}

func FromEnv() (*Conf, error) {
	conf := getDefaultConf()

	return populateStruct(
		conf,
		&EnvGetter{},
	)
}
