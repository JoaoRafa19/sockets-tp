package domain

import (
	"errors"
	"time"
)

var (
	ErrEventNotFound               = errors.New("event not found")
	ErrInvalidTicketType           = errors.New("invalid ticket type")
	ErrNotEnoughTickets            = errors.New("not enough tickets")
	ErrReservationNotFound         = errors.New("reservation not found")
	ErrReservationAlreadyCancelled = errors.New("reservation already cancelled")
)

type Ticket struct {
	Type     string `json:"type"`
	Quantity int    `json:"quantity"`
}

type Event struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Tickets   []Ticket  `json:"tickets"`
	CreatedAt time.Time `json:"created_at"`
}

const (
	ReservationActive    = "active"
	ReservationCancelled = "cancelled"
)

type Reservation struct {
	ID        string    `json:"id"`
	EventID   string    `json:"event_id"`
	Type      string    `json:"type"`
	Email     string    `json:"email"`
	Quantity  int       `json:"quantity"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}