package models

type LoginUser struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type DeleteUser struct {
	Username string `json:"username"`
}

type UserViewFields struct {
	Admin    bool   `json:"admin"`
	Username string `json:"username"`
}
