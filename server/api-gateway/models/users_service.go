package models

type UpdateUserRequest struct {
	FirstName string `json:"first_name" validate:"omitempty,min=1,max=50"`
	LastName  string `json:"last_name" validate:"omitempty,min=1,max=50"`
	Username  string `json:"username" validate:"required,min=5,max=50"`
	Password  string `json:"password" validate:"omitempty,min=8,max=20"`
	Email     string `json:"email" validate:"omitempty,email,min=10,max=100"`
}

type User struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	Active    int    `json:"active"`
	Admin     int    `json:"admin"`
}
