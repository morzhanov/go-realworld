package dto

type CreateUserDto struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type GetUserDto struct {
	ID       string `json:"id"`
	Username string `json:"username"`
}
