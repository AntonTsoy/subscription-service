package models

import (
	"time"

	"github.com/google/uuid"
)

type Subscription struct {
	ID          int        `db:"id"`
	ServiceName string     `db:"service_name"`
	Price       int        `db:"price"`
	UserID      uuid.UUID  `db:"user_id"`
	StartDate   time.Time  `db:"start_date"`
	EndDate     *time.Time `db:"end_date"`
}

type ListSubscriptionsParams struct {
	StartDate   time.Time  `db:"start_date"`
	EndDate     time.Time  `db:"end_date"`
	UserID      *uuid.UUID `db:"user_id"`
	ServiceName *string    `db:"service_name"`
}
