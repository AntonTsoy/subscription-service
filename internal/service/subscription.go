package service

import (
	"context"

	"github.com/AntonTsoy/subscription-service/internal/models"
)

type SubscriptionRepository interface {
	Create(ctx context.Context, sub *models.Subscription) error
	GetByID(ctx context.Context, id int) (*models.Subscription, error)
	GetAll(ctx context.Context, limit, offset int) ([]models.Subscription, error)
	Update(ctx context.Context, sub *models.Subscription) error
	Delete(ctx context.Context, id int) error
	ListByUserAndService(ctx context.Context, params *models.ListSubscriptionsParams) ([]models.Subscription, error)
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
	return s.repo.GetAll(ctx, limit, offset)
}

func (s *SubsService) Update(ctx context.Context, sub *models.Subscription) error {
	return s.repo.Update(ctx, sub)
}

func (s *SubsService) Delete(ctx context.Context, id int) error {
	return s.repo.Delete(ctx, id)
}

func (s *SubsService) EvaluateTotalServiceSubscriptionsCost(ctx context.Context, subParams *models.ListSubscriptionsParams) (int, error) {
	subs, err := s.repo.ListByUserAndService(ctx, subParams)
	if err != nil {
		return 0, err
	}

	totalCost := 0
	for _, sub := range subs {
		if sub.StartDate.Before(subParams.StartDate) {
			sub.StartDate = subParams.StartDate
		}
		if sub.EndDate == nil || sub.EndDate.After(subParams.EndDate) {
			sub.EndDate = &subParams.EndDate
		}

		totalCost += (1 + int(sub.EndDate.Month()) - int(sub.StartDate.Month())) * sub.Price
	}
	return totalCost, nil
}
