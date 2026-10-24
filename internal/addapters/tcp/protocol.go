package addapters

import "encoding/json"

type Method string

const (
	CreateEvent Method = "create-event"
)

type HandlerFunc func(Payload)

var MethodMap = map[Method]HandlerFunc{}

type Payload struct {
	Method Method          `json:"method"`
	Data   json.RawMessage `data:"data"`
}
