package service

import (
	"context"
	"time"

	"github.com/AntonTsoy/subscription-service/internal/models"
	"github.com/google/uuid"
)

type SubscriptionRepository interface {
	Create(ctx context.Context, sub *models.Subscription) error
	GetByID(ctx context.Context, id int) (*models.Subscription, error)
	GetAll(ctx context.Context, limit, offset int) ([]models.Subscription, error)
	Update(ctx context.Context, sub *models.Subscription) error
	Delete(ctx context.Context, id int) error
	ListByUserAndService(ctx context.Context, userID uuid.UUID, serviceName string, start, end time.Time) ([]models.Subscription, error)
}

type SubsService struct {
	repo SubscriptionRepository
}

func NewSubsService(repo SubscriptionRepository) *SubsService {
	return &SubsService{repo: repo}
}

func (s *SubsService) Create(ctx context.Context, sub *models.Subscription) error {
	return s.repo.Create(ctx, sub)
}

func (s *SubsService) GetByID(ctx context.Context, id int) (*models.Subscription, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *SubsService) GetAll(ctx context.Context, limit, offset int) ([]models.Subscription, error) {
	// TODO: улучшить логику обработки случая, когда limit и offset не указан
	if limit == 0 {
		limit = 200
	}
	return s.repo.GetAll(ctx, limit, offset)
}

func (s *SubsService) Update(ctx context.Context, sub *models.Subscription) error {
	return s.repo.Update(ctx, sub)
}

func (s *SubsService) Delete(ctx context.Context, id int) error {
	return s.repo.Delete(ctx, id)
}

func (s *SubsService) EvaluateCostOfServiceIntervalSubscriptions() {
	// TODO: реализовать метод
}
