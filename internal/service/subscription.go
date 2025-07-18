package service

import (
	"context"
	"gorestsubs/internal/models"
	"gorestsubs/internal/repository"
	"time"

	"github.com/google/uuid"
)

type SubscriptionService interface {
	Create(ctx context.Context, sub *models.Subscription) (uuid.UUID, error)
	GetByID(ctx context.Context, id uuid.UUID) (*models.Subscription, error)
	GetAll(ctx context.Context) ([]*models.Subscription, error)
	Update(ctx context.Context, sub *models.Subscription) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetTotalCost(ctx context.Context, userID uuid.UUID, serviceName, startPeriod, endPeriod string) (int, error)
}

type subscriptionService struct {
	repo repository.SubscriptionRepository
}

func NewSubscriptionService(repo repository.SubscriptionRepository) SubscriptionService {
	return &subscriptionService{repo: repo}
}

func (s *subscriptionService) Create(ctx context.Context, sub *models.Subscription) (uuid.UUID, error) {
	return s.repo.Create(ctx, sub)
}

func (s *subscriptionService) GetByID(ctx context.Context, id uuid.UUID) (*models.Subscription, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *subscriptionService) GetAll(ctx context.Context) ([]*models.Subscription, error) {
	return s.repo.GetAll(ctx)
}

func (s *subscriptionService) Update(ctx context.Context, sub *models.Subscription) error {
	return s.repo.Update(ctx, sub)
}

func (s *subscriptionService) Delete(ctx context.Context, id uuid.UUID) error {
	return s.repo.Delete(ctx, id)
}

func (s *subscriptionService) GetTotalCost(ctx context.Context, userID uuid.UUID, serviceName, startPeriod, endPeriod string) (int, error) {
	subs, err := s.repo.GetTotalCost(ctx, userID, serviceName, startPeriod, endPeriod)
	if err != nil {
		return 0, err
	}

	const layout = "01-2006"
	start, err := time.Parse(layout, startPeriod)
	if err != nil {
		return 0, err
	}

	end, err := time.Parse(layout, endPeriod)
	if err != nil {
		return 0, err
	}

	totalCost := 0
	for _, sub := range subs {
		subStart, err := time.Parse(layout, sub.StartDate)
		if err != nil {
			continue
		}

		subEnd := end.AddDate(0, 1, 0)
		if sub.EndDate != nil && *sub.EndDate != "" {
			subEnd, err = time.Parse(layout, *sub.EndDate)
			if err != nil {
				continue
			}
		}

		overlapStart := max(start, subStart)
		overlapEnd := min(end, subEnd)

		if !overlapStart.After(overlapEnd) {
			months := 0
			for d := overlapStart; d.Before(overlapEnd) || d.Equal(overlapEnd); d = d.AddDate(0, 1, 0) {
				if (d.Before(subEnd) || d.Equal(subEnd)) && (d.After(subStart) || d.Equal(subStart)) {
					months++
				}
			}
			totalCost += months * sub.Price
		}
	}

	return totalCost, nil
}

func min(a, b time.Time) time.Time {
	if a.Before(b) {
		return a
	}
	return b
}

func max(a, b time.Time) time.Time {
	if a.After(b) {
		return a
	}
	return b
}
