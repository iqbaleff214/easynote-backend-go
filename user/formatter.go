package user

type UserFormatter struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Token string `json:"token,omitempty"`
}

func FormatUser(user User, token string) UserFormatter {
	return UserFormatter{
		ID: user.ID,
		Name: user.Name,
		Email: user.Email,
		Token: token,
	}
}
