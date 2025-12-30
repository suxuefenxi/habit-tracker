package service

import (
	"context"
	"errors"
	"time"

	"gorm.io/gorm"

	"habit-tracker/internal/models"
	"habit-tracker/internal/repository"
)

const baseCheckinPoints = 1

type CheckinService struct {
	habitRepo          *repository.HabitRepository
	userRepo           *repository.UserRepository
	checkinRepo        *repository.CheckinRepository
	pointService       *PointsService
	achievementService *AchievementService
}

var (
	ErrCheckinForbidden = errors.New("habit does not belong to user")
	ErrHabitMissing     = gorm.ErrRecordNotFound
)

type CheckinResult struct {
	TodayCount     int
	ReachedTarget  bool
	StreakDays     int
	TotalCheckins  int
	PointsAwarded  int
	UnlockedAwards []models.UserAchievement
}

func NewCheckinService(habitRepo *repository.HabitRepository, users *repository.UserRepository, checkins *repository.CheckinRepository, points *PointsService, achievements *AchievementService) *CheckinService {
	return &CheckinService{habitRepo: habitRepo, userRepo: users, checkinRepo: checkins, pointService: points, achievementService: achievements}
}

func (s *CheckinService) Checkin(ctx context.Context, userID, habitID uint64, countInc int) (*CheckinResult, error) {
	if countInc <= 0 {
		return nil, errors.New("check-in count must be greater than 0")
	}

	habit, err := s.getOwnedHabit(ctx, userID, habitID)
	if err != nil {
		return nil, err
	}

	today := todayDate()
	if err := s.upsertToday(ctx, userID, habitID, countInc, today); err != nil {
		return nil, err
	}

	todayRec, err := s.checkinRepo.GetByHabitAndDate(ctx, habitID, today)
	if err != nil {
		return nil, err
	}

	reached := todayRec.Count >= habit.TargetTimes
	var pointsAwarded int
	if reached {
		if todayRec.Count == countInc { // first time reaching target today
			pointsAwarded, err = s.awardPoints(ctx, userID, habitID)
			if err != nil {
				return nil, err
			}
		}
		if err := s.userRepo.IncrementCheckins(ctx, userID, 1); err != nil {
			return nil, err
		}
	}

	streak := s.calculateStreak(ctx, habitID, habit.TargetTimes, today)
	totalCheckins, err := s.checkinRepo.SumCountByHabit(ctx, habitID)
	if err != nil {
		return nil, err
	}
	totalPoints, err := s.pointService.GetUserPoints(ctx, userID)
	if err != nil {
		return nil, err
	}

	newly, err := s.achievementService.EvaluateAndUnlock(ctx, userID, AchievementMetrics{
		CurrentStreakDays: streak,
		TotalCheckins:     int(totalCheckins),
		TotalPoints:       totalPoints,
	})
	if err != nil {
		return nil, err
	}

	return &CheckinResult{
		TodayCount:     todayRec.Count,
		ReachedTarget:  reached,
		StreakDays:     streak,
		TotalCheckins:  int(totalCheckins),
		PointsAwarded:  pointsAwarded,
		UnlockedAwards: newly,
	}, nil
}

func (s *CheckinService) getOwnedHabit(ctx context.Context, userID, habitID uint64) (*models.Habit, error) {
	habit, err := s.habitRepo.GetByID(ctx, habitID)
	if err != nil {
		return nil, err
	}
	if habit.UserID != userID {
		return nil, ErrCheckinForbidden
	}
	return habit, nil
}

func todayDate() time.Time {
	return time.Now().In(time.Local).Truncate(24 * time.Hour)
}

func (s *CheckinService) upsertToday(ctx context.Context, userID, habitID uint64, countInc int, today time.Time) error {
	rec := &models.HabitCheckin{
		HabitID:     habitID,
		UserID:      userID,
		CheckinDate: today,
		Count:       countInc,
		CreatedAt:   time.Now(),
	}
	return s.checkinRepo.Upsert(ctx, rec)
}

func (s *CheckinService) awardPoints(ctx context.Context, userID, habitID uint64) (int, error) {
	delta := int64(baseCheckinPoints)
	if err := s.pointService.AddPoints(ctx, userID, delta, "checkin", &habitID); err != nil {
		return 0, err
	}
	return int(delta), nil
}

func (s *CheckinService) ListHistory(ctx context.Context, userID, habitID uint64, start, end time.Time) ([]models.HabitCheckin, error) {
	habit, err := s.habitRepo.GetByID(ctx, habitID)
	if err != nil {
		return nil, err
	}
	if habit.UserID != userID {
		return nil, ErrCheckinForbidden
	}
	return s.checkinRepo.ListByHabitAndDateRange(ctx, habitID, start, end)
}

func (s *CheckinService) calculateStreak(ctx context.Context, habitID uint64, targetTimes int, today time.Time) int {
	records, err := s.checkinRepo.ListByHabitDesc(ctx, habitID)
	if err != nil {
		return 0
	}
	return countConsecutiveFromToday(records, targetTimes, today)
}

func countConsecutiveFromToday(records []models.HabitCheckin, targetTimes int, today time.Time) int {
	streak := 0
	expected := today
	for _, rec := range records {
		if rec.CheckinDate.After(expected) { // future date, skip
			continue
		}
		if !sameDate(rec.CheckinDate, expected) {
			break
		}
		if rec.Count < targetTimes {
			break
		}
		streak++
		expected = expected.AddDate(0, 0, -1)
	}
	return streak
}

func sameDate(a, b time.Time) bool {
	ay, am, ad := a.Date()
	by, bm, bd := b.Date()
	return ay == by && am == bm && ad == bd
}
