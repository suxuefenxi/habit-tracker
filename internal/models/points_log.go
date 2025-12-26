package models

import "time"

type UserPointsLog struct {
	ID             uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID         uint64    `gorm:"column:user_id;not null;index" json:"user_id"`
	ChangeAmount   int       `gorm:"column:change_amount;not null" json:"change_amount"`
	Reason         string    `gorm:"column:reason;type:varchar(32);not null" json:"reason"`
	RelatedHabitID *uint64   `gorm:"column:related_habit_id;index" json:"related_habit_id"`
	CreatedAt      time.Time `gorm:"column:created_at;not null;index" json:"created_at"`
}

func (UserPointsLog) TableName() string { return "user_points_log" }
