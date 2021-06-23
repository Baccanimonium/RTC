package Repos

// TODO разобраться с SQL Запросом и получать сущности раздельно
type Participant struct {
	Id             int     `json:"id" db:"id"`
	IdUser         int     `json:"id_user" db:"id_user" binding:"required"`
	Salary         float32 `json:"salary"`
	Qualifications string  `json:"qualifications"`
	Contacts       string  `json:"contacts"`
	Name           string  `json:"name" binding:"required"`
	About          string  `json:"about"`
	Address        string  `json:"address"`
	Phone          string  `json:"phone"`
	//UserCreate `json:"user"`
	//Doctor `json:"doctor"`
}
