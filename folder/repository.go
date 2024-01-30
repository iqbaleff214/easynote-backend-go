package folder

import (
	"database/sql"
	"time"
)

type Repository interface {
	FindByID(userID, id int) (Folder, error)
	FindByUserID(userID int) ([]Folder, error)
	FindByParentID(userID, parentID int) ([]Folder, error)
	Save(folder Folder) (Folder, error)
	SaveWithParentID(folder Folder) (Folder, error)
	Update(folder Folder) (Folder, error)
	UpdateWithParentID(folder Folder, parentID any) (Folder, error)
	Delete(folder Folder) error
}

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *repository {
	return &repository{db}
}

func (r *repository) FindByID(userID, id int) (Folder, error) {
	var folder Folder

	query := "SELECT f.id, f.name, COALESCE(f.parent_id, 0), f.user_id, f.created_at, f.updated_at, COALESCE(p.name, '') " +
		"FROM folders f LEFT JOIN folders p ON f.parent_id = p.id WHERE f.user_id = ? AND f.id = ?"

	err := r.db.QueryRow(query, userID, id).Scan(
		&folder.ID, &folder.Name, &folder.ParentID,
		&folder.UserID, &folder.CreatedAt, &folder.UpdatedAt, &folder.ParentName,
	)
	if err != nil {
		return folder, err
	}

	return folder, nil
}

func (r *repository) FindByUserID(userID int) ([]Folder, error) {

	var folders []Folder

	query := "SELECT f.id, f.name, COALESCE(f.parent_id, 0), f.user_id, f.created_at, f.updated_at, COALESCE(p.name, '') " +
		"FROM folders f LEFT JOIN folders p ON f.parent_id = p.id WHERE f.user_id = ?"

	rows, err := r.db.Query(query, userID)
	if err != nil {
		return folders, err
	}

	for rows.Next() {
		var folder Folder
		if err := rows.Scan(
			&folder.ID, &folder.Name, &folder.ParentID,
			&folder.UserID, &folder.CreatedAt, &folder.UpdatedAt, &folder.ParentName,
		); err != nil {
			return folders, err
		}

		folders = append(folders, folder)
	}

	return folders, nil
}

func (r *repository) FindByParentID(userID, parentID int) ([]Folder, error) {

	var folders []Folder

	query := "SELECT f.id, f.name, COALESCE(f.parent_id, 0), f.user_id, f.created_at, f.updated_at, COALESCE(p.name, '') " +
		"FROM folders f LEFT JOIN folders p ON f.parent_id = p.id WHERE f.parent_id = ? AND f.user_id = ?"

	rows, err := r.db.Query(query, parentID, userID)
	if err != nil {
		return folders, err
	}

	for rows.Next() {
		var folder Folder
		if err := rows.Scan(
			&folder.ID, &folder.Name, &folder.ParentID,
			&folder.UserID, &folder.CreatedAt, &folder.UpdatedAt, &folder.ParentName,
		); err != nil {
			return folders, err
		}

		folders = append(folders, folder)
	}

	return folders, nil
}

func (r *repository) Save(folder Folder) (Folder, error) {
	query := "INSERT INTO folders SET " +
		"name = ?, user_id = ?, created_at = NOW(), updated_at = NOW()"

	res, err := r.db.Exec(query, folder.Name, folder.UserID)
	if err != nil {
		return folder, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return folder, err
	}

	folder.ID = int(id)
	folder.CreatedAt = time.Now()
	folder.UpdatedAt = time.Now()

	return folder, nil
}

func (r *repository) SaveWithParentID(folder Folder) (Folder, error) {
	query := "INSERT INTO folders SET " +
		"name = ?, user_id = ?, parent_id = ?, created_at = NOW(), updated_at = NOW()"

	res, err := r.db.Exec(query, folder.Name, folder.UserID, folder.ParentID)
	if err != nil {
		return folder, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return folder, err
	}

	folder.ID = int(id)
	folder.CreatedAt = time.Now()
	folder.UpdatedAt = time.Now()

	return folder, nil
}

func (r *repository) Update(folder Folder) (Folder, error) {
	query := "UPDATE folders SET " +
		"name = ?, updated_at = NOW() " +
		"WHERE id = ?"

	folder.UpdatedAt = time.Now()
	_, err := r.db.Exec(query, folder.Name, folder.ID)
	if err != nil {
		return folder, err
	}

	return folder, nil
}

func (r *repository) UpdateWithParentID(folder Folder, parentID any) (Folder, error) {
	query := "UPDATE folders SET " +
		"name = ?, parent_id = ?, updated_at = NOW() " +
		"WHERE id = ?"

	folder.UpdatedAt = time.Now()
	_, err := r.db.Exec(query, folder.Name, parentID, folder.ID)
	if err != nil {
		return folder, err
	}

	return folder, nil
}

func (r *repository) Delete(folder Folder) error {
	query := "DELETE FROM folders WHERE id = ?"

	_, err := r.db.Exec(query, folder.ID)
	return err
}
