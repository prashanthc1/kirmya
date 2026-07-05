package application

import (
	"context"

	"workspace-app/internal/network/domain"
)

type Service struct {
	repo domain.Repository
}

func NewService(repo domain.Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) SendRequest(ctx context.Context, requesterID, receiverID string) (*domain.Connection, error) {
	if requesterID == receiverID {
		return nil, domain.ErrSelfConnection
	}

	status, _, err := s.repo.GetConnectionStatus(ctx, requesterID, receiverID)
	if err != nil {
		return nil, err
	}

	if status == domain.StatusPending || status == domain.StatusAccepted {
		return nil, domain.ErrDuplicateRequest
	}

	if status == domain.StatusRejected {
		if err := s.repo.Delete(ctx, requesterID, receiverID); err != nil {
			return nil, err
		}
	}

	return s.repo.Create(ctx, requesterID, receiverID)
}

func (s *Service) AcceptRequest(ctx context.Context, receiverID, connectionID string) error {
	c, err := s.repo.GetByID(ctx, connectionID)
	if err != nil {
		return err
	}

	if c.ReceiverID != receiverID {
		return domain.ErrNotFound
	}

	if c.Status != domain.StatusPending {
		return domain.ErrInvalidTransition
	}

	return s.repo.UpdateStatus(ctx, connectionID, domain.StatusAccepted)
}

func (s *Service) RejectRequest(ctx context.Context, receiverID, connectionID string) error {
	c, err := s.repo.GetByID(ctx, connectionID)
	if err != nil {
		return err
	}

	if c.ReceiverID != receiverID {
		return domain.ErrNotFound
	}

	if c.Status != domain.StatusPending {
		return domain.ErrInvalidTransition
	}

	return s.repo.UpdateStatus(ctx, connectionID, domain.StatusRejected)
}

func (s *Service) GetConnections(ctx context.Context, userID string) ([]domain.Connection, error) {
	return s.repo.GetConnections(ctx, userID)
}

func (s *Service) GetIncomingRequests(ctx context.Context, userID string) ([]domain.Connection, error) {
	return s.repo.GetIncomingRequests(ctx, userID)
}

func (s *Service) GetConnectionStatus(ctx context.Context, userA, userB string) (domain.ConnectionStatus, string, error) {
	return s.repo.GetConnectionStatus(ctx, userA, userB)
}
