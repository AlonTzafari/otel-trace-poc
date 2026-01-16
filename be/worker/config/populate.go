package config

import (
	"fmt"
	"reflect"
	"strconv"
)

type IGetter interface {
	Get(key string) string
}

func populateStruct[T any](conf T, getter IGetter) (T, error) {
	val := reflect.ValueOf(conf).Elem()
	t := val.Type()

	for i := range t.NumField() {
		field := t.Field(i)

		fieldEnv := getter.Get(field.Name)
		if fieldEnv == "" {
			continue
		}

		switch field.Type {
		case reflect.TypeFor[int]():
			envInt, err := strconv.Atoi(fieldEnv)
			if err != nil {
				return conf, fmt.Errorf("expected type int for field %s, found %s", field.Name, fieldEnv)
			}
			val.Field(i).SetInt(int64(envInt))
		case reflect.TypeFor[string]():
			val.Field(i).SetString(fieldEnv)
		default:
			return conf, fmt.Errorf("unknown field type %v", field.Type)
		}
	}

	return conf, nil
}
