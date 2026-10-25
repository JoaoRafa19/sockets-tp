package tcp

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"log"

	"sockets-tc/internal/domain"
	"sockets-tc/internal/service"
)

type HandlerFunc func(ctx context.Context, data json.RawMessage) Response

type Handler struct {
	events       *service.EventService
	reservations *service.ReservationService
	routes       map[string]HandlerFunc
}

func NewHandler(events *service.EventService, reservations *service.ReservationService) *Handler {
	h := &Handler{events: events, reservations: reservations}
	h.routes = h.newRoutes()
	return h
}

func (h *Handler) Handle(ctx context.Context, req Request) (resp Response) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("panic in %q: %v", req.Method, r)
			resp = fail("internal error")
		}
	}()

	fn, found := h.routes[req.Method]
	if !found {
		return fail("invalid method")
	}
	return fn(ctx, req.Data)
}

func (h *Handler) createEvent(ctx context.Context, data json.RawMessage) Response {
	var in createEventRequest
	if err := decodeData(data, &in); err != nil {
		return fail("invalid payload: " + err.Error())
	}

	event, err := h.events.CreateEvent(ctx, service.CreateEventInput{
		Name:    in.Name,
		Tickets: in.Tickets,
	})
	if err != nil {
		return errorResponse(err)
	}
	return ok(event)
}

func (h *Handler) reserveTicket(ctx context.Context, data json.RawMessage) Response {
	var in reserveTicketRequest
	if err := decodeData(data, &in); err != nil {
		return fail("invalid payload: " + err.Error())
	}

	res, err := h.reservations.Reserve(ctx, service.ReserveInput{
		EventID:  in.EventID,
		Email:    in.Email,
		Type:     in.Type,
		Quantity: in.Quantity,
	})
	if err != nil {
		return errorResponse(err)
	}
	return ok(res)
}

func (h *Handler) getReservation(ctx context.Context, data json.RawMessage) Response {
	var in reservationRequest
	if err := decodeData(data, &in); err != nil {
		return fail("invalid payload: " + err.Error())
	}

	res, err := h.reservations.Get(ctx, in.ReservationID, in.Email)
	if err != nil {
		return errorResponse(err)
	}
	return ok(res)
}

func (h *Handler) cancelReservation(ctx context.Context, data json.RawMessage) Response {
	var in reservationRequest
	if err := decodeData(data, &in); err != nil {
		return fail("invalid payload: " + err.Error())
	}

	res, err := h.reservations.Cancel(ctx, in.ReservationID, in.Email)
	if err != nil {
		return errorResponse(err)
	}
	return ok(res)
}

func decodeData(data json.RawMessage, v any) error {
	if len(data) == 0 {
		return errors.New("missing data")
	}
	dec := json.NewDecoder(bytes.NewReader(data))
	dec.DisallowUnknownFields()
	return dec.Decode(v)
}

func errorResponse(err error) Response {
	if domain.IsBusinessError(err) {
		return fail(err.Error())
	}
	log.Printf("internal error: %v", err)
	return fail("internal error")
}