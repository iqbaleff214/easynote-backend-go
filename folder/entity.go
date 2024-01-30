package folder

import "time"

type Folder struct {
	ID         int
	Name       string
	ParentName string
	ParentID   int
	UserID     int
	CreatedAt  time.Time
	UpdatedAt  time.Time
}
