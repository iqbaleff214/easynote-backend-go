package folder

type FolderFormatter struct {
	ID             int    `json:"id"`
	Name           string `json:"name"`
	ParentFolder   string `json:"parent_folder,omitempty"`
	ParentFolderID int    `json:"parent_folder_id,omitempty"`
}

func FormatFolder(folder Folder) FolderFormatter {
	return FolderFormatter{
		ID:             folder.ID,
		Name:           folder.Name,
		ParentFolderID: folder.ParentID,
		ParentFolder:   folder.ParentName,
	}
}

func FormatFolders(folders []Folder) []FolderFormatter {
	folderFormatters := []FolderFormatter{}

	for _, folder := range folders {
		folderFormatter := FormatFolder(folder)
		folderFormatters = append(folderFormatters, folderFormatter)
	}

	return folderFormatters
}
