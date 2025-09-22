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
	EvaluateTotalServiceSubscriptionsCost(ctx context.Context, subParams *models.ListSubscriptionsParams) (int, error)
}

type SubsHandler struct {
	service SubscriptionService
}

func NewSubsHandler(service SubscriptionService) *SubsHandler {
	return &SubsHandler{service: service}
}

// CreateSubscription godoc
// @Summary      Создать подписку
// @Description  Создаёт новую подписку пользователю
// @Tags         subscriptions
// @Accept       json
// @Produce      json
// @Param        request body dto.SubscriptionRequest true "Subscription data"
// @Success      201 {object} dto.SubscriptionResponse
// @Failure      400 {string} string "invalid request"
// @Failure      500 {string} string "failed to create subscription"
// @Router       /subscriptions [post]
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

// GetSubscription godoc
// @Summary      Получить подписку
// @Description  Возвращает подписку по её ID
// @Tags         subscriptions
// @Accept       json
// @Produce      json
// @Param        id path int true "Subscription ID"
// @Success      200 {object} dto.SubscriptionResponse
// @Failure      400 {string} string "invalid id"
// @Failure      500 {string} string "failed to get subscription"
// @Router       /subscriptions/{id} [get]
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

// GetAllSubscriptions godoc
// @Summary      Список подписок
// @Description  Возвращает все подписки с пагинацией
// @Tags         subscriptions
// @Accept       json
// @Produce      json
// @Param        limit query int false "Max items (default 100)"
// @Param        offset query int false "Offset (default 0)"
// @Success      200 {array} dto.SubscriptionResponse
// @Failure      500 {string} string "failed to get subscriptions"
// @Router       /subscriptions [get]
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

// UpdateSubscription godoc
// @Summary      Обновить подписку
// @Description  Обновляет данные подписки по ID
// @Tags         subscriptions
// @Accept       json
// @Produce      json
// @Param        id path int true "Subscription ID"
// @Param        request body dto.SubscriptionRequest true "Updated subscription"
// @Success      204 {string} string "no content"
// @Failure      400 {string} string "invalid request"
// @Failure      500 {string} string "failed to update subscription"
// @Router       /subscriptions/{id} [put]
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

// DeleteSubscription godoc
// @Summary      Удалить подписку
// @Description  Удаляет подписку по ID
// @Tags         subscriptions
// @Accept       json
// @Produce      json
// @Param        id path int true "Subscription ID"
// @Success      204 {string} string "no content"
// @Failure      400 {string} string "invalid id"
// @Failure      500 {string} string "failed to delete subscription"
// @Router       /subscriptions/{id} [delete]
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

// TotalServiceSubscriptionsCost godoc
// @Summary      Общая стоимость подписок
// @Description  Считает суммарную стоимость подписок пользователя на сервис за период [start; end]
// @Tags         subscriptions
// @Accept       json
// @Produce      json
// @Param        start path string true "Start date (MM-YYYY)"
// @Param        end path string true "End date (MM-YYYY)"
// @Param        user_id query string false "User UUID (optional)"
// @Param        service_name query string false "Service name (optional)"
// @Success      200 {object} map[string]int "total cost"
// @Failure      400 {string} string "invalid parameters"
// @Failure      500 {string} string "failed to calculate cost"
// @Router       /subscriptions/{start}/{end}/total-cost [get]
func (h *SubsHandler) TotalServiceSubscriptionsCost(w http.ResponseWriter, r *http.Request) {
	var req dto.TotalSubscriptionsCostRequest
	req.StartDate = chi.URLParam(r, "start")
	req.EndDate = chi.URLParam(r, "end")
	if req.StartDate == "" || req.EndDate == "" {
		http.Error(w, "invalid subscription perion in path parameter", http.StatusBadRequest)
		return
	}

	req.UserID = r.URL.Query().Get("user_id")
	req.ServiceName = r.URL.Query().Get("service_name")

	subParams, err := dto.ToListSubscriptionsParams(&req)
	if err != nil {
		http.Error(w, "invalid request body parameter", http.StatusBadRequest)
		return
	}

	totalCost, err := h.service.EvaluateTotalServiceSubscriptionsCost(r.Context(), subParams)
	if err != nil {
		http.Error(w, "failed to get subscriptions cost for period", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]int{"totalCost": totalCost})
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
