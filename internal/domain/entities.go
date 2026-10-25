package domain

import (
	"errors"
	"net/mail"
	"strings"
	"time"
)

var (
	ErrEventNotFound               = errors.New("event not found")
	ErrInvalidEventName            = errors.New("invalid event name")
	ErrNoTickets                   = errors.New("event must have at least one ticket type")
	ErrDuplicateTicketType         = errors.New("duplicate ticket type")
	ErrInvalidTicketType           = errors.New("invalid ticket type")
	ErrInvalidQuantity             = errors.New("quantity must be greater than zero")
	ErrInvalidEmail                = errors.New("invalid email")
	ErrNotEnoughTickets            = errors.New("not enough tickets")
	ErrReservationNotFound         = errors.New("reservation not found")
	ErrReservationAlreadyCancelled = errors.New("reservation already cancelled")
)

var BusinessErrors = []error{
	ErrEventNotFound, ErrInvalidEventName, ErrNoTickets, ErrDuplicateTicketType,
	ErrInvalidTicketType, ErrInvalidQuantity, ErrInvalidEmail, ErrNotEnoughTickets,
	ErrReservationNotFound, ErrReservationAlreadyCancelled,
}

func IsBusinessError(err error) bool {
	for _, e := range BusinessErrors {
		if errors.Is(err, e) {
			return true
		}
	}
	return false
}

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

func NewEvent(name string, tickets []Ticket) (Event, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return Event{}, ErrInvalidEventName
	}
	if len(tickets) == 0 {
		return Event{}, ErrNoTickets
	}

	clean := make([]Ticket, 0, len(tickets))
	seen := make(map[string]bool, len(tickets))
	for _, t := range tickets {
		t.Type = strings.TrimSpace(t.Type)
		if t.Type == "" {
			return Event{}, ErrInvalidTicketType
		}
		if t.Quantity <= 0 {
			return Event{}, ErrInvalidQuantity
		}
		if seen[t.Type] {
			return Event{}, ErrDuplicateTicketType
		}
		seen[t.Type] = true
		clean = append(clean, t)
	}

	return Event{Name: name, Tickets: clean}, nil
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

func NewReservation(eventID, email, ticketType string, quantity int) (Reservation, error) {
	eventID = strings.TrimSpace(eventID)
	if eventID == "" {
		return Reservation{}, ErrEventNotFound
	}
	email, err := NormalizeEmail(email)
	if err != nil {
		return Reservation{}, err
	}
	ticketType = strings.TrimSpace(ticketType)
	if ticketType == "" {
		return Reservation{}, ErrInvalidTicketType
	}
	if quantity <= 0 {
		return Reservation{}, ErrInvalidQuantity
	}

	return Reservation{
		EventID:  eventID,
		Email:    email,
		Type:     ticketType,
		Quantity: quantity,
		Status:   ReservationActive,
	}, nil
}

func NormalizeEmail(email string) (string, error) {
	email = strings.ToLower(strings.TrimSpace(email))
	addr, err := mail.ParseAddress(email)
	if err != nil || addr.Address != email {
		return "", ErrInvalidEmail
	}
	return email, nil
}
