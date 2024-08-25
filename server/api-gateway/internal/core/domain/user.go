package domain

type User struct {
	Username  string
	Password  string
	FirstName string
	LastName  string
	Email     string
	Token     string
	Active    int
	Admin     int
}
