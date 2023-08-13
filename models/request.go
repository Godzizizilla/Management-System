package models

type LoginRequest struct {
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
}

type RegisterRequest struct {
	Username  string `json:"user_name,omitempty"`
	Password  string `json:"password,omitempty"`
	Grade     string `json:"grade,omitempty"`
	StudentID uint   `json:"student_id,omitempty"`
}

type UpdateInfoRequest struct {
	NewName     string `json:"new_name,omitempty"`
	Email       string `json:"email,omitempty"`
	Phone       string `json:"phone,omitempty"`
	OldPassword string `json:"old_password,omitempty"`
	NewPassword string `json:"new_password,omitempty"`
}

type DeleteRequest struct {
	Password string `json:"password,omitempty"`
}
