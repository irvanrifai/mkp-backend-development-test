package schedule

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/irvanrifai/mkp-backend-development-test/api-ticketing/internal/usecase/schedule"
	"github.com/irvanrifai/mkp-backend-development-test/api-ticketing/models"
	"github.com/irvanrifai/mkp-backend-development-test/api-ticketing/pkg/utils"
)

type ScheduleHandler struct {
	usecase schedule.IScheduleUsecase
}

func NewScheduleHandler(uc schedule.IScheduleUsecase) *ScheduleHandler {
	return &ScheduleHandler{usecase: uc}
}

func (h *ScheduleHandler) Create(c *fiber.Ctx) error {
	type CreateRequest struct {
		MovieID        uint      `json:"movie_id" validate:"required"`
		StudioID       uint      `json:"studio_id" validate:"required"`
		ShowTime       time.Time `json:"show_time" validate:"required"`
		PricePerTicket float64   `json:"price_per_ticket" validate:"required"`
		Status         string    `json:"status"`
	}

	req := new(CreateRequest)
	if err := c.BodyParser(req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Bad Request"})
	}

	if errs := utils.Validate(req); len(errs) > 0 {
		return c.Status(422).JSON(fiber.Map{"errors": errs})
	}

	if req.Status == "" {
		req.Status = "ACTIVE"
	}

	res, err := h.usecase.Create(models.Schedule{
		MovieID:        req.MovieID,
		StudioID:       req.StudioID,
		ShowTime:       req.ShowTime,
		PricePerTicket: req.PricePerTicket,
		Status:         req.Status,
	})
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(201).JSON(res)
}

func (h *ScheduleHandler) List(c *fiber.Ctx) error {
	res, err := h.usecase.FindAll()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(200).JSON(res)
}

func (h *ScheduleHandler) Detail(c *fiber.Ctx) error {
	id := c.Params("id")
	scheduleID, err := utils.StringToUint(id)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid schedule ID"})
	}

	res, err := h.usecase.FindByID(scheduleID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(200).JSON(res)
}

func (h *ScheduleHandler) Update(c *fiber.Ctx) error {
	id := c.Params("id")
	scheduleID, err := utils.StringToUint(id)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid schedule ID"})
	}

	type UpdateRequest struct {
		MovieID        uint      `json:"movie_id" validate:"required"`
		StudioID       uint      `json:"studio_id" validate:"required"`
		ShowTime       time.Time `json:"show_time" validate:"required"`
		PricePerTicket float64   `json:"price_per_ticket" validate:"required"`
		Status         string    `json:"status"`
	}

	req := new(UpdateRequest)
	if err := c.BodyParser(req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Bad Request"})
	}

	if errs := utils.Validate(req); len(errs) > 0 {
		return c.Status(422).JSON(fiber.Map{"errors": errs})
	}

	if req.Status == "" {
		req.Status = "ACTIVE"
	}

	res, err := h.usecase.Update(scheduleID, models.Schedule{
		MovieID:        req.MovieID,
		StudioID:       req.StudioID,
		ShowTime:       req.ShowTime,
		PricePerTicket: req.PricePerTicket,
		Status:         req.Status,
	})
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(200).JSON(res)
}

func (h *ScheduleHandler) Delete(c *fiber.Ctx) error {
	id := c.Params("id")
	scheduleID, err := utils.StringToUint(id)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid schedule ID"})
	}

	err = h.usecase.Delete(scheduleID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(200).JSON(fiber.Map{"message": "ok"})
}
