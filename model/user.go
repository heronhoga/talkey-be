package model

type User struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserLogin struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type UserRegister struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}