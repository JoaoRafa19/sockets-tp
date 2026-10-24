package addapters

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net"
	"sockets-tc/internal/service"
)

type CreateEventInput struct {
	Name    string `json:"name"`
	Tickets []struct {
		Type     string `json:"type"`
		Quantity int    `json:"quantity"`
	}
}

type Handler struct {
	events *service.EventService
}

func CreateEventHandler(conn net.Conn, payl Request) {

	var data CreateEventInput

	if err := json.NewDecoder(bytes.NewReader(payl.Data)).Decode(&data); err != nil {
		conn.Write(NewResponse(Response{
			Status: 400,
			Data:   "invalid payload",
		}))
	}

	fmt.Println(data)

}
