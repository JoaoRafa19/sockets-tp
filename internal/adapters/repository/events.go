package postgres

import (
	"context"
	"fmt"

	"sockets-tc/internal/adapters/repository/postgres/db"
	"sockets-tc/internal/domain"
)

func (r *Repository) CreateEvent(ctx context.Context, e domain.Event) (domain.Event, error) {
	var out domain.Event

	err := r.withTx(ctx, func(q *db.Queries) error {
		row, err := q.CreateEvent(ctx, e.Name)
		if err != nil {
			return fmt.Errorf("create event: %w", err)
		}

		for _, t := range e.Tickets {
			err := q.CreateEventTicket(ctx, db.CreateEventTicketParams{
				EventID:  row.ID,
				Type:     t.Type,
				Quantity: int32(t.Quantity),
			})
			if err != nil {
				return fmt.Errorf("create ticket %q: %w", t.Type, err)
			}
		}

		out = domain.Event{
			ID:        row.ID.String(),
			Name:      row.Name,
			Tickets:   e.Tickets,
			CreatedAt: row.CreatedAt.Time,
		}
		return nil
	})

	return out, err
}
