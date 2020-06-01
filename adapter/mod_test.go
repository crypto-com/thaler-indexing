package adapter_test

import (
	"fmt"

	jsoniter "github.com/json-iterator/go"
)

func JsonMustMarshal(value interface{}) string {
	b, err := jsoniter.Marshal(value)
	if err != nil {
		panic(fmt.Sprintf("error when JSON marshalling: %v", err))
	}

	return string(b)
}
