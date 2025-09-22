package dto

import (
	"fmt"
	"time"

	"github.com/AntonTsoy/subscription-service/internal/models"
	"github.com/google/uuid"
)

const layout = "01-2006"

func ToSubscription(req *SubscriptionRequest) (*models.Subscription, error) {
	userID, err := uuid.Parse(req.UserID)
	if err != nil {
		return nil, fmt.Errorf("неверный формат user_id: %w", err)
	}

	start, err := time.Parse(layout, req.StartDate)
	if err != nil {
		return nil, fmt.Errorf("неверный формат start_date: %w", err)
	}

	var end *time.Time
	if req.EndDate != "" {
		t, err := time.Parse(layout, req.EndDate)
		if err != nil {
			return nil, fmt.Errorf("неверный формат end_date: %w", err)
		}
		end = &t
	}

	return &models.Subscription{
		ServiceName: req.ServiceName,
		Price:       req.Price,
		UserID:      userID,
		StartDate:   start,
		EndDate:     end,
	}, nil
}

func ToSubscriptionResponse(sub *models.Subscription) *SubscriptionResponse {
	resp := SubscriptionResponse{
		ID:          sub.ID,
		ServiceName: sub.ServiceName,
		Price:       sub.Price,
		UserID:      sub.UserID.String(),
		StartDate:   sub.StartDate.Format(layout),
	}
	if sub.EndDate != nil {
		resp.EndDate = sub.EndDate.Format(layout)
	}
	return &resp
}

func ToListSubscriptionsParams(req *TotalSubscriptionsCostRequest) (*models.ListSubscriptionsParams, error) {
	start, err := time.Parse(layout, req.StartDate)
	if err != nil {
		return nil, fmt.Errorf("неверный формат start_date: %w", err)
	}

	end, err := time.Parse(layout, req.EndDate)
	if err != nil {
		return nil, fmt.Errorf("неверный формат end_date: %w", err)
	}

	model := &models.ListSubscriptionsParams{
		StartDate: start,
		EndDate:   end,
	}

	if req.UserID != "" {
		id, err := uuid.Parse(req.UserID)
		if err != nil {
			return nil, fmt.Errorf("неверный формат user_id: %w", err)
		}
		model.UserID = &id
	}

	if req.ServiceName != "" {
		model.ServiceName = &req.ServiceName
	}

	return model, nil
}
