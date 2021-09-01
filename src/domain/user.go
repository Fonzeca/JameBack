package domain

type User struct {
	UserName       string
	Password       string
	FirstName      string
	LastName       string
	DocumentType   int //1: DNI
	DocumentNumber string
}
