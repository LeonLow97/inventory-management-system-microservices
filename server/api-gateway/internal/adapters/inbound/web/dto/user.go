package dto

type User struct {
	Username  string `json:"username,omitempty"`
	FirstName string `json:"first_name,omitempty"`
	LastName  string `json:"last_name,omitempty"`
	Email     string `json:"email,omitempty"`
	Token     string `json:"token,omitempty"`
	Active    int    `json:"active,omitempty"`
	Admin     int    `json:"admin,omitempty"`
}

type GetUsersResponse struct {
	Users []User `json:"users"`
}

type UpdateUserRequest struct {
	FirstName string `json:"first_name" validate:"omitempty,min=1,max=50"`
	LastName  string `json:"last_name" validate:"omitempty,min=1,max=50"`
	Username  string `json:"username" validate:"required,min=5,max=50"`
	Password  string `json:"password" validate:"omitempty,min=8,max=20"`
	Email     string `json:"email" validate:"omitempty,email,min=10,max=100"`
}
