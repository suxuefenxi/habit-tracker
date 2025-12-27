package service

import (
	"context"
	"time"

	"habit-tracker/internal/models"
	"habit-tracker/internal/repository"
)

type AchievementMetrics struct {
	CurrentStreakDays int
	TotalCheckins     int
	TotalPoints       int64
}

type AchievementService struct {
	achievements *repository.AchievementRepository
	userAch      *repository.UserAchievementRepository
}

func NewAchievementService(ach *repository.AchievementRepository, userAch *repository.UserAchievementRepository) *AchievementService {
	return &AchievementService{achievements: ach, userAch: userAch}
}

// EvaluateAndUnlock checks achievements and inserts newly unlocked ones.
func (s *AchievementService) EvaluateAndUnlock(ctx context.Context, userID uint64, m AchievementMetrics) ([]models.UserAchievement, error) {
	all, err := s.achievements.ListAll(ctx)
	if err != nil {
		return nil, err
	}
	unlocked, err := s.userAch.ListByUser(ctx, userID)
	if err != nil {
		return nil, err
	}
	unlockedSet := make(map[uint64]struct{}, len(unlocked))
	for _, ua := range unlocked {
		unlockedSet[ua.AchievementID] = struct{}{}
	}

	var newly []models.UserAchievement
	now := time.Now()
	for _, ach := range all {
		if _, ok := unlockedSet[ach.ID]; ok {
			continue
		}
		if !s.meetCondition(ach, m) {
			continue
		}
		ua := models.UserAchievement{
			UserID:        userID,
			AchievementID: ach.ID,
			UnlockedAt:    now,
		}
		if err := s.userAch.Create(ctx, &ua); err != nil {
			return newly, err
		}
		newly = append(newly, ua)
	}
	return newly, nil
}

func (s *AchievementService) meetCondition(ach models.Achievement, m AchievementMetrics) bool {
	switch ach.ConditionType {
	case "streak_days":
		return m.CurrentStreakDays >= ach.ConditionValue
	case "total_checkins":
		return m.TotalCheckins >= ach.ConditionValue
	case "points":
		return m.TotalPoints >= int64(ach.ConditionValue)
	default:
		return false
	}
}
