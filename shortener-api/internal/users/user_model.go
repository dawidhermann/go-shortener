package users

type UserViewModel struct {
	UserId   int    `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

type UserCreateViewModel struct {
	Username        string `json:"username"`
	Email           string `json:"email"`
	Password        string `json:"password"`
	PasswordConfirm string `json:"passwordConfirm"`
}

type UserPatchModel struct {
	Email           *string `json:"email"`
	Password        *string `json:"password"`
	PasswordConfirm *string `json:"passwordConfirm"`
}

type user struct {
	UserId   int
	Username string
	Password string
	Email    string
}

func NewUserViewModel(user user) UserViewModel {
	return UserViewModel{
		UserId:   user.UserId,
		Username: user.Username,
		Email:    user.Email,
	}
}
