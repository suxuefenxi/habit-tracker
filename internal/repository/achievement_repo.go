package repository

import (
	"context"

	"gorm.io/gorm"

	"habit-tracker/internal/models"
)

type AchievementRepository struct {
	db *gorm.DB
}

func NewAchievementRepository(db *gorm.DB) *AchievementRepository {
	return &AchievementRepository{db: db}
}

func (r *AchievementRepository) ListAll(ctx context.Context) ([]models.Achievement, error) {
	var items []models.Achievement
	err := r.db.WithContext(ctx).Order("id asc").Find(&items).Error
	return items, err
}

func (r *AchievementRepository) ListByConditionType(ctx context.Context, conditionType string) ([]models.Achievement, error) {
	var items []models.Achievement
	err := r.db.WithContext(ctx).
		Where("condition_type = ?", conditionType).
		Order("id asc").
		Find(&items).Error
	return items, err
}
