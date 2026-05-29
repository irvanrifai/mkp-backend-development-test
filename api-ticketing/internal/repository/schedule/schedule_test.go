package schedule

import (
	"testing"
	"time"

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

	repo := NewScheduleRepository(db)
	schedule := &models.Schedule{
		MovieID:        1,
		StudioID:       1,
		ShowTime:       time.Now(),
		PricePerTicket: 125000,
		Status:         "ACTIVE",
	}

	mock.ExpectBegin()
	// GORM PostgreSQL driver menggunakan $1, $2, dst.
	// Sesuaikan jumlah AnyArg dengan jumlah kolom di INSERT statement.
	mock.ExpectQuery(`^INSERT INTO "schedules"`).
		WithArgs(
			sqlmock.AnyArg(), // movie_id
			sqlmock.AnyArg(), // studio_id
			sqlmock.AnyArg(), // show_time
			sqlmock.AnyArg(), // price_per_ticket
			sqlmock.AnyArg(), // status
			sqlmock.AnyArg(), // created_at
			sqlmock.AnyArg(), // updated_at
			sqlmock.AnyArg(), // deleted_at
		).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	mock.ExpectCommit()

	err := repo.Create(schedule)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}
