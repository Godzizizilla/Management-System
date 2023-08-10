package models

type Response struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
}

type TokenResponse struct {
	Response
	Token string `json:"token"`
}

type GetUserResponse struct {
	Response
	User User `json:"user,omitempty"`
}

type GetUserListResponse struct {
	Response
	UserList []User `json:"user_list,omitempty"`
}
