package note

import "time"

type Note struct {
	ID         int
	Title      string
	Content    string
	IsPublic   bool
	UserID     int
	UserName   string
	FolderID   int
	FolderName string
	Tags       []Tag
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

type Tag struct {
	ID        int
	Name      string
	NoteID    int
	CreatedAt time.Time
	UpdatedAt time.Time
}
