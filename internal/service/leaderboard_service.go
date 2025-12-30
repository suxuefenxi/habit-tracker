package service

import (
	"context"
	"sort"
	"time"

	"habit-tracker/internal/repository"
	"habit-tracker/internal/utils"
)

type LeaderboardEntry struct {
	UserID   uint64 `json:"user_id"`
	Nickname string `json:"nickname"`
	Points   int64  `json:"points"`
	Rank     int    `json:"rank"`
}

type LeaderboardService struct {
	users  *repository.UserRepository
	points *repository.PointsRepository
}

func NewLeaderboardService(users *repository.UserRepository, points *repository.PointsRepository) *LeaderboardService {
	return &LeaderboardService{users: users, points: points}
}

func (s *LeaderboardService) Weekly(ctx context.Context) ([]LeaderboardEntry, error) {
	today := utils.TruncateDate(time.Now())
	weekday := today.Weekday()
	// Calculate the start of the current week (Monday 00:00)
	start := today.AddDate(0, 0, -int(weekday)+1)
	if weekday == time.Sunday {
		start = today.AddDate(0, 0, -6) // Special case for Sunday
	}
	end := start.AddDate(0, 0, 7) // End of the week (next Monday 00:00)
	return s.build(ctx, start, end)
}

func (s *LeaderboardService) Monthly(ctx context.Context) ([]LeaderboardEntry, error) {
	today := utils.TruncateDate(time.Now())
	start := time.Date(today.Year(), today.Month(), 1, 0, 0, 0, 0, today.Location())
	end := start.AddDate(0, 1, 0)
	return s.build(ctx, start, end)
}

func (s *LeaderboardService) build(ctx context.Context, start, end time.Time) ([]LeaderboardEntry, error) {
	users, err := s.users.ListAll(ctx)
	if err != nil {
		return nil, err
	}

	entries := make([]LeaderboardEntry, 0, len(users))
	for _, u := range users {
		sum, err := s.points.SumByUserAndRange(ctx, u.ID, start, end)
		if err != nil {
			return nil, err
		}
		if sum == 0 { // 目前的处理：没有积分变动的用户不上榜
			continue
		}
		entries = append(entries, LeaderboardEntry{
			UserID:   u.ID,
			Nickname: u.Nickname,
			Points:   sum,
		})
	}

	sort.Slice(entries, func(i, j int) bool {
		if entries[i].Points == entries[j].Points {
			return entries[i].UserID < entries[j].UserID
		}
		return entries[i].Points > entries[j].Points
	})

	for i := range entries {
		entries[i].Rank = i + 1
	}
	return entries, nil
}
