package dto

type RestLoginUser struct {
	User  RestUser `json:"user"`
	Token string   `json:"token"`
}
