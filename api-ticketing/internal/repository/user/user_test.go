package user

import (
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/irvanrifai/mkp-backend-development-test/api-ticketing/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func TestUserRepository_Create(t *testing.T) {
	dbMock, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer dbMock.Close()

	dialector := postgres.New(postgres.Config{Conn: dbMock})
	db, err := gorm.Open(dialector, &gorm.Config{})
	require.NoError(t, err)

	repo := NewUserRepository(db)
	user := &models.User{
		Name:     "John Doe",
		Username: "johndoe",
		Email:    "john@example.com",
		Password: "secret",
	}

	mock.ExpectBegin()
	mock.ExpectQuery(`^INSERT INTO "users"`).
		WithArgs(
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
		).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	mock.ExpectCommit()

	err = repo.Create(user)

	assert.NoError(t, err)
	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestUserRepository_FindByEmail(t *testing.T) {
	dbMock, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer dbMock.Close()

	dialector := postgres.New(postgres.Config{Conn: dbMock})
	db, err := gorm.Open(dialector, &gorm.Config{})
	require.NoError(t, err)

	repo := NewUserRepository(db)

	now := time.Now()
	mock.ExpectQuery(`SELECT (.+) FROM "users" WHERE`).
		WithArgs("john@example.com", sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "username", "email", "password", "phone", "created_at", "updated_at", "deleted_at"}).
			AddRow(1, "John Doe", "johndoe", "john@example.com", "secret", "", now, now, nil))

	result, err := repo.FindByEmail("john@example.com")

	assert.NoError(t, err)
	assert.Equal(t, uint(1), result.ID)
	assert.Equal(t, "john@example.com", result.Email)
}
