package schedule

import (
	"errors"
	"testing"
	"time"

	"github.com/irvanrifai/mkp-backend-development-test/api-ticketing/mocks"
	"github.com/irvanrifai/mkp-backend-development-test/api-ticketing/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateSchedule_Success(t *testing.T) {
	repo := mocks.NewIScheduleRepository(t)
	uc := NewScheduleUsecase(repo)

	payload := models.Schedule{
		MovieID:        1,
		StudioID:       2,
		ShowTime:       time.Now(),
		PricePerTicket: 100000,
		Status:         "ACTIVE",
	}

	repo.EXPECT().
		Create(mock.AnythingOfType("*models.Schedule")).
		Return(nil)

	res, err := uc.Create(payload)

	assert.NoError(t, err)
	assert.Equal(t, payload.MovieID, res.MovieID)
	assert.Equal(t, payload.StudioID, res.StudioID)
	assert.Equal(t, payload.Status, res.Status)
}

func TestCreateSchedule_DatabaseError(t *testing.T) {
	repo := mocks.NewIScheduleRepository(t)
	uc := NewScheduleUsecase(repo)

	repo.EXPECT().
		Create(mock.AnythingOfType("*models.Schedule")).
		Return(errors.New("db error"))

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
	repo := mocks.NewIScheduleRepository(t)
	uc := NewScheduleUsecase(repo)

	existing := models.Schedule{ID: 1}
	updatedPayload := models.Schedule{
		MovieID:        2,
		StudioID:       3,
		ShowTime:       time.Now(),
		PricePerTicket: 120000,
		Status:         "INACTIVE",
	}

	repo.EXPECT().
		FindByID(uint(1)).
		Return(existing, nil)

	repo.EXPECT().
		Update(mock.AnythingOfType("*models.Schedule")).
		Return(nil)

	res, err := uc.Update(1, updatedPayload)

	assert.NoError(t, err)
	assert.Equal(t, uint(1), res.ID)
	assert.Equal(t, updatedPayload.MovieID, res.MovieID)
	assert.Equal(t, updatedPayload.Status, res.Status)
}

func TestUpdateSchedule_FindByIDError(t *testing.T) {
	repo := mocks.NewIScheduleRepository(t)
	uc := NewScheduleUsecase(repo)

	repo.EXPECT().
		FindByID(uint(1)).
		Return(models.Schedule{}, errors.New("not found"))

	_, err := uc.Update(1, models.Schedule{})

	assert.Error(t, err)
	assert.Equal(t, "not found", err.Error())
}
