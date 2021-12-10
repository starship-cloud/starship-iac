package models

type AuthEntity struct {
	UserId string  `json:"user_id"`
	AuthToken string `json:"auth_token"`
}
