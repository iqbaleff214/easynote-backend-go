package handler

import (
	"log"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/iqbaleff214/easynote-backend-go/helper"
	"github.com/iqbaleff214/easynote-backend-go/note"
	"github.com/iqbaleff214/easynote-backend-go/user"
)

type noteHandler struct {
	noteService note.Service
}

func NewNoteHandler(noteService note.Service) *noteHandler {
	return &noteHandler{noteService}
}

func (h *noteHandler) FindPublicNotes(c *fiber.Ctx) error {
	search := c.Query("q")

	notes, err := h.noteService.PublicNotes(search)
	if err != nil {
		log.Println(err)
		return c.Status(fiber.StatusBadRequest).JSON(
			helper.APIResponse("Cannot fetch notes", "error", fiber.StatusBadRequest, nil),
		)
	}

	return c.Status(fiber.StatusOK).JSON(
		helper.APIResponse("Successfully fetched public notes", "success", fiber.StatusOK, note.FormatPublicNotes(notes)),
	)
}

func (h *noteHandler) FindNotes(c *fiber.Ctx) error {
	search := c.Query("q")
	folderID, _ := strconv.Atoi(c.Query("folder_id"))

	currentUser := c.Locals("currentUser").(user.User)

	notes, err := h.noteService.FindNotes(currentUser.ID, folderID, search)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			helper.APIResponse("Cannot fetch notes", "error", fiber.StatusBadRequest, nil),
		)
	}

	return c.Status(fiber.StatusOK).JSON(
		helper.APIResponse("Successfully fetched notes", "success", fiber.StatusOK, note.FormatNotes(notes)),
	)
}

func (h *noteHandler) FindNote(c *fiber.Ctx) error {
	noteID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			helper.APIResponse("There's something wrong with your note id", "error", fiber.StatusBadRequest, nil),
		)
	}

	currentUser := c.Locals("currentUser").(user.User)

	fetchedNote, err := h.noteService.FindNote(currentUser.ID, noteID)
	if err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(
			helper.APIResponse("Cannot fetch the note", "error", fiber.StatusUnprocessableEntity, nil),
		)
	}

	return c.Status(fiber.StatusOK).JSON(
		helper.APIResponse("Successfully fetched the note", "success", fiber.StatusOK, note.FormatNote(fetchedNote)),
	)
}

func (h *noteHandler) CreateNote(c *fiber.Ctx) error {

	var input note.CreateNoteInput

	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			helper.APIResponse("There's something wrong with request body", "error", fiber.StatusBadRequest, nil),
		)
	}

	currentUser := c.Locals("currentUser").(user.User)

	newNote, err := h.noteService.CreateNote(input, currentUser.ID)
	if err != nil {
		log.Println(err)
		return c.Status(fiber.StatusUnprocessableEntity).JSON(
			helper.APIResponse("Cannot create new note", "error", fiber.StatusUnprocessableEntity, nil),
		)
	}

	return c.Status(fiber.StatusOK).JSON(
		helper.APIResponse("Successfully created new note", "success", fiber.StatusOK, note.FormatNote(newNote)),
	)
}

func (h *noteHandler) UpdateNote(c *fiber.Ctx) error {

	var input note.UpdateNoteInput

	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			helper.APIResponse("There's something wrong with request body", "error", fiber.StatusBadRequest, nil),
		)
	}

	noteID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			helper.APIResponse("There's something wrong with your note id", "error", fiber.StatusBadRequest, nil),
		)
	}

	currentUser := c.Locals("currentUser").(user.User)

	updatedNote, err := h.noteService.UpdateNote(input, currentUser.ID, noteID)
	if err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(
			helper.APIResponse("Cannot update the note", "error", fiber.StatusUnprocessableEntity, nil),
		)
	}

	return c.Status(fiber.StatusOK).JSON(
		helper.APIResponse("Successfully updated the note", "success", fiber.StatusOK, note.FormatNote(updatedNote)),
	)
}

func (h *noteHandler) DeleteNote(c *fiber.Ctx) error {
	noteID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			helper.APIResponse("There's something wrong with your note id", "error", fiber.StatusBadRequest, nil),
		)
	}

	currentUser := c.Locals("currentUser").(user.User)

	if err := h.noteService.DeleteNote(currentUser.ID, noteID); err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(
			helper.APIResponse("Cannot delete the note", "error", fiber.StatusUnprocessableEntity, nil),
		)
	}

	return c.Status(fiber.StatusOK).JSON(
		helper.APIResponse("Successfully deleted the note", "success", fiber.StatusOK, nil),
	)
}
