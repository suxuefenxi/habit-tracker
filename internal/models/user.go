package models

import "time"

type User struct {
	ID            uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	Username      string    `gorm:"column:username;type:varchar(64);not null;uniqueIndex" json:"username"`
	PasswordHash  string    `gorm:"column:password_hash;type:varchar(255);not null" json:"-"`
	Nickname      string    `gorm:"column:nickname;type:varchar(64)" json:"nickname"`
	Points        int64     `gorm:"column:points;not null;default:0" json:"points"`
	TotalCheckins int64     `gorm:"column:total_checkins;not null;default:0" json:"total_checkins"`
	CreatedAt     time.Time `gorm:"column:created_at;not null" json:"created_at"`
}

func (User) TableName() string { return "users" }
