package domain

type User struct {
	ID        int64
	Password  string
	FirstName string
	LastName  string
	Email     string
	Token     string
	Active    bool
	Admin     bool
}
