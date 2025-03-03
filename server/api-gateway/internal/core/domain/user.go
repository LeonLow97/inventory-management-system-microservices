package domain

type User struct {
	Password  string
	FirstName string
	LastName  string
	Email     string
	Token     string
	Active    bool
	Admin     bool
}
