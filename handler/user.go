package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v3/log"
	"github.com/iqbaleff214/easynote-backend-go/helper"
	"github.com/iqbaleff214/easynote-backend-go/user"
)

type userHandler struct {
	userService user.Service
}

func NewUserHandler(userService user.Service) *userHandler {
	return &userHandler{userService}
}

func (h *userHandler) RegisterUser(c *fiber.Ctx) error {

	var input user.RegisterUserInput

	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			helper.APIResponse("There's something wrong with request body", "error", fiber.StatusBadRequest, nil),
		)
	}

	newUser, err := h.userService.RegisterUser(input)
	if err != nil {
		log.Info(err)
		return c.Status(fiber.StatusInternalServerError).JSON(
			helper.APIResponse("Cannot create new user", "error", fiber.StatusInternalServerError, nil),
		)
	}

	token := "dummy-token"
	response := helper.APIResponse("New user has been registered", "success", fiber.StatusCreated, user.FormatUser(newUser, token))

	return c.Status(fiber.StatusOK).JSON(response)
}
