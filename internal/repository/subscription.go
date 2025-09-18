package repository

import (
	"context"
	"fmt"

	"github.com/AntonTsoy/subscription-service/internal/models"
	"github.com/jmoiron/sqlx"
)

type SubsRepo struct {
	db *sqlx.DB
}

func NewSubsRepo(db *sqlx.DB) *SubsRepo {
	return &SubsRepo{db: db}
}

func (r *SubsRepo) Create(ctx context.Context, sub *models.Subscription) error {
	query := `
        INSERT INTO subscriptions (service_name, price, user_id, start_date, end_date)
        VALUES ($1, $2, $3, $4, $5)
        RETURNING id
    `

	var endDate any = nil
	if sub.EndDate != nil {
		endDate = *sub.EndDate
	}

	err := r.db.QueryRowContext(
		ctx,
		query,
		sub.ServiceName,
		sub.Price,
		sub.UserID,
		sub.StartDate,
		endDate,
	).Scan(&sub.ID)

	if err != nil {
		return fmt.Errorf("failed to create subscription: %w", err)
	}

	return nil
}
