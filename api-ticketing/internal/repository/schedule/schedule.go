package schedule

import (
	"time"

	"github.com/irvanrifai/mkp-backend-development-test/api-ticketing/models"
	"gorm.io/gorm"
)

type IScheduleRepository interface {
	Create(schedule *models.Schedule) error
	FindByID(id uint) (models.Schedule, error)
	FindAll() ([]models.Schedule, error)
	Update(schedule *models.Schedule) error
	Delete(id uint) error
}

type scheduleRepository struct {
	db *gorm.DB
}

func NewScheduleRepository(db *gorm.DB) IScheduleRepository {
	return &scheduleRepository{db: db}
}

func (r *scheduleRepository) Create(schedule *models.Schedule) error {
	return r.db.Create(schedule).Error
}

func (r *scheduleRepository) FindByID(id uint) (models.Schedule, error) {
	var schedule models.Schedule
	err := r.db.First(&schedule, id).Error
	return schedule, err
}

func (r *scheduleRepository) FindAll() ([]models.Schedule, error) {
	var schedules []models.Schedule
	err := r.db.Find(&schedules).Error
	return schedules, err
}

func (r *scheduleRepository) Update(schedule *models.Schedule) error {
	return r.db.Model(&models.Schedule{}).
		Where("id = ?", schedule.ID).
		Updates(map[string]interface{}{
			"movie_id":         schedule.MovieID,
			"studio_id":        schedule.StudioID,
			"show_time":        schedule.ShowTime,
			"price_per_ticket": schedule.PricePerTicket,
			"status":           schedule.Status,
			"updated_at":       time.Now(),
		}).Error
}

func (r *scheduleRepository) Delete(id uint) error {
	return r.db.Delete(&models.Schedule{}, id).Error
}
