package tcp

import (
	"encoding/json"

	"sockets-tc/internal/domain"
)

// Métodos aceitos pelo servidor.
const (
	MethodCreateEvent       = "create-event"
	MethodReserveTicket     = "reserve-ticket"
	MethodGetReservation    = "get-reservation"
	MethodCancelReservation = "cancel-reservation"
)

type Request struct {
	Method string          `json:"method"`
	Data   json.RawMessage `json:"data"`
}

type Response struct {
	Status string `json:"status"`
	Data   any    `json:"data,omitempty"`
	Error  string `json:"error,omitempty"`
}

func ok(data any) Response {
	return Response{Status: "ok", Data: data}
}

func fail(msg string) Response {
	return Response{Status: "error", Error: msg}
}

type createEventRequest struct {
	Name    string          `json:"name"`
	Tickets []domain.Ticket `json:"tickets"`
}

type reserveTicketRequest struct {
	EventID  string `json:"event_id"`
	Email    string `json:"email"`
	Type     string `json:"type"`
	Quantity int    `json:"quantity"`
}

type reservationRequest struct {
	ReservationID string `json:"reservation_id"`
	Email         string `json:"email"`
}
