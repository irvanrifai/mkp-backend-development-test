package auth

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/irvanrifai/mkp-backend-development-test/api-ticketing/mocks"
	"github.com/irvanrifai/mkp-backend-development-test/api-ticketing/models"
	"github.com/stretchr/testify/assert"
)

func TestAuthHandler_Register_Success(t *testing.T) {
	app := fiber.New()
	usecase := mocks.NewIUserUsecase(t)
	handler := NewAuthHandler(usecase)
	app.Post("/auth/register", handler.Register)

	payload := map[string]any{
		"name":     "John Doe",
		"username": "johndoe",
		"email":    "john@example.com",
		"password": "password",
	}
	body, _ := json.Marshal(payload)

	usecase.EXPECT().
		Register("John Doe", "johndoe", "john@example.com", "password").
		Return(models.User{ID: 1, Name: "John Doe", Email: "john@example.com", Username: "johndoe"}, nil)

	req := httptest.NewRequest(http.MethodPost, "/auth/register", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	var got models.User
	err = json.NewDecoder(resp.Body).Decode(&got)
	assert.NoError(t, err)
	assert.Equal(t, "John Doe", got.Name)
	assert.Equal(t, "john@example.com", got.Email)
}

func TestAuthHandler_Login_Success(t *testing.T) {
	app := fiber.New()
	usecase := mocks.NewIUserUsecase(t)
	handler := NewAuthHandler(usecase)
	app.Post("/auth/login", handler.Login)

	payload := map[string]any{
		"email":    "john@example.com",
		"password": "password",
	}
	body, _ := json.Marshal(payload)

	usecase.EXPECT().
		Login("john@example.com", "password").
		Return("token-value", nil)

	req := httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var got map[string]string
	err = json.NewDecoder(resp.Body).Decode(&got)
	assert.NoError(t, err)
	assert.Equal(t, "token-value", got["token"])
}
