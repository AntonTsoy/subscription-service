package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/AntonTsoy/subscription-service/internal/models"
	"github.com/google/uuid"
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
        VALUES (:service_name, :price, :user_id, :start_date, :end_date)
        RETURNING id
    `

	stmt, err := r.db.PrepareNamed(query)
	if err != nil {
		return fmt.Errorf("подготовка записи подписки: %w", err)
	}

	return stmt.QueryRowxContext(ctx, sub).Scan(&sub.ID)
}

func (r *SubsRepo) GetByID(ctx context.Context, id int) (*models.Subscription, error) {
	query := `SELECT * FROM subscriptions WHERE id=$1`

	var sub models.Subscription
	if err := r.db.GetContext(ctx, &sub, query, id); err != nil {
		return nil, fmt.Errorf("ошибка получения подписки: %w", err)
	}

	return &sub, nil
}

func (r *SubsRepo) GetAll(ctx context.Context, limit, offset int) ([]models.Subscription, error) {
	query := `
        SELECT * FROM subscriptions
        ORDER BY id
        LIMIT $1 OFFSET $2
    `

	var subs []models.Subscription
	if err := r.db.SelectContext(ctx, &subs, query, limit, offset); err != nil {
		return nil, fmt.Errorf("ошибка получения подписок: %w", err)
	}

	return subs, nil
}

func (r *SubsRepo) Update(ctx context.Context, sub *models.Subscription) error {
	query := `
        UPDATE subscriptions
        SET service_name = :service_name,
            price = :price,
            user_id = :user_id,
            start_date = :start_date,
            end_date = :end_date
        WHERE id = :id
    `

	res, err := r.db.NamedExecContext(ctx, query, sub)
	if err != nil {
		return fmt.Errorf("ошибка обновления записи: %w", err)
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("не удалось обновить запись: %w", err)
	}
	if rows == 0 {
		return fmt.Errorf("не найдена подписка c id=%d", sub.ID)
	}
	return nil
}

func (r *SubsRepo) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM subscriptions WHERE id=$1`

	res, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("ошибка удаления записи: %w", err)
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("не удалось удалить запись: %w", err)
	}
	if rows == 0 {
		return fmt.Errorf("не найдена подписка c id=%d", id)
	}
	return nil
}

func (r *SubsRepo) ListByUserAndService(
	ctx context.Context,
	userID uuid.UUID,
	serviceName string,
	start, end time.Time,
) ([]models.Subscription, error) {
	query := `
		SELECT * FROM subscriptions
		WHERE user_id = $1
			AND service_name = $2
			AND start_date <= $4
			AND (end_date IS NULL OR end_date >= $3);
	`

	var subs []models.Subscription
	if err := r.db.SelectContext(ctx, &subs, query, userID, serviceName, start, end); err != nil {
		return nil, fmt.Errorf("ошибка получения подписок: %w", err)
	}

	return subs, nil
}
