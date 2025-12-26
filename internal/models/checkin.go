package models

import "time"

type HabitCheckin struct {
	ID          uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	HabitID     uint64    `gorm:"column:habit_id;not null;index;uniqueIndex:uq_habit_date" json:"habit_id"`
	UserID      uint64    `gorm:"column:user_id;not null;index" json:"user_id"`
	CheckinDate time.Time `gorm:"column:checkin_date;type:date;not null;index;uniqueIndex:uq_habit_date" json:"checkin_date"`
	Count       int       `gorm:"column:count;not null;default:0" json:"count"`
	CreatedAt   time.Time `gorm:"column:created_at;not null" json:"created_at"`
}

func (HabitCheckin) TableName() string { return "habit_checkins" }
