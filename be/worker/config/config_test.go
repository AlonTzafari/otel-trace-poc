package config

import (
	"testing"
)

type TestGetter struct {
	m map[string]string
}

var _ IGetter = (*TestGetter)(nil)

func (t *TestGetter) Get(key string) string {
	return t.m[key]
}

func TestGetConf(t *testing.T) {
	tg := &TestGetter{
		m: map[string]string{
			"PORT":  "8090",
			"IFACE": "127.0.0.1",
		},
	}

	conf, err := populateStruct(&Conf{COLLECTOR: "localhost:4317"}, tg)
	if err != nil {
		t.Errorf("ERR: %v", err)
	}
	t.Logf("conf: %+v", conf)
}
