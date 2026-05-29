package event

import (
	"github.com/irvanrifai/mkp-backend-development-test/api-ticketing/internal/repository/event"
	"github.com/irvanrifai/mkp-backend-development-test/api-ticketing/models"
)

type IEventUsecase interface {
	Create(title, desc string, creatorID uint) (models.Event, error)
	FindAll() ([]models.Event, error)
	FindByID(id uint) (models.Event, error)
	Update(id uint, title, desc string) (models.Event, error)
	Delete(id uint) error
}

type eventUsecase struct {
	repo event.IEventRepository
}

func NewEventUsecase(repo event.IEventRepository) IEventUsecase {
	return &eventUsecase{repo}
}

func (u *eventUsecase) Create(title, desc string, creatorID uint) (models.Event, error) {
	event := models.Event{
		Title:       title,
		Description: desc,
		CreatorID:   creatorID,
	}
	err := u.repo.Create(&event)
	return event, err
}

func (u *eventUsecase) FindAll() ([]models.Event, error) {
	return u.repo.FindAll()
}

func (u *eventUsecase) FindByID(id uint) (models.Event, error) {
	return u.repo.FindByID(id)
}

func (u *eventUsecase) Update(id uint, title, desc string) (models.Event, error) {
	event, err := u.repo.FindByID(id)
	if err != nil {
		return models.Event{}, err
	}

	event.Title = title
	event.Description = desc

	err = u.repo.Update(&event)
	return event, err
}

func (u *eventUsecase) Delete(id uint) error {
	return u.repo.Delete(id)
}
