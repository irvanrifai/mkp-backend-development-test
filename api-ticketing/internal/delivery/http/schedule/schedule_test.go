package schedule

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/irvanrifai/mkp-backend-development-test/api-ticketing/mocks"
	"github.com/irvanrifai/mkp-backend-development-test/api-ticketing/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestScheduleHandler_Update_Success(t *testing.T) {
	app := fiber.New()
	usecase := mocks.NewIScheduleUsecase(t)
	handler := NewScheduleHandler(usecase)
	app.Put("/schedules/:id", handler.Update)

	payload := map[string]any{
		"movie_id":         2,
		"studio_id":        3,
		"show_time":        "2026-05-29T15:00:00Z",
		"price_per_ticket": 125000,
		"status":           "ACTIVE",
	}
	body, _ := json.Marshal(payload)

	updatedSchedule := models.Schedule{
		ID:             1,
		MovieID:        2,
		StudioID:       3,
		ShowTime:       time.Date(2026, 5, 29, 15, 0, 0, 0, time.UTC),
		PricePerTicket: 125000,
		Status:         "ACTIVE",
	}

	usecase.EXPECT().
		Update(uint(1), mock.Anything).
		Return(updatedSchedule, nil)

	req := httptest.NewRequest(http.MethodPut, "/schedules/1", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var got models.Schedule
	err = json.NewDecoder(resp.Body).Decode(&got)
	assert.NoError(t, err)
	assert.Equal(t, updatedSchedule.ID, got.ID)
	assert.Equal(t, updatedSchedule.MovieID, got.MovieID)
	assert.Equal(t, updatedSchedule.StudioID, got.StudioID)
	assert.Equal(t, updatedSchedule.Status, got.Status)
}

func TestScheduleHandler_Update_InvalidID(t *testing.T) {
	app := fiber.New()
	usecase := mocks.NewIScheduleUsecase(t)
	handler := NewScheduleHandler(usecase)
	app.Put("/schedules/:id", handler.Update)

	req := httptest.NewRequest(http.MethodPut, "/schedules/abc", bytes.NewReader([]byte(`{"movie_id":1,"studio_id":1,"show_time":"2026-05-29T15:00:00Z","price_per_ticket":120000}`)))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}
