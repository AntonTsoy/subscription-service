package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

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

	err := r.db.QueryRowContext(ctx, query, sub.ServiceName, sub.Price, sub.UserID, sub.StartDate, sub.EndDate).Scan(&sub.ID)
	if err != nil {
		return fmt.Errorf("не удалось записать данные подписки: %w", err)
	}
	return nil
}

func (r *SubsRepo) GetByID(ctx context.Context, id int) (*models.Subscription, error) {
	query := `SELECT * FROM subscriptions WHERE id=$1`

	var sub models.Subscription
	if err := r.db.GetContext(ctx, &sub, query, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("%w: не удалось получить подписку id %d", models.ErrSubscriptionNotFound, id)
		}
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
		return fmt.Errorf("%w: обновление данных подписки id %d", models.ErrSubscriptionNotFound, sub.ID)
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
		return fmt.Errorf("%w: удаление подписки id %d", models.ErrSubscriptionNotFound, id)
	}
	return nil
}

func (r *SubsRepo) ListByUserAndService(ctx context.Context, params *models.ListSubscriptionsParams) ([]models.Subscription, error) {
	query := `
		SELECT * FROM subscriptions
		WHERE start_date <= $2
			AND (end_date IS NULL OR end_date >= $1)
	`

	args := []any{params.StartDate, params.EndDate}
	argIndex := 3

	var conditions []string
	if params.UserID != nil {
		conditions = append(conditions, fmt.Sprintf("user_id = $%d", argIndex))
		args = append(args, *params.UserID)
		argIndex++
	}
	if params.ServiceName != nil {
		conditions = append(conditions, fmt.Sprintf("service_name = $%d", argIndex))
		args = append(args, *params.ServiceName)
		argIndex++
	}

	if len(conditions) > 0 {
		query += " AND " + strings.Join(conditions, " AND ")
	}
	query += ";"

	var subs []models.Subscription
	err := r.db.SelectContext(ctx, &subs, query, args...)
	if err != nil {
		return nil, fmt.Errorf("ошибка получения подписок: %w", err)
	}
	return subs, nil
}
