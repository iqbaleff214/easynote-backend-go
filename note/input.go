package note

type CreateNoteInput struct {
	Title    string   `json:"title"`
	Content  string   `json:"content"`
	IsPublic bool     `json:"is_public"`
	FolderID int      `json:"folder_id"`
	Tags     []string `json:"tags"`
}

type UpdateNoteInput struct {
	Title    string `json:"title"`
	Content  string `json:"content"`
	IsPublic bool   `json:"is_public"`
	FolderID int    `json:"folder_id"`
}
