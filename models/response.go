package models

type Response struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
}

type TokenResponse struct {
	Response
	Role  string `json:"role,omitempty"`
	Token string `json:"token"`
}

type GetUserResponse struct {
	Response
	Role string `json:"role,omitempty"`
	User User   `json:"user,omitempty"`
}

type GetAdminResponse struct {
	Response
	Role  string `json:"role,omitempty"`
	Admin Admin  `json:"user,omitempty"`
}

type GetUserListResponse struct {
	Response
	UserList []User `json:"user_list,omitempty"`
}
