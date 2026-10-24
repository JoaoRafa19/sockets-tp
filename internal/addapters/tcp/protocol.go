package addapters

import (
	"bytes"
	"encoding/json"
)

type Method string

const (
	CreateEvent Method = "create-event"
)

type Request struct {
	Method Method          `json:"method"`
	Data   json.RawMessage `data:"data"`
}

type Response struct {
	Status int    `json:"status"`
	Data   string `json:"data"`
}

func NewResponse(resp Response) []byte {
	res := bytes.NewBuffer([]byte(""))
	json.NewEncoder(res).Encode(resp)
	return res.Bytes()
}
