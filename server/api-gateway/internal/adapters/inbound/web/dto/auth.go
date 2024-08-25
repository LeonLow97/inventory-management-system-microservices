package dto

type LoginRequest struct {
	Username string `json:"username" validate:"required,min=5,max=50"`
	Password string `json:"password" validate:"required,min=8,max=20"`
}

type LoginResponse struct {
	FirstName string `json:"first_name,omitempty"`
	LastName  string `json:"last_name,omitempty"`
	Username  string `json:"username,omitempty"`
	Email     string `json:"email,omitempty"`
	Active    int32  `json:"active,omitempty"`
	Admin     int32  `json:"admin,omitempty"`
	Token     string `json:"token,omitempty"`
}

type SignUpRequest struct {
	FirstName string `json:"first_name" validate:"required,min=1,max=50"`
	LastName  string `json:"last_name" validate:"required,min=1,max=50"`
	Username  string `json:"username" validate:"required,min=5,max=50"`
	Password  string `json:"password" validate:"required,min=8,max=20"`
	Email     string `json:"email" validate:"required,email,min=10,max=100"`
}
