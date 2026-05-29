package event

import (
	"github.com/gofiber/fiber/v2"
	"github.com/irvanrifai/mkp-backend-development-test/api-ticketing/internal/usecase/event"
	"github.com/irvanrifai/mkp-backend-development-test/api-ticketing/pkg/utils"
)

type EventHandler struct {
	usecase event.IEventUsecase
}

func NewEventHandler(uc event.IEventUsecase) *EventHandler {
	return &EventHandler{usecase: uc}
}

func (h *EventHandler) Create(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint) // Dari JWT Middleware

	type CreateRequest struct {
		Title string `json:"title" validate:"required,min=5"`
		Desc  string `json:"description" validate:"required"`
	}

	req := new(CreateRequest)
	if err := c.BodyParser(req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Bad Request"})
	}

	if errs := utils.Validate(req); len(errs) > 0 {
		return c.Status(422).JSON(fiber.Map{"errors": errs})
	}

	res, err := h.usecase.Create(req.Title, req.Desc, userID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(201).JSON(res)
}

func (h *EventHandler) List(c *fiber.Ctx) error {
	res, err := h.usecase.FindAll()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(201).JSON(res)
}

func (h *EventHandler) Detail(c *fiber.Ctx) error {
	id := c.Params("id")
	eventID, err := utils.StringToUint(id)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid event ID"})
	}

	res, err := h.usecase.FindByID(eventID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(201).JSON(res)
}

func (h *EventHandler) Update(c *fiber.Ctx) error {
	return nil
}

func (h *EventHandler) Delete(c *fiber.Ctx) error {
	id := c.Params("id")
	eventID, err := utils.StringToUint(id)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid event ID"})
	}

	err = h.usecase.Delete(eventID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(200).JSON(fiber.Map{"message": "ok"})
}
