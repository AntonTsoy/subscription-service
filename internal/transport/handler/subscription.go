package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/AntonTsoy/subscription-service/internal/models"
	"github.com/AntonTsoy/subscription-service/internal/transport/dto"
	"github.com/go-chi/chi/v5"
)

type SubscriptionService interface {
	Create(ctx context.Context, sub *models.Subscription) error
	GetByID(ctx context.Context, id int) (*models.Subscription, error)
	GetAll(ctx context.Context, limit, offset int) ([]models.Subscription, error)
	Update(ctx context.Context, sub *models.Subscription) error
	Delete(ctx context.Context, id int) error
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
		return
	}

	if err := h.service.Create(r.Context(), sub); err != nil {
		http.Error(w, "failed to create subscription", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(dto.ToSubscriptionResponse(sub))
}

func (h *SubsHandler) GetSubscription(w http.ResponseWriter, r *http.Request) {
	id, err := getIntPathParam(r, "id")
	if err != nil {
		http.Error(w, "missing or invalid subscription id path parameter value", http.StatusBadRequest)
		return
	}

	sub, err := h.service.GetByID(r.Context(), id)
	if err != nil {
		http.Error(w, "failed to get subscription", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(dto.ToSubscriptionResponse(sub))
}

func (h *SubsHandler) GetAllSubscriptions(w http.ResponseWriter, r *http.Request) {
	limit := getIntQueryParam(r, "limit", 100)
	offset := getIntQueryParam(r, "offset", 0)

	subscriptions, err := h.service.GetAll(r.Context(), limit, offset)
	if err != nil {
		http.Error(w, "failed to get all subscriptions", http.StatusInternalServerError)
		return
	}

	response := make([]dto.SubscriptionResponse, len(subscriptions))
	for i, sub := range subscriptions {
		response[i] = *dto.ToSubscriptionResponse(&sub)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func (h *SubsHandler) UpdateSubscription(w http.ResponseWriter, r *http.Request) {
	var req dto.SubscriptionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	newSubData, err := dto.ToSubscription(&req)
	if err != nil {
		http.Error(w, "invalid request body parameter", http.StatusBadRequest)
		return
	}

	newSubData.ID, err = getIntPathParam(r, "id")
	if err != nil {
		http.Error(w, "missing or invalid subscription id path parameter value", http.StatusBadRequest)
		return
	}

	if err := h.service.Update(r.Context(), newSubData); err != nil {
		http.Error(w, "failed to update subscription", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *SubsHandler) DeleteSubscription(w http.ResponseWriter, r *http.Request) {
	id, err := getIntPathParam(r, "id")
	if err != nil {
		http.Error(w, "missing or invalid subscription id path parameter value", http.StatusBadRequest)
		return
	}

	if err := h.service.Delete(r.Context(), id); err != nil {
		http.Error(w, "failed to delete subscription", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *SubsHandler) GetTotalServiceSubscriptionsCost(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	// json.NewEncoder(w).Encode(response)
}

func getIntPathParam(r *http.Request, key string) (int, error) {
	valueStr := chi.URLParam(r, key)
	if valueStr == "" {
		return 0, fmt.Errorf("missing path parameter")
	}

	value, err := strconv.Atoi(valueStr)
	if err != nil {
		return 0, err
	} else if value <= 0 {
		return 0, fmt.Errorf("incorrect path parameter value")
	}

	return value, nil
}

func getIntQueryParam(r *http.Request, key string, defaultValue int) int {
	valueStr := r.URL.Query().Get(key)
	if valueStr == "" {
		return defaultValue
	}

	value, err := strconv.Atoi(valueStr)
	if err != nil || value < 0 {
		return defaultValue
	}

	return value
}
