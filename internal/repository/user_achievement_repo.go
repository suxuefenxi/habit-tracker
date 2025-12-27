package repository

import (
	"context"

	"gorm.io/gorm"

	"habit-tracker/internal/models"
)

type UserAchievementRepository struct {
	db *gorm.DB
}

func NewUserAchievementRepository(db *gorm.DB) *UserAchievementRepository {
	return &UserAchievementRepository{db: db}
}

func (r *UserAchievementRepository) ListByUser(ctx context.Context, userID uint64) ([]models.UserAchievement, error) {
	var items []models.UserAchievement
	err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("unlocked_at desc").
		Find(&items).Error
	return items, err
}

func (r *UserAchievementRepository) Create(ctx context.Context, ua *models.UserAchievement) error {
	return r.db.WithContext(ctx).Create(ua).Error
}
