package note

import (
	"database/sql"
	"fmt"
	"strings"
	"time"
)

type Repository interface {
	FindByID(userID, id int) (Note, error)
	FindAll(search string) ([]Note, error)
	FindByUserID(userID int, search string) ([]Note, error)
	FindByFolderID(userID, folderID int, search string) ([]Note, error)
	Save(note Note) (Note, error)
	SaveWithFolderID(note Note) (Note, error)
	Update(note Note) (Note, error)
	UpdateWithFolderID(note Note, parentID any) (Note, error)
	Delete(note Note) error
	FindTagsByNoteIDs(noteIDs []int) ([]Tag, error)
	FindTagsByName(tagNames []string) ([]Tag, error)
	SaveTags(tags []Tag) (lastID int, err error)
	SaveNoteTags(noteID int, tagIDs []int) error
}

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *repository {
	return &repository{db}
}

func (r *repository) FindByID(userID int, id int) (Note, error) {
	var note Note

	query := "SELECT n.id, n.title, n.content, n.is_public, n.user_id, u.name, COALESCE(n.folder_id, 0), COALESCE(f.name, ''), " +
		"n.created_at, n.updated_at FROM notes n LEFT JOIN folders f ON n.folder_id = f.id AND n.user_id = f.user_id JOIN users u ON n.user_id = u.id WHERE n.user_id = ? AND n.id = ?"

	err := r.db.QueryRow(query, userID, id).Scan(
		&note.ID, &note.Title, &note.Content, &note.IsPublic, &note.UserID, &note.UserName,
		&note.FolderID, &note.FolderName, &note.CreatedAt, &note.UpdatedAt,
	)
	if err != nil {
		return note, err
	}

	return note, nil
}

func (r *repository) FindAll(search string) ([]Note, error) {
	var notes []Note

	query := "SELECT n.id, n.title, n.content, n.is_public, n.user_id, u.name, COALESCE(n.folder_id, 0), COALESCE(f.name, ''), " +
		"n.created_at, n.updated_at FROM notes n LEFT JOIN folders f ON n.folder_id = f.id AND n.user_id = f.user_id JOIN users u ON n.user_id = u.id WHERE n.is_public = 1 AND n.title LIKE ?"

	rows, err := r.db.Query(query, search+"%")
	if err != nil {
		return notes, err
	}

	for rows.Next() {
		var note Note

		if err := rows.Scan(
			&note.ID, &note.Title, &note.Content, &note.IsPublic, &note.UserID, &note.UserName,
			&note.FolderID, &note.FolderName, &note.CreatedAt, &note.UpdatedAt,
		); err != nil {
			return notes, err
		}

		notes = append(notes, note)
	}

	return notes, nil
}

func (r *repository) FindByUserID(userID int, search string) ([]Note, error) {
	var notes []Note

	query := "SELECT n.id, n.title, n.content, n.is_public, n.user_id, u.name, COALESCE(n.folder_id, 0), COALESCE(f.name, ''), " +
		"n.created_at, n.updated_at FROM notes n LEFT JOIN folders f ON n.folder_id = f.id AND n.user_id = f.user_id JOIN users u ON n.user_id = u.id WHERE n.user_id = ? AND n.title LIKE ?"

	rows, err := r.db.Query(query, userID, search+"%")
	if err != nil {
		return notes, err
	}

	for rows.Next() {
		var note Note

		if err := rows.Scan(
			&note.ID, &note.Title, &note.Content, &note.IsPublic, &note.UserID, &note.UserName,
			&note.FolderID, &note.FolderName, &note.CreatedAt, &note.UpdatedAt,
		); err != nil {
			return notes, err
		}

		notes = append(notes, note)
	}

	return notes, nil
}

func (r *repository) FindByFolderID(userID int, folderID int, search string) ([]Note, error) {
	var notes []Note

	query := "SELECT n.id, n.title, n.content, n.is_public, n.user_id, u.name, COALESCE(n.folder_id, 0), COALESCE(f.name, ''), " +
		"n.created_at, n.updated_at FROM notes n LEFT JOIN folders f ON n.folder_id = f.id AND n.user_id = f.user_id JOIN users u ON n.user_id = u.id WHERE n.user_id = ? AND n.folder_id = ? AND n.title LIKE ?"

	rows, err := r.db.Query(query, userID, folderID, search+"%")
	if err != nil {
		return notes, err
	}

	for rows.Next() {
		var note Note

		if err := rows.Scan(
			&note.ID, &note.Title, &note.Content, &note.IsPublic, &note.UserID, &note.UserName,
			&note.FolderID, &note.FolderName, &note.CreatedAt, &note.UpdatedAt,
		); err != nil {
			return notes, err
		}

		notes = append(notes, note)
	}

	return notes, nil
}

func (r *repository) Save(note Note) (Note, error) {
	query := "INSERT INTO notes SET " +
		"title = ?, content = ?, is_public = ?, user_id = ?, created_at = now(), updated_at = now()"
	res, err := r.db.Exec(query, note.Title, note.Content, note.IsPublic, note.UserID)
	if err != nil {
		return note, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return note, err
	}

	note.ID = int(id)
	note.CreatedAt = time.Now()
	note.UpdatedAt = time.Now()

	return note, nil
}

func (r *repository) SaveWithFolderID(note Note) (Note, error) {
	query := "INSERT INTO notes SET " +
		"title = ?, content = ?, is_public = ?, user_id = ?, folder_id = ?, created_at = now(), updated_at = now()"

	res, err := r.db.Exec(query, note.Title, note.Content, note.IsPublic, note.UserID, note.FolderID)
	if err != nil {
		return note, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return note, err
	}

	note.ID = int(id)
	note.CreatedAt = time.Now()
	note.UpdatedAt = time.Now()

	return note, nil
}

func (r *repository) Update(note Note) (Note, error) {
	query := "UPDATE notes SET " +
		"title = ?, content = ?, is_public = ?, updated_at = NOW() " +
		"WHERE id = ?"

	note.UpdatedAt = time.Now()
	_, err := r.db.Exec(query, note.Title, note.Content, note.IsPublic, note.ID)
	if err != nil {
		return note, err
	}

	return note, nil
}

func (r *repository) UpdateWithFolderID(note Note, parentID any) (Note, error) {
	query := "UPDATE notes SET " +
		"title = ?, content = ?, is_public = ?, folder_id = ?, updated_at = NOW() " +
		"WHERE id = ?"

	note.UpdatedAt = time.Now()
	_, err := r.db.Exec(query, note.Title, note.Content, note.IsPublic, parentID, note.ID)
	if err != nil {
		return note, err
	}

	return note, nil
}

func (r *repository) Delete(note Note) error {
	query := "DELETE FROM notes WHERE id = ?"

	_, err := r.db.Exec(query, note.ID)
	return err
}

func (r *repository) FindTagsByNoteIDs(noteIDs []int) ([]Tag, error) {
	var tags []Tag

	query := "SELECT t.id, t.name, nt.note_id, t.created_at, t.updated_at " +
		"FROM tags t LEFT JOIN note_tags nt ON t.id = nt.tag_id " +
		"WHERE nt.note_id IN (%s)"

	questionMarks := []string{}
	fields := []any{}
	for _, id := range noteIDs {
		questionMarks = append(questionMarks, "?")
		fields = append(fields, id)
	}

	query = fmt.Sprintf(query, strings.Join(questionMarks, ","))
	rows, err := r.db.Query(query, fields...)
	if err != nil {
		return tags, err
	}

	for rows.Next() {
		var tag Tag

		if err := rows.Scan(
			&tag.ID, &tag.Name, &tag.NoteID, &tag.CreatedAt, &tag.UpdatedAt,
		); err != nil {
			return tags, err
		}

		tags = append(tags, tag)
	}

	return tags, nil
}

func (r *repository) FindTagsByName(tagNames []string) ([]Tag, error) {
	var tags []Tag

	query := "SELECT t.id, t.name, t.created_at, t.updated_at " +
		"FROM tags t WHERE t.name IN (%s)"

	questionMarks := []string{}
	fields := []any{}
	for _, tagName := range tagNames {
		questionMarks = append(questionMarks, "?")
		fields = append(fields, tagName)
	}

	query = fmt.Sprintf(query, strings.Join(questionMarks, ","))
	rows, err := r.db.Query(query, fields...)
	if err != nil {
		return tags, err
	}

	for rows.Next() {
		var tag Tag

		if err := rows.Scan(
			&tag.ID, &tag.Name, &tag.CreatedAt, &tag.UpdatedAt,
		); err != nil {
			return tags, err
		}

		tags = append(tags, tag)
	}

	return tags, nil
}

func (r *repository) SaveTags(tags []Tag) (lastID int, err error) {
	query := "INSERT INTO tags (name, created_at, updated_at) VALUES "

	questionMarks := []string{}
	fields := []any{}

	for _, tag := range tags {
		questionMarks = append(questionMarks, "(?, now(), now())")
		fields = append(fields, tag.Name)
	}

	query += strings.Join(questionMarks, ",")

	res, err := r.db.Exec(query, fields...)
	if err != nil {
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	lastID = int(id)

	return
}

func (r *repository) SaveNoteTags(noteID int, tagIDs []int) error {
	query := "INSERT INTO note_tags (note_id, tag_id, created_at, updated_at) VALUES "

	questionMarks := []string{}
	fields := []any{}

	for _, tag := range tagIDs {
		questionMarks = append(questionMarks, fmt.Sprintf("(%d, ?, now(), now())", noteID))
		fields = append(fields, tag)
	}

	query += strings.Join(questionMarks, ",")

	_, err := r.db.Exec(query, fields...)
	if err != nil {
		return err
	}

	return nil
}
