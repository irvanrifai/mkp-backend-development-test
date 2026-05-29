package user

import (
	"testing"

	"github.com/irvanrifai/mkp-backend-development-test/api-ticketing/mocks"
	"github.com/irvanrifai/mkp-backend-development-test/api-ticketing/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
)

func TestUserUsecase_Register_Success(t *testing.T) {
	repo := mocks.NewIUserRepository(t)
	uc := NewUserUsecase(repo, "secret")

	repo.EXPECT().
		Create(mock.AnythingOfType("*models.User")).
		Return(nil)

	result, err := uc.Register("John Doe", "johndoe", "john@example.com", "password")

	assert.NoError(t, err)
	assert.Equal(t, "John Doe", result.Name)
	assert.Equal(t, "johndoe", result.Username)
	assert.Equal(t, "john@example.com", result.Email)
}

func TestUserUsecase_Login_Success(t *testing.T) {
	repo := mocks.NewIUserRepository(t)
	uc := NewUserUsecase(repo, "secret")

	hashed, err := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)
	assert.NoError(t, err)

	repo.EXPECT().
		FindByEmail("john@example.com").
		Return(models.User{ID: 1, Email: "john@example.com", Password: string(hashed)}, nil)

	token, err := uc.Login("john@example.com", "password")

	assert.NoError(t, err)
	assert.NotEmpty(t, token)
}
