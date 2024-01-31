package note

import (
	"errors"
)

type Service interface {
	PublicNotes(search string) ([]Note, error)
	FindNotes(userID int, folderID int, search string) ([]Note, error)
	FindNote(userID int, noteID int) (Note, error)
	CreateNote(input CreateNoteInput, userID int) (Note, error)
	UpdateNote(input UpdateNoteInput, userID, noteID int) (Note, error)
	DeleteNote(userID int, noteID int) error
}

type service struct {
	repository Repository
}

func NewService(repository Repository) *service {
	return &service{repository}
}

func (s *service) PublicNotes(search string) ([]Note, error) {
	var notes []Note

	notes, err := s.repository.FindAll(search)
	if err != nil {
		return notes, err
	}

	if len(notes) > 0 {
		var noteIDs []int

		for _, note := range notes {
			noteIDs = append(noteIDs, note.ID)
		}

		tags, err := s.repository.FindTagsByNoteIDs(noteIDs)
		if err != nil {
			return notes, err
		}

		mappedTags := map[int][]Tag{}
		for _, tag := range tags {
			mappedTags[tag.NoteID] = append(mappedTags[tag.NoteID], tag)
		}

		for i, note := range notes {
			if existingTag, ok := mappedTags[note.ID]; ok {
				notes[i].Tags = existingTag
			}
		}
	}

	return notes, nil
}

func (s *service) FindNotes(userID int, folderID int, search string) ([]Note, error) {
	var notes []Note

	if userID == 0 {
		return notes, errors.New("no user available on this session")
	}

	var err error

	if folderID != 0 {
		notes, err = s.repository.FindByFolderID(userID, folderID, search)
		if err != nil {
			return notes, err
		}
	} else {
		notes, err = s.repository.FindByUserID(userID, search)
		if err != nil {
			return notes, err
		}
	}

	if len(notes) > 0 {
		var noteIDs []int

		for _, note := range notes {
			noteIDs = append(noteIDs, note.ID)
		}

		tags, err := s.repository.FindTagsByNoteIDs(noteIDs)
		if err != nil {
			return notes, err
		}

		mappedTags := map[int][]Tag{}
		for _, tag := range tags {
			mappedTags[tag.NoteID] = append(mappedTags[tag.NoteID], tag)
		}

		for i, note := range notes {
			if existingTag, ok := mappedTags[note.ID]; ok {
				notes[i].Tags = existingTag
			}
		}
	}

	return notes, nil
}

func (s *service) FindNote(userID int, noteID int) (Note, error) {
	note, err := s.repository.FindByID(userID, noteID)
	if err != nil {
		return note, err
	}

	note.Tags, err = s.repository.FindTagsByNoteIDs([]int{noteID})
	if err != nil {
		return note, err
	}

	return note, nil
}

func (s *service) CreateNote(input CreateNoteInput, userID int) (Note, error) {
	var note Note

	note.Title = input.Title
	note.Content = input.Content
	note.IsPublic = input.IsPublic
	note.FolderID = input.FolderID
	note.UserID = userID

	var err error
	if note.FolderID == 0 {
		note, err = s.repository.Save(note)
		if err != nil {
			return note, err
		}
	} else {
		note, err = s.repository.SaveWithFolderID(note)
		if err != nil {
			return note, err
		}
	}

	if len(input.Tags) == 0 {
		return note, nil
	}

	mappingInputTags := map[string]bool{}
	for _, tag := range input.Tags {
		mappingInputTags[tag] = true
	}

	existingTags, err := s.repository.FindTagsByName(input.Tags)
	if err != nil {
		return note, err
	}
	note.Tags = append(note.Tags, existingTags...)

	var tagsIDs []int
	for _, tag := range existingTags {
		delete(mappingInputTags, tag.Name)
		tagsIDs = append(tagsIDs, tag.ID)
	}

	var newTags []Tag
	for tag := range mappingInputTags {
		newTags = append(newTags, Tag{Name: tag})
	}
	note.Tags = append(note.Tags, newTags...)

	if len(newTags) > 0 {
		lastID, err := s.repository.SaveTags(newTags)
		if err != nil {
			return note, err
		}

		for i := 0; i < len(newTags); i++ {
			tagsIDs = append(tagsIDs, lastID+i)
		}
	}

	if err := s.repository.SaveNoteTags(note.ID, tagsIDs); err != nil {
		return note, err
	}

	return note, nil
}

func (s *service) UpdateNote(input UpdateNoteInput, userID, noteID int) (Note, error) {
	oldNote, err := s.repository.FindByID(userID, noteID)
	if err != nil {
		return oldNote, err
	}

	oldNote.Title = input.Title
	oldNote.Content = input.Content
	oldNote.IsPublic = input.IsPublic
	oldNote.FolderID = input.FolderID

	oldNote.Tags, err = s.repository.FindTagsByNoteIDs([]int{noteID})
	if err != nil {
		return oldNote, err
	}

	if oldNote.FolderID == 0 {
		newNote, err := s.repository.Update(oldNote)
		if err != nil {
			return oldNote, err
		}

		return newNote, nil
	}

	var folderID any
	if oldNote.FolderID > 0 {
		folderID = oldNote.FolderID
	} else {
		folderID = nil
		oldNote.FolderID = 0
		oldNote.FolderName = ""
	}

	newNote, err := s.repository.UpdateWithFolderID(oldNote, folderID)
	if err != nil {
		return oldNote, err
	}

	return newNote, nil
}

func (s *service) DeleteNote(userID int, noteID int) error {
	note, err := s.repository.FindByID(userID, noteID)
	if err != nil {
		return err
	}

	return s.repository.Delete(note)
}
