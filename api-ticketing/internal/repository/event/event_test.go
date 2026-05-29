package event

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/irvanrifai/mkp-backend-development-test/api-ticketing/models"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func TestRepository_Create(t *testing.T) {
	dbMock, mock, _ := sqlmock.New()
	dialector := postgres.New(postgres.Config{Conn: dbMock})
	db, _ := gorm.Open(dialector, &gorm.Config{})

	repo := NewEventRepository(db)
	event := &models.Event{Title: "Konser Musik"}

	mock.ExpectBegin()
	// GORM PostgreSQL driver menggunakan $1, $2, dst.
	// Sesuaikan jumlah AnyArg dengan jumlah kolom di INSERT statement (biasanya 6)
	mock.ExpectQuery(`^INSERT INTO "events"`).
		WithArgs(
			sqlmock.AnyArg(), // title
			sqlmock.AnyArg(), // description
			sqlmock.AnyArg(), // creator_id
			sqlmock.AnyArg(), // created_at
			sqlmock.AnyArg(), // updated_at
			sqlmock.AnyArg(), // deleted_at
		).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	mock.ExpectCommit()

	err := repo.Create(event)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}
