package event

import (
	"github.com/irvanrifai/mkp-backend-development-test/api-ticketing/models"
	"gorm.io/gorm"
)

type IEventRepository interface {
	Create(event *models.Event) error
	FindByID(id uint) (models.Event, error)
	FindAll() ([]models.Event, error)
	Update(event *models.Event) error
	Delete(id uint) error
}

type eventRepository struct {
	db *gorm.DB
}

func NewEventRepository(db *gorm.DB) IEventRepository {
	return &eventRepository{db: db}
}

func (r *eventRepository) Create(event *models.Event) error {
	return r.db.Create(event).Error
}

func (r *eventRepository) FindByID(id uint) (models.Event, error) {
	var event models.Event
	err := r.db.First(&event, id).Error
	return event, err
}

func (r *eventRepository) FindAll() ([]models.Event, error) {
	var events []models.Event
	err := r.db.Find(&events).Error
	return events, err
}

func (r *eventRepository) Update(event *models.Event) error {
	return r.db.Save(event).Error
}

func (r *eventRepository) Delete(id uint) error {
	return r.db.Delete(&models.Event{}, id).Error
}
