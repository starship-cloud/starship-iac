package models

type UserEntity struct {
	UserId       string `json:"user_id""`
	UserName     string `json:"user_name"`
	Password     string `json:"_"`
	Email        string `json:"email"`
	ExternalUser string `json:"external_user"`
	UserLocal    bool   `json:"user_local"`
	Salt         string `json:"_"`
	CreateAt     int64  `json:"create_at"`
	UpdateAt     int64  `json:"update_at"`
}
