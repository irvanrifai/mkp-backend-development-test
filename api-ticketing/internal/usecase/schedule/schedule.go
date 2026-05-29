package schedule

import (
	"github.com/irvanrifai/mkp-backend-development-test/api-ticketing/internal/repository/schedule"
	"github.com/irvanrifai/mkp-backend-development-test/api-ticketing/models"
)

type IScheduleUsecase interface {
	Create(models.Schedule) (models.Schedule, error)
	FindAll() ([]models.Schedule, error)
	FindByID(id uint) (models.Schedule, error)
	Update(id uint, schedule models.Schedule) (models.Schedule, error)
	Delete(id uint) error
}

type scheduleUsecase struct {
	repo schedule.IScheduleRepository
}

func NewScheduleUsecase(repo schedule.IScheduleRepository) IScheduleUsecase {
	return &scheduleUsecase{repo}
}

func (u *scheduleUsecase) Create(schedule models.Schedule) (models.Schedule, error) {
	err := u.repo.Create(&schedule)
	return schedule, err
}

func (u *scheduleUsecase) FindAll() ([]models.Schedule, error) {
	return u.repo.FindAll()
}

func (u *scheduleUsecase) FindByID(id uint) (models.Schedule, error) {
	return u.repo.FindByID(id)
}

func (u *scheduleUsecase) Update(id uint, schedule models.Schedule) (models.Schedule, error) {
	_, err := u.repo.FindByID(id)
	if err != nil {
		return models.Schedule{}, err
	}

	schedule.ID = id
	err = u.repo.Update(&schedule)
	return schedule, err
}

func (u *scheduleUsecase) Delete(id uint) error {
	return u.repo.Delete(id)
}
