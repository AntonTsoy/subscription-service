package http

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/AntonTsoy/subscription-service/internal/models"
	"github.com/AntonTsoy/subscription-service/internal/transport/dto"
)

type SubscriptionService interface {
	Create(ctx context.Context, sub *models.Subscription) error
}

type SubsHandler struct {
	service SubscriptionService
}

func NewSubsHandler(service SubscriptionService) *SubsHandler {
	return &SubsHandler{service: service}
}

func (h *SubsHandler) CreateSubscription(w http.ResponseWriter, r *http.Request) {
	var req dto.SubscriptionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	sub, err := dto.ToSubscription(&req)
	if err != nil {
		http.Error(w, "invalid request body parameter", http.StatusBadRequest)
	}

	if err := h.service.Create(r.Context(), sub); err != nil {
		http.Error(w, "failed to create subscription", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(dto.ToSubscriptionResponse(sub))
}
