package service

import (
	"context"
	"sort"
	"time"

	"habit-tracker/internal/models"
	"habit-tracker/internal/repository"
	"habit-tracker/internal/utils"
)

type UserStats struct {
	TotalCheckins   int64 `json:"total_checkins"`
	WeeklyCheckins  int64 `json:"weekly_checkins"`
	MonthlyCheckins int64 `json:"monthly_checkins"`
	WeeklyPoints    int64 `json:"weekly_points"`
	MonthlyPoints   int64 `json:"monthly_points"`
	LongestStreak   int   `json:"longest_streak"`
}

type UserStatsService struct {
	users    *repository.UserRepository
	habits   *repository.HabitRepository
	checkins *repository.CheckinRepository
	points   *PointsService
}

func NewUserStatsService(users *repository.UserRepository, habits *repository.HabitRepository, checkins *repository.CheckinRepository, points *PointsService) *UserStatsService {
	return &UserStatsService{users: users, habits: habits, checkins: checkins, points: points}
}

func (s *UserStatsService) GetStats(ctx context.Context, userID uint64) (*UserStats, error) {
	today := utils.TruncateDate(time.Now())

	weekStart := startOfWeek(today)
	weekEnd := weekStart.AddDate(0, 0, 7)

	monthStart := time.Date(today.Year(), today.Month(), 1, 0, 0, 0, 0, today.Location())
	monthEnd := monthStart.AddDate(0, 1, 0)

	totalCheckins, err := s.users.GetTotalCheckins(ctx, userID)
	if err != nil {
		return nil, err
	}

	weeklyCheckins, err := s.checkins.SumCountByUserAndRange(ctx, userID, weekStart, weekEnd)
	if err != nil {
		return nil, err
	}

	monthlyCheckins, err := s.checkins.SumCountByUserAndRange(ctx, userID, monthStart, monthEnd)
	if err != nil {
		return nil, err
	}

	weeklyPoints, err := s.points.SumByRange(ctx, userID, weekStart, weekEnd)
	if err != nil {
		return nil, err
	}

	monthlyPoints, err := s.points.SumByRange(ctx, userID, monthStart, monthEnd)
	if err != nil {
		return nil, err
	}

	longestStreak, err := s.longestStreak(ctx, userID)
	if err != nil {
		return nil, err
	}

	return &UserStats{
		TotalCheckins:   totalCheckins,
		WeeklyCheckins:  weeklyCheckins,
		MonthlyCheckins: monthlyCheckins,
		WeeklyPoints:    weeklyPoints,
		MonthlyPoints:   monthlyPoints,
		LongestStreak:   longestStreak,
	}, nil
}

func (s *UserStatsService) longestStreak(ctx context.Context, userID uint64) (int, error) {
	habits, err := s.habits.ListByUser(ctx, userID)
	if err != nil {
		return 0, err
	}

	maxStreak := 0
	for _, h := range habits {
		records, err := s.checkins.ListByHabitDesc(ctx, h.ID)
		if err != nil {
			return 0, err
		}
		streak := longestStreakForHabit(records, h.TargetTimes)
		if streak > maxStreak {
			maxStreak = streak
		}
	}
	return maxStreak, nil
}

func longestStreakForHabit(records []models.HabitCheckin, target int) int {
	if target <= 0 {
		return 0
	}
	if len(records) == 0 {
		return 0
	}

	sort.Slice(records, func(i, j int) bool {
		return records[i].CheckinDate.Before(records[j].CheckinDate)
	})

	best, current := 0, 0
	var prevDate time.Time

	for _, rec := range records {
		if rec.Count < target {
			current = 0
			continue
		}

		if current == 0 {
			current = 1
		} else {
			delta := daysBetween(prevDate, rec.CheckinDate)
			if delta == 1 {
				current++
			} else if delta > 1 {
				current = 1
			}
		}

		prevDate = rec.CheckinDate
		if current > best {
			best = current
		}
	}

	return best
}

func startOfWeek(t time.Time) time.Time {
	weekday := t.Weekday()
	if weekday == time.Sunday {
		return t.AddDate(0, 0, -6)
	}
	return t.AddDate(0, 0, -int(weekday)+1)
}

func daysBetween(a, b time.Time) int {
	ay, am, ad := a.Date()
	by, bm, bd := b.Date()
	da := time.Date(ay, am, ad, 0, 0, 0, 0, a.Location())
	db := time.Date(by, bm, bd, 0, 0, 0, 0, b.Location())
	return int(db.Sub(da).Hours() / 24)
}
