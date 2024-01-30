package folder

import "errors"

type Service interface {
	FindFolders(userID int, folderID int) ([]Folder, error)
	CreateFolder(input CreateFolderInput, userID int) (Folder, error)
	UpdateFolder(input UpdateFolderInput, userID, folderID int) (Folder, error)
	DeleteFolder(userID, folderID int) error
}

type service struct {
	repository Repository
}

func NewService(repository Repository) *service {
	return &service{repository}
}

func (s *service) FindFolders(userID int, folderID int) ([]Folder, error) {
	var folders []Folder

	if userID == 0 {
		return folders, errors.New("no user available on this session")
	}

	if folderID != 0 {
		folders, err := s.repository.FindByParentID(userID, folderID)
		if err != nil {
			return folders, err
		}

		return folders, nil
	}

	folders, err := s.repository.FindByUserID(userID)
	if err != nil {
		return folders, err
	}

	return folders, nil
}

func (s *service) CreateFolder(input CreateFolderInput, userID int) (Folder, error) {
	var folder Folder

	folder.Name = input.Name
	folder.ParentID = input.ParentID
	folder.UserID = userID

	if folder.ParentID == 0 {
		newFolder, err := s.repository.Save(folder)
		if err != nil {
			return newFolder, err
		}

		return newFolder, nil
	}

	newFolder, err := s.repository.SaveWithParentID(folder)
	if err != nil {
		return newFolder, err
	}

	return newFolder, nil
}

func (s *service) UpdateFolder(input UpdateFolderInput, userID, folderID int) (Folder, error) {
	currentFolder, err := s.repository.FindByID(userID, folderID)
	if err != nil {
		return currentFolder, err
	}

	currentFolder.Name = input.Name
	currentFolder.ParentID = input.ParentID

	if currentFolder.ParentID == 0 {
		newFolder, err := s.repository.Update(currentFolder)
		if err != nil {
			return currentFolder, err
		}

		return newFolder, nil
	}

	var parentID any
	if currentFolder.ParentID > 0 {
		parentID = currentFolder.ParentID
	} else {
		parentID = nil
		currentFolder.ParentID = 0
		currentFolder.ParentName = ""
	}

	newFolder, err := s.repository.UpdateWithParentID(currentFolder, parentID)
	if err != nil {
		return currentFolder, err
	}

	return newFolder, nil
}

func (s *service) DeleteFolder(userID, folderID int) error {
	currentFolder, err := s.repository.FindByID(userID, folderID)
	if err != nil {
		return err
	}

	return s.repository.Delete(currentFolder)
}
