package service

import (
	"context"
	"errors"
	"time"

	"gorm.io/gorm"

	"habit-tracker/internal/models"
	"habit-tracker/internal/repository"
)

var (
	ErrHabitNotFound  = gorm.ErrRecordNotFound
	ErrHabitForbidden = errors.New("habit does not belong to user")
	validTargetTypes  = map[string]struct{}{
		"daily":  {},
		"weekly": {},
		"custom": {},
	}
)

type HabitService struct {
	habitRepo *repository.HabitRepository
}

func NewHabitService(habitRepo *repository.HabitRepository) *HabitService {
	return &HabitService{habitRepo: habitRepo}
}

type HabitInput struct {
	Name        string
	Description string
	TargetType  string
	TargetTimes int
	StartDate   time.Time
	IsActive    *bool // optional for update
}

func (s *HabitService) Create(ctx context.Context, userID uint64, in HabitInput) (*models.Habit, error) {
	if err := validateHabitInput(in, false); err != nil {
		return nil, err
	}

	habit := &models.Habit{
		UserID:      userID,
		Name:        in.Name,
		Description: in.Description,
		TargetType:  in.TargetType,
		TargetTimes: in.TargetTimes,
		StartDate:   in.StartDate,
		IsActive:    true,
	}
	if err := s.habitRepo.Create(ctx, habit); err != nil {
		return nil, err
	}
	return habit, nil
}

func (s *HabitService) Update(ctx context.Context, userID, habitID uint64, in HabitInput) (*models.Habit, error) {
	if err := validateHabitInput(in, true); err != nil {
		return nil, err
	}

	habit, err := s.habitRepo.GetByID(ctx, habitID)
	if err != nil {
		return nil, err
	}
	if habit.UserID != userID {
		return nil, ErrHabitForbidden
	}

	habit.Name = in.Name
	habit.Description = in.Description
	habit.TargetType = in.TargetType
	habit.TargetTimes = in.TargetTimes
	habit.StartDate = in.StartDate
	if in.IsActive != nil {
		habit.IsActive = *in.IsActive
	}

	if err := s.habitRepo.Update(ctx, habit); err != nil {
		return nil, err
	}
	return habit, nil
}

func (s *HabitService) List(ctx context.Context, userID uint64, isActive *bool) ([]models.Habit, error) {
	return s.habitRepo.ListByUserWithActive(ctx, userID, isActive)
}

func (s *HabitService) Get(ctx context.Context, userID, habitID uint64) (*models.Habit, error) {
	habit, err := s.habitRepo.GetByID(ctx, habitID)
	if err != nil {
		return nil, err
	}
	if habit.UserID != userID {
		return nil, ErrHabitForbidden
	}
	return habit, nil
}

func (s *HabitService) SetActive(ctx context.Context, userID, habitID uint64, active bool) error {
	habit, err := s.habitRepo.GetByID(ctx, habitID)
	if err != nil {
		return err
	}
	if habit.UserID != userID {
		return ErrHabitForbidden
	}
	habit.IsActive = active
	if err := s.habitRepo.UpdateStatus(ctx, habitID, active); err != nil {
		return err
	}
	return nil
}

func validateHabitInput(in HabitInput, allowZeroStart bool) error {
	if in.Name == "" {
		return errors.New("name is required")
	}
	if _, ok := validTargetTypes[in.TargetType]; !ok {
		return errors.New("invalid target_type")
	}
	if in.TargetTimes <= 0 {
		return errors.New("target_times must be > 0")
	}
	if !allowZeroStart && in.StartDate.IsZero() {
		return errors.New("start_date is required")
	}
	return nil
}
