package models

type UserEntity struct {
	UserId string `json:userid`
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}
