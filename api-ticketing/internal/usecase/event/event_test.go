package event

import (
	"errors"
	"testing"

	"github.com/irvanrifai/mkp-backend-development-test/api-ticketing/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateEvent_Success(t *testing.T) {
	// Inisialisasi Mock Repository hasil generate Mockery
	repo := mocks.NewIEventRepository(t)
	uc := NewEventUsecase(repo)

	// Setup Ekspektasi menggunakan fitur EXPECT() (Expecter)
	repo.EXPECT().
		Create(mock.AnythingOfType("*models.Event")).
		Return(nil)

	res, err := uc.Create("Workshop Go", "Belajar Clean Arch", 1)

	assert.NoError(t, err)
	assert.Equal(t, "Workshop Go", res.Title)
}

func TestCreateEvent_DatabaseError(t *testing.T) {
	repo := mocks.NewIEventRepository(t)
	uc := NewEventUsecase(repo)

	// Simulasi error database
	repo.EXPECT().
		Create(mock.Anything).
		Return(errors.New("db error"))

	_, err := uc.Create("Title", "Desc", 1)

	assert.Error(t, err)
	assert.Equal(t, "db error", err.Error())
}
