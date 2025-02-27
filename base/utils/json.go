package utils

import (
	"fmt"
	jsoniter "github.com/json-iterator/go"
)

func MustMarshal(v interface{}) []byte {
	bytes, err := jsoniter.Marshal(v)
	if err != nil {
		panic(fmt.Sprintf("utils.MustMarshal failed to Marshal: %v", v))
	}
	return bytes
}

func MustUnmarshal(bytes []byte, v interface{}) {
	err := jsoniter.Unmarshal(bytes, v)
	if err != nil {
		panic(fmt.Sprintf("utils.MustUnmarshal failed to Unmarshal: %s", string(bytes)))
	}
}
