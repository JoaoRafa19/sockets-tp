package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"sockets-tc/internal/adapters/repository/postgres/db"
	"sockets-tc/internal/domain"
)

func (r *Repository) Reserve(ctx context.Context, res domain.Reservation) (domain.Reservation, error) {
	eventID, err := uuid.Parse(res.EventID)
	if err != nil {
		return domain.Reservation{}, domain.ErrEventNotFound
	}

	var out domain.Reservation

	err = r.withTx(ctx, func(q *db.Queries) error {
		n, err := q.DecrementTickets(ctx, db.DecrementTicketsParams{
			Amount:  int32(res.Quantity),
			EventID: eventID,
			Type:    res.Type,
		})
		if err != nil {
			return fmt.Errorf("decrement tickets: %w", err)
		}
		if n == 0 {
			return reserveFailure(ctx, q, eventID, res.Type)
		}

		row, err := q.CreateReservation(ctx, db.CreateReservationParams{
			EventID:  eventID,
			Type:     res.Type,
			Email:    res.Email,
			Quantity: int32(res.Quantity),
		})
		if err != nil {
			return fmt.Errorf("create reservation: %w", err)
		}

		out = toDomainReservation(row)
		return nil
	})

	return out, err
}

func reserveFailure(ctx context.Context, q *db.Queries, eventID uuid.UUID, ticketType string) error {
	_, err := q.GetEventTicket(ctx, db.GetEventTicketParams{EventID: eventID, Type: ticketType})
	if err == nil {
		return domain.ErrNotEnoughTickets
	}
	if !errors.Is(err, pgx.ErrNoRows) {
		return fmt.Errorf("get event ticket: %w", err)
	}

	exists, err := q.EventExists(ctx, eventID)
	if err != nil {
		return fmt.Errorf("event exists: %w", err)
	}
	if !exists {
		return domain.ErrEventNotFound
	}
	return domain.ErrInvalidTicketType
}

func (r *Repository) GetReservation(ctx context.Context, id, email string) (domain.Reservation, error) {
	rid, err := uuid.Parse(id)
	if err != nil {
		return domain.Reservation{}, domain.ErrReservationNotFound
	}

	row, err := r.queries.GetReservation(ctx, db.GetReservationParams{ID: rid, Email: email})
	if errors.Is(err, pgx.ErrNoRows) {
		return domain.Reservation{}, domain.ErrReservationNotFound
	}
	if err != nil {
		return domain.Reservation{}, fmt.Errorf("get reservation: %w", err)
	}

	return toDomainReservation(row), nil
}

func (r *Repository) CancelReservation(ctx context.Context, id, email string) (domain.Reservation, error) {
	rid, err := uuid.Parse(id)
	if err != nil {
		return domain.Reservation{}, domain.ErrReservationNotFound
	}

	var out domain.Reservation

	err = r.withTx(ctx, func(q *db.Queries) error {
		// Só cancela se estiver ativa; a linha fica travada até o commit.
		row, err := q.CancelReservation(ctx, db.CancelReservationParams{ID: rid, Email: email})
		if errors.Is(err, pgx.ErrNoRows) {
			return cancelFailure(ctx, q, rid, email)
		}
		if err != nil {
			return fmt.Errorf("cancel reservation: %w", err)
		}

		// Devolve os ingressos ao estoque.
		err = q.IncrementTickets(ctx, db.IncrementTicketsParams{
			Amount:  row.Quantity,
			EventID: row.EventID,
			Type:    row.Type,
		})
		if err != nil {
			return fmt.Errorf("increment tickets: %w", err)
		}

		out = toDomainReservation(row)
		return nil
	})

	return out, err
}

func cancelFailure(ctx context.Context, q *db.Queries, id uuid.UUID, email string) error {
	_, err := q.GetReservation(ctx, db.GetReservationParams{ID: id, Email: email})
	if errors.Is(err, pgx.ErrNoRows) {
		return domain.ErrReservationNotFound
	}
	if err != nil {
		return fmt.Errorf("get reservation: %w", err)
	}
	return domain.ErrReservationAlreadyCancelled
}

func toDomainReservation(row db.Reservation) domain.Reservation {
	return domain.Reservation{
		ID:        row.ID.String(),
		EventID:   row.EventID.String(),
		Type:      row.Type,
		Email:     row.Email,
		Quantity:  int(row.Quantity),
		Status:    row.Status,
		CreatedAt: row.CreatedAt.Time,
	}
}
