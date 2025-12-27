package service

import (
	"context"
	"time"

	"habit-tracker/internal/models"
	"habit-tracker/internal/repository"
)

type PointsService struct {
	users  *repository.UserRepository
	points *repository.PointsRepository
}

func NewPointsService(users *repository.UserRepository, points *repository.PointsRepository) *PointsService {
	return &PointsService{users: users, points: points}
}

// AddPoints applies delta to user's points and logs the change.
func (s *PointsService) AddPoints(ctx context.Context, userID uint64, delta int64, reason string, relatedHabitID *uint64) error {
	if err := s.users.UpdatePoints(ctx, userID, delta); err != nil {
		return err
	}
	log := &models.UserPointsLog{
		UserID:         userID,
		ChangeAmount:   int(delta),
		Reason:         reason,
		RelatedHabitID: relatedHabitID,
		CreatedAt:      time.Now(),
	}
	if err := s.points.AddLog(ctx, log); err != nil {
		return err
	}
	return nil
}

func (s *PointsService) SumByRange(ctx context.Context, userID uint64, start, end time.Time) (int64, error) {
	return s.points.SumByUserAndRange(ctx, userID, start, end)
}

func (s *PointsService) GetUserPoints(ctx context.Context, userID uint64) (int64, error) {
	user, err := s.users.GetByID(ctx, userID)
	if err != nil {
		return 0, err
	}
	return user.Points, nil
}
