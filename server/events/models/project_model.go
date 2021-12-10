package models

type ProjectEntity struct {
	ProjectId   string `json:"project_id""`
	ProjectName    string `json:"project_name"`
	Discription string `json:"project_description"`
	CreateAt    int64  `json:"create_at"`
	UpdateAt    int64  `json:"update_at"`
}
