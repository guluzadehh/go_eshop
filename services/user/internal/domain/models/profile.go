package models

type Profile struct {
	*User
	FirstName string
	LastName  string
	Phone     string
	Picture   string
}
