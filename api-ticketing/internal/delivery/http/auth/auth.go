package auth

import (
	"github.com/gofiber/fiber/v2"
	"github.com/irvanrifai/mkp-backend-development-test/api-ticketing/internal/usecase/user"
	"github.com/irvanrifai/mkp-backend-development-test/api-ticketing/pkg/utils"
)

type AuthHandler struct {
	usecase user.IUserUsecase
}

func NewAuthHandler(uc user.IUserUsecase) *AuthHandler {
	return &AuthHandler{usecase: uc}
}

func (h *AuthHandler) Register(c *fiber.Ctx) error {
	type RegisterRequest struct {
		Name     string `json:"name" validate:"required"`
		Email    string `json:"email" validate:"required,email"`
		Password string `json:"password" validate:"required,min=6"`
	}
	req := new(RegisterRequest)
	if err := c.BodyParser(req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Bad Request"})
	}

	if errs := utils.Validate(req); len(errs) > 0 {
		return c.Status(422).JSON(fiber.Map{"errors": errs})
	}

	res, err := h.usecase.Register(req.Name, req.Email, req.Password)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(201).JSON(res)
}

func (h *AuthHandler) Login(c *fiber.Ctx) error {
	type LoginRequest struct {
		Email    string `json:"email" validate:"required,email"`
		Password string `json:"password" validate:"required,min=6"`
	}

	req := new(LoginRequest)
	if err := c.BodyParser(req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Bad Request"})
	}

	if errs := utils.Validate(req); len(errs) > 0 {
		return c.Status(422).JSON(fiber.Map{"errors": errs})
	}

	res, err := h.usecase.Login(req.Email, req.Password)
	if err != nil {
		return c.Status(401).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"token": res})
}
