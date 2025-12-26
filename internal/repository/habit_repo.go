package repository

import (
	"context"

	"gorm.io/gorm"

	"habit-tracker/internal/models"
)

type HabitRepository struct {
	db *gorm.DB
}

func NewHabitRepository(db *gorm.DB) *HabitRepository {
	return &HabitRepository{db: db}
}

func (r *HabitRepository) ListByUser(ctx context.Context, userID uint64) ([]models.Habit, error) {
	var habits []models.Habit
	err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("start_date asc").
		Find(&habits).Error
	return habits, err
}

func (r *HabitRepository) GetByID(ctx context.Context, id uint64) (*models.Habit, error) {
	var habit models.Habit
	if err := r.db.WithContext(ctx).First(&habit, id).Error; err != nil {
		return nil, err
	}
	return &habit, nil
}

func (r *HabitRepository) Create(ctx context.Context, habit *models.Habit) error {
	return r.db.WithContext(ctx).Create(habit).Error
}

func (r *HabitRepository) Update(ctx context.Context, habit *models.Habit) error {
	return r.db.WithContext(ctx).Save(habit).Error
}
