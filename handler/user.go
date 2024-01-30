package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/iqbaleff214/easynote-backend-go/auth"
	"github.com/iqbaleff214/easynote-backend-go/helper"
	"github.com/iqbaleff214/easynote-backend-go/user"
)

type userHandler struct {
	userService user.Service
	authService auth.Service
}

func NewUserHandler(userService user.Service, authService auth.Service) *userHandler {
	return &userHandler{userService, authService}
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
		return c.Status(fiber.StatusInternalServerError).JSON(
			helper.APIResponse("Cannot create new user", "error", fiber.StatusInternalServerError, nil),
		)
	}

	token, err := h.authService.GenerateToken(newUser.ID)
	if err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(
			helper.APIResponse("Cannot generate token for new user", "error", fiber.StatusUnprocessableEntity, nil),
		)
	}

	response := helper.APIResponse("New user has been registered", "success", fiber.StatusCreated, user.FormatUser(newUser, token))

	return c.Status(fiber.StatusOK).JSON(response)
}

func (h *userHandler) Login(c *fiber.Ctx) error {

	var input user.LoginInput

	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			helper.APIResponse("There's something wrong with request body", "error", fiber.StatusBadRequest, nil),
		)
	}

	loggedUser, err := h.userService.Login(input)
	if err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(
			helper.APIResponse(err.Error(), "error", fiber.StatusUnprocessableEntity, nil),
		)
	}

	token, err := h.authService.GenerateToken(loggedUser.ID)
	if err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(
			helper.APIResponse("Cannot generate token for current user", "error", fiber.StatusUnprocessableEntity, nil),
		)
	}

	response := helper.APIResponse("Successfully logged in", "success", fiber.StatusOK, user.FormatUser(loggedUser, token))

	return c.Status(fiber.StatusOK).JSON(response)
}

func (h *userHandler) CurrentUser(c *fiber.Ctx) error {
	currentUser := c.Locals("currentUser").(user.User)

	return c.Status(fiber.StatusOK).JSON(
		helper.APIResponse("Successfully fetched current user's profile", "success", fiber.StatusOK, user.FormatUser(currentUser, "")),
	)
}

func (h *userHandler) UpdateUser(c *fiber.Ctx) error {

	var input user.UpdateUserInput

	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			helper.APIResponse("There's something wrong with request body", "error", fiber.StatusBadRequest, nil),
		)
	}

	currentUser := c.Locals("currentUser").(user.User)

	updatedUser, err := h.userService.UpdateUser(input, currentUser)
	if err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(
			helper.APIResponse("Cannot update current user's profile", "error", fiber.StatusUnprocessableEntity, nil),
		)
	}

	return c.Status(fiber.StatusOK).JSON(
		helper.APIResponse("Successfully updated current user's profile", "success", fiber.StatusOK, user.FormatUser(updatedUser, "")),
	)
}
