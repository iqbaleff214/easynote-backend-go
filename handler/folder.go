package handler

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/iqbaleff214/easynote-backend-go/folder"
	"github.com/iqbaleff214/easynote-backend-go/helper"
	"github.com/iqbaleff214/easynote-backend-go/user"
)

type folderHandler struct {
	folderService folder.Service
}

func NewFolderHandler(folderService folder.Service) *folderHandler {
	return &folderHandler{folderService}
}

func (h *folderHandler) FindFolders(c *fiber.Ctx) error {
	currentUser := c.Locals("currentUser").(user.User)
	folderID, _ := strconv.Atoi(c.Query("parent_id"))

	folders, err := h.folderService.FindFolders(currentUser.ID, folderID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			helper.APIResponse("Cannot fetch folders", "error", fiber.StatusBadRequest, nil),
		)
	}

	return c.Status(fiber.StatusOK).JSON(
		helper.APIResponse("Successfully fetched folder's list", "success", fiber.StatusOK, folder.FormatFolders(folders)),
	)
}

func (h *folderHandler) CreateFolder(c *fiber.Ctx) error {

	var input folder.CreateFolderInput

	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			helper.APIResponse("There's something wrong with request body", "error", fiber.StatusBadRequest, nil),
		)
	}

	currentUser := c.Locals("currentUser").(user.User)

	newFolder, err := h.folderService.CreateFolder(input, currentUser.ID)
	if err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(
			helper.APIResponse("Cannot create new folder", "error", fiber.StatusUnprocessableEntity, nil),
		)
	}

	return c.Status(fiber.StatusOK).JSON(
		helper.APIResponse("Successfully created new folder", "success", fiber.StatusOK, folder.FormatFolder(newFolder)),
	)
}

func (h *folderHandler) UpdateFolder(c *fiber.Ctx) error {

	var input folder.UpdateFolderInput

	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			helper.APIResponse("There's something wrong with request body", "error", fiber.StatusBadRequest, nil),
		)
	}

	folderID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			helper.APIResponse("There's something wrong with your folder id", "error", fiber.StatusBadRequest, nil),
		)
	}

	currentUser := c.Locals("currentUser").(user.User)

	updatedFolder, err := h.folderService.UpdateFolder(input, currentUser.ID, folderID)
	if err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(
			helper.APIResponse("Cannot update the folder", "error", fiber.StatusUnprocessableEntity, nil),
		)
	}

	return c.Status(fiber.StatusOK).JSON(
		helper.APIResponse("Successfully updated the folder", "success", fiber.StatusOK, folder.FormatFolder(updatedFolder)),
	)
}

func (h *folderHandler) DeleteFolder(c *fiber.Ctx) error {
	folderID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			helper.APIResponse("There's something wrong with your folder id", "error", fiber.StatusBadRequest, nil),
		)
	}

	currentUser := c.Locals("currentUser").(user.User)

	if err := h.folderService.DeleteFolder(currentUser.ID, folderID); err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(
			helper.APIResponse("Cannot delete the folder", "error", fiber.StatusUnprocessableEntity, nil),
		)
	}

	return c.Status(fiber.StatusOK).JSON(
		helper.APIResponse("Successfully deleted the folder", "success", fiber.StatusOK, nil),
	)
}