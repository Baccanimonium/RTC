package Repos

type Participant struct {
	UserCreate `json:"user"`
	Doctor     `json:"doctor"`
	Patient    `json:"patient"`
}
