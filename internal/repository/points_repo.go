package repository

import (
	"context"
	"time"

	"gorm.io/gorm"

	"habit-tracker/internal/models"
)

type PointsRepository struct {
	db *gorm.DB
}

func NewPointsRepository(db *gorm.DB) *PointsRepository {
	return &PointsRepository{db: db}
}

// AddLog inserts a points change record.
func (r *PointsRepository) AddLog(ctx context.Context, log *models.UserPointsLog) error {
	return r.db.WithContext(ctx).Create(log).Error
}

// SumByUserAndRange aggregates points change amount in a time window.
func (r *PointsRepository) SumByUserAndRange(ctx context.Context, userID uint64, start, end time.Time) (int64, error) {
	var total int64
	err := r.db.WithContext(ctx).
		Model(&models.UserPointsLog{}).
		Select("COALESCE(SUM(change_amount),0)").
		Where("user_id = ? AND created_at BETWEEN ? AND ?", userID, start, end).
		Scan(&total).Error
	return total, err
}
