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
			Columns: []clause.Column{{Name: "habit_id"}, {Name: "checkin_date"}},
			DoUpdates: clause.Assignments(map[string]interface{}{
				"count":   gorm.Expr("habit_checkins.count + EXCLUDED.count"),
				"user_id": gorm.Expr("EXCLUDED.user_id"),
			}),
		}).
		Create(checkin).Error
}

func (r *CheckinRepository) GetByHabitAndDate(ctx context.Context, habitID uint64, date time.Time) (*models.HabitCheckin, error) {
	var rec models.HabitCheckin
	if err := r.db.WithContext(ctx).
		Where("habit_id = ? AND checkin_date = ?", habitID, date).
		First(&rec).Error; err != nil {
		return nil, err
	}
	return &rec, nil
}

func (r *CheckinRepository) ListByHabitAndDateRange(ctx context.Context, habitID uint64, start, end time.Time) ([]models.HabitCheckin, error) {
	var records []models.HabitCheckin
	err := r.db.WithContext(ctx).
		Where("habit_id = ? AND checkin_date BETWEEN ? AND ?", habitID, start, end).
		Order("checkin_date desc").
		Find(&records).Error
	return records, err
}

func (r *CheckinRepository) SumCountByHabit(ctx context.Context, habitID uint64) (int64, error) {
	var total int64
	err := r.db.WithContext(ctx).
		Model(&models.HabitCheckin{}).
		Select("COALESCE(SUM(count),0)").
		Where("habit_id = ?", habitID).
		Scan(&total).Error
	return total, err
}

func (r *CheckinRepository) SumCountByUser(ctx context.Context, userID uint64) (int64, error) {
	var total int64
	err := r.db.WithContext(ctx).
		Model(&models.HabitCheckin{}).
		Select("COALESCE(SUM(count),0)").
		Where("user_id = ?", userID).
		Scan(&total).Error
	return total, err
}

func (r *CheckinRepository) SumCountByUserAndRange(ctx context.Context, userID uint64, start, end time.Time) (int64, error) {
	var total int64
	err := r.db.WithContext(ctx).
		Model(&models.HabitCheckin{}).
		Select("COALESCE(SUM(count),0)").
		Where("user_id = ? AND checkin_date BETWEEN ? AND ?", userID, start, end).
		Scan(&total).Error
	return total, err
}

func (r *CheckinRepository) ListByHabitDesc(ctx context.Context, habitID uint64) ([]models.HabitCheckin, error) {
	var records []models.HabitCheckin
	err := r.db.WithContext(ctx).
		Where("habit_id = ?", habitID).
		Order("checkin_date desc").
		Find(&records).Error
	return records, err
}
