package models

type UserEntity struct {
	Userid       string `json:"userid""`
	Username     string `json:"username"`
	Password     string `json:"password"`
	Email        string `json:"email"`
	Externaluser string `json:"externaluser"`
	UserLocal    bool   `json:"userlocal"`
	Salt         string `json:"salt"`
	CreateAt     int    `json:"create_at"`
	UpdateAt     int    `json:"update_at"`
}
