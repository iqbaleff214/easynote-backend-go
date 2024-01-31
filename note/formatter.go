package note

import "time"

type NoteFormatter struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	IsPublic  bool      `json:"is_public"`
	Folder    string    `json:"folder,omitempty"`
	FolderID  int       `json:"folder_id,omitempty"`
	Tags      []string  `json:"tags"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type NotePublicFormatter struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	Author    string    `json:"author"`
	Tags      []string  `json:"tags"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type TagFormatter struct {
	ID   int    `json:"id"`
	Name string `json:"tag"`
}

func FormatNote(note Note) NoteFormatter {
	tags := []string{}

	for _, tag := range note.Tags {
		tags = append(tags, tag.Name)
	}

	return NoteFormatter{
		ID:        note.ID,
		Title:     note.Title,
		Content:   note.Content,
		IsPublic:  note.IsPublic,
		Folder:    note.FolderName,
		FolderID:  note.FolderID,
		Tags:      tags,
		CreatedAt: note.CreatedAt,
		UpdatedAt: note.UpdatedAt,
	}
}

func FormatNotes(notes []Note) []NoteFormatter {
	noteFormatters := []NoteFormatter{}

	for _, note := range notes {
		noteFormatter := FormatNote(note)
		noteFormatters = append(noteFormatters, noteFormatter)
	}

	return noteFormatters
}

func FormatPublicNote(note Note) NotePublicFormatter {
	tags := []string{}

	for _, tag := range note.Tags {
		tags = append(tags, tag.Name)
	}

	return NotePublicFormatter{
		ID:        note.ID,
		Title:     note.Title,
		Content:   note.Content,
		Author:    note.UserName,
		Tags:      tags,
		CreatedAt: note.CreatedAt,
		UpdatedAt: note.UpdatedAt,
	}
}

func FormatPublicNotes(notes []Note) []NotePublicFormatter {
	noteFormatters := []NotePublicFormatter{}

	for _, note := range notes {
		noteFormatter := FormatPublicNote(note)
		noteFormatters = append(noteFormatters, noteFormatter)
	}

	return noteFormatters
}

func FormatTag(tag Tag) TagFormatter {
	return TagFormatter{
		ID:   tag.ID,
		Name: tag.Name,
	}
}

func FormatTags(tags []Tag) []TagFormatter {
	tagFormatters := []TagFormatter{}

	for _, tag := range tags {
		tagFormatter := FormatTag(tag)
		tagFormatters = append(tagFormatters, tagFormatter)
	}

	return tagFormatters
}
