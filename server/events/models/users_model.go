package models

type UserEntity struct {
	Userid       string `json:"userid""`
	Username     string `json:"username"`
	Password     string `json:"_"`
	Email        string `json:"email"`
	Externaluser string `json:"externaluser"`
	UserLocal    bool   `json:"userlocal"`
	Salt         string `json:"_"`
	CreateAt     int64  `json:"create_at"`
	UpdateAt     int64  `json:"update_at"`
}

type AuthEntity struct {
	Userid string  `json:"userid"`
	AuthToken string `json:"auth_token"`
}
