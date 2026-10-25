package ports

import (
	"context"

	"sockets-tc/internal/domain"
)

type EventRepository interface {
	CreateEvent(ctx context.Context, e domain.Event) (domain.Event, error)
}

type ReservationRepository interface {
	Reserve(ctx context.Context, r domain.Reservation) (domain.Reservation, error)
	GetReservation(ctx context.Context, id, email string) (domain.Reservation, error)
	CancelReservation(ctx context.Context, id, email string) (domain.Reservation, error)
}