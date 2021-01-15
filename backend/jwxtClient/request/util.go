package request

import (
	"encoding/json"
)

func PanicIf(err error) {
	if err != nil {
		panic(err)
	}
}

func JsonToMap(data []byte) map[string]interface{} {
	var d map[string]interface{}
	PanicIf(json.Unmarshal(data, &d))
	return d
}

func JsonToStruct(data []byte, v interface{}) {
	PanicIf(json.Unmarshal(data, v))
}

func JsonConvert(data []byte, v interface{}) {
	PanicIf(json.Unmarshal(data, v))
}
