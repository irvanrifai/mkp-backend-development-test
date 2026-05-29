package schedule

import (
	"errors"
	"testing"
	"time"

	"github.com/irvanrifai/mkp-backend-development-test/api-ticketing/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type scheduleRepoMock struct {
	mock.Mock
}

func (m *scheduleRepoMock) Create(schedule *models.Schedule) error {
	args := m.Called(schedule)
	return args.Error(0)
}

func (m *scheduleRepoMock) FindByID(id uint) (models.Schedule, error) {
	args := m.Called(id)
	return args.Get(0).(models.Schedule), args.Error(1)
}

func (m *scheduleRepoMock) FindAll() ([]models.Schedule, error) {
	args := m.Called()
	return args.Get(0).([]models.Schedule), args.Error(1)
}

func (m *scheduleRepoMock) Update(schedule *models.Schedule) error {
	args := m.Called(schedule)
	return args.Error(0)
}

func (m *scheduleRepoMock) Delete(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func TestCreateSchedule_Success(t *testing.T) {
	repo := &scheduleRepoMock{}
	uc := NewScheduleUsecase(repo)

	repo.On("Create", mock.AnythingOfType("*models.Schedule")).Return(nil)

	payload := models.Schedule{
		MovieID:        1,
		StudioID:       2,
		ShowTime:       time.Now(),
		PricePerTicket: 100000,
		Status:         "ACTIVE",
	}

	res, err := uc.Create(payload)

	assert.NoError(t, err)
	assert.Equal(t, payload.MovieID, res.MovieID)
	assert.Equal(t, payload.StudioID, res.StudioID)
	assert.Equal(t, payload.Status, res.Status)
}

func TestCreateSchedule_DatabaseError(t *testing.T) {
	repo := &scheduleRepoMock{}
	uc := NewScheduleUsecase(repo)

	repo.On("Create", mock.Anything).Return(errors.New("db error"))

	_, err := uc.Create(models.Schedule{
		MovieID:        1,
		StudioID:       2,
		ShowTime:       time.Now(),
		PricePerTicket: 100000,
		Status:         "ACTIVE",
	})

	assert.Error(t, err)
	assert.Equal(t, "db error", err.Error())
}

func TestUpdateSchedule_Success(t *testing.T) {
	repo := &scheduleRepoMock{}
	uc := NewScheduleUsecase(repo)

	existing := models.Schedule{
		ID:             1,
		MovieID:        1,
		StudioID:       1,
		ShowTime:       time.Now().Add(-24 * time.Hour),
		PricePerTicket: 90000,
		Status:         "ACTIVE",
	}

	repo.On("FindByID", uint(1)).Return(existing, nil)
	repo.On("Update", mock.AnythingOfType("*models.Schedule")).Return(nil)

	updatedPayload := models.Schedule{
		MovieID:        2,
		StudioID:       3,
		ShowTime:       time.Now(),
		PricePerTicket: 120000,
		Status:         "INACTIVE",
	}

	res, err := uc.Update(1, updatedPayload)

	assert.NoError(t, err)
	assert.Equal(t, uint(1), res.ID)
	assert.Equal(t, updatedPayload.MovieID, res.MovieID)
	assert.Equal(t, updatedPayload.StudioID, res.StudioID)
	assert.Equal(t, updatedPayload.Status, res.Status)
}

func TestUpdateSchedule_FindByIDError(t *testing.T) {
	repo := &scheduleRepoMock{}
	uc := NewScheduleUsecase(repo)

	repo.On("FindByID", uint(1)).Return(models.Schedule{}, errors.New("not found"))

	_, err := uc.Update(1, models.Schedule{})

	assert.Error(t, err)
	assert.Equal(t, "not found", err.Error())
}
