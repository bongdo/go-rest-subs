package repository

import (
	"context"
	"gorestsubs/internal/models"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type SubscriptionRepository interface {
	Create(ctx context.Context, sub *models.Subscription) (uuid.UUID, error)
	GetByID(ctx context.Context, id uuid.UUID) (*models.Subscription, error)
	GetAll(ctx context.Context) ([]*models.Subscription, error)
	Update(ctx context.Context, sub *models.Subscription) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetTotalCost(ctx context.Context, userID uuid.UUID, serviceName, startPeriod, endPeriod string) ([]*models.Subscription, error)
}

type subscriptionRepository struct {
	db *pgxpool.Pool
}

func NewSubscriptionRepository(db *pgxpool.Pool) SubscriptionRepository {
	return &subscriptionRepository{db: db}
}

func (r *subscriptionRepository) Create(ctx context.Context, sub *models.Subscription) (uuid.UUID, error) {
	var id uuid.UUID
	err := r.db.QueryRow(ctx, `INSERT INTO subscriptions (service_name, price, user_id, start_date, end_date) VALUES ($1, $2, $3, $4, $5) RETURNING id`,
		sub.ServiceName, sub.Price, sub.UserID, sub.StartDate, sub.EndDate).Scan(&id)
	return id, err
}

func (r *subscriptionRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Subscription, error) {
	row := r.db.QueryRow(ctx, `SELECT id, service_name, price, user_id, start_date, end_date FROM subscriptions WHERE id = $1`, id)
	sub := &models.Subscription{}
	err := row.Scan(&sub.ID, &sub.ServiceName, &sub.Price, &sub.UserID, &sub.StartDate, &sub.EndDate)
	return sub, err
}

func (r *subscriptionRepository) GetAll(ctx context.Context) ([]*models.Subscription, error) {
	rows, err := r.db.Query(ctx, `SELECT id, service_name, price, user_id, start_date, end_date FROM subscriptions`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var subs []*models.Subscription
	for rows.Next() {
		sub := &models.Subscription{}
		if err := rows.Scan(&sub.ID, &sub.ServiceName, &sub.Price, &sub.UserID, &sub.StartDate, &sub.EndDate); err != nil {
			return nil, err
		}
		subs = append(subs, sub)
	}
	return subs, nil
}

func (r *subscriptionRepository) Update(ctx context.Context, sub *models.Subscription) error {
	_, err := r.db.Exec(ctx, `UPDATE subscriptions SET service_name = $1, price = $2, user_id = $3, start_date = $4, end_date = $5 WHERE id = $6`,
		sub.ServiceName, sub.Price, sub.UserID, sub.StartDate, sub.EndDate, sub.ID)
	return err
}

func (r *subscriptionRepository) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := r.db.Exec(ctx, `DELETE FROM subscriptions WHERE id = $1`, id)
	return err
}

func (r *subscriptionRepository) GetTotalCost(ctx context.Context, userID uuid.UUID, serviceName, startPeriod, endPeriod string) ([]*models.Subscription, error) {
	query := `SELECT id, service_name, price, user_id, start_date, end_date FROM subscriptions WHERE user_id = $1 AND ($2 = '' OR service_name = $2)`
	rows, err := r.db.Query(ctx, query, userID, serviceName)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var subs []*models.Subscription
	for rows.Next() {
		sub := &models.Subscription{}
		if err := rows.Scan(&sub.ID, &sub.ServiceName, &sub.Price, &sub.UserID, &sub.StartDate, &sub.EndDate); err != nil {
			return nil, err
		}
		subs = append(subs, sub)
	}
	return subs, nil
}
