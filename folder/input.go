package folder

type CreateFolderInput struct {
	Name     string `json:"name"`
	ParentID int    `json:"parent_folder_id"`
}

type UpdateFolderInput struct {
	Name     string `json:"name"`
	ParentID int    `json:"parent_folder_id"`
}
