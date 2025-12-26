package models

import "time"

type UserAchievement struct {
	ID            uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID        uint64    `gorm:"column:user_id;not null;index;uniqueIndex:uq_user_ach" json:"user_id"`
	AchievementID uint64    `gorm:"column:achievement_id;not null;index;uniqueIndex:uq_user_ach" json:"achievement_id"`
	UnlockedAt    time.Time `gorm:"column:unlocked_at;not null" json:"unlocked_at"`
}

func (UserAchievement) TableName() string { return "user_achievements" }
