package domain

type LoginInput struct {
	Username string
	Password string
}

type SignUpInput struct {
	FirstName string
	LastName  string
	Username  string
	Password  string
	Email     string
}
