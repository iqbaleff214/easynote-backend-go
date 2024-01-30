package user

import "database/sql"

type Repository interface {
	Save(user User) (User, error)
}

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *repository {
	return &repository{db}
}

func (r *repository) Save(user User) (User, error) {
	query := "INSERT INTO users SET "+
	"name = ?, email = ?, password = ?, created_at = NOW(), updated_at = NOW()"

	res, err := r.db.Exec(query, user.Name, user.Email, user.Password)
	if err != nil {
		return user, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return user, err
	}

	user.ID = int(id)


	return user, nil
}