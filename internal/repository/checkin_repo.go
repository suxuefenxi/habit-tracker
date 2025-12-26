package repository

import (
	"context"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"habit-tracker/internal/models"
)

type CheckinRepository struct {
	db *gorm.DB
}

func NewCheckinRepository(db *gorm.DB) *CheckinRepository {
	return &CheckinRepository{db: db}
}

// Upsert by (habit_id, checkin_date) to avoid duplicate daily records.
func (r *CheckinRepository) Upsert(ctx context.Context, checkin *models.HabitCheckin) error {
	return r.db.WithContext(ctx).
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "habit_id"}, {Name: "checkin_date"}},
			DoUpdates: clause.AssignmentColumns([]string{"count", "user_id"}),
		}).
		Create(checkin).Error
}

func (r *CheckinRepository) ListByHabitAndDateRange(ctx context.Context, habitID uint64, start, end time.Time) ([]models.HabitCheckin, error) {
	var records []models.HabitCheckin
	err := r.db.WithContext(ctx).
		Where("habit_id = ? AND checkin_date BETWEEN ? AND ?", habitID, start, end).
		Order("checkin_date desc").
		Find(&records).Error
	return records, err
}
