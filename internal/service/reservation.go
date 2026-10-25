package service

import (
	"context"
	"strings"

	"sockets-tc/internal/domain"
	"sockets-tc/internal/ports"
)

type ReserveInput struct {
	EventID  string
	Email    string
	Type     string
	Quantity int
}

type ReservationService struct {
	repo ports.ReservationRepository
}

func NewReservationService(repo ports.ReservationRepository) *ReservationService {
	return &ReservationService{repo: repo}
}

func (s *ReservationService) Reserve(ctx context.Context, in ReserveInput) (domain.Reservation, error) {
	res, err := domain.NewReservation(in.EventID, in.Email, in.Type, in.Quantity)
	if err != nil {
		return domain.Reservation{}, err
	}
	return s.repo.Reserve(ctx, res)
}

func (s *ReservationService) Get(ctx context.Context, id, email string) (domain.Reservation, error) {
	id, email, err := normalizeLookup(id, email)
	if err != nil {
		return domain.Reservation{}, err
	}
	return s.repo.GetReservation(ctx, id, email)
}

func (s *ReservationService) Cancel(ctx context.Context, id, email string) (domain.Reservation, error) {
	id, email, err := normalizeLookup(id, email)
	if err != nil {
		return domain.Reservation{}, err
	}
	return s.repo.CancelReservation(ctx, id, email)
}

func normalizeLookup(id, email string) (string, string, error) {
	id = strings.TrimSpace(id)
	if id == "" {
		return "", "", domain.ErrReservationNotFound
	}
	email, err := domain.NormalizeEmail(email)
	if err != nil {
		return "", "", err
	}
	return id, email, nil
}