package models

type UserEntity struct {
	Userid string `json:userid`
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}
