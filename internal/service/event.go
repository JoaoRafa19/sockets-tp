package service

import (
	"context"

	"sockets-tc/internal/domain"
	"sockets-tc/internal/ports"
)

type CreateEventInput struct {
	Name    string
	Tickets []domain.Ticket
}

type EventService struct {
	repo ports.EventRepository
}

func NewEventService(repo ports.EventRepository) *EventService {
	return &EventService{repo: repo}
}

func (s *EventService) CreateEvent(ctx context.Context, in CreateEventInput) (domain.Event, error) {
	event, err := domain.NewEvent(in.Name, in.Tickets)
	if err != nil {
		return domain.Event{}, err
	}
	return s.repo.CreateEvent(ctx, event)
}