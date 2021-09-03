package domain

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
)

type User struct {
	Id             interface{}
	UserName       string `bson:"userName"`
	Password       string `bson:"password"`
	FirstName      string `bson:"firstName"`
	LastName       string `bson:"lastName"`
	DocumentType   int    `bson:"documentType"` //1: DNI
	DocumentNumber string `bson:"documentNumber"`
}

type DocumentType int

const (
	DNI DocumentType = iota
)

func (u *User) ValidateData() (bool, error) {

	if u.UserName == "" {
		return false, echo.NewHTTPError(http.StatusBadRequest, "userName empty")
	}

	if u.Password == "" {
		return false, echo.NewHTTPError(http.StatusBadRequest, "password empty")
	}

	if u.DocumentType < 1 || u.DocumentNumber == "" {
		return false, echo.NewHTTPError(http.StatusBadRequest, "document type empty")
	}

	return true, nil
}

type UserRepository interface {
	GetAll(ctx context.Context) ([]User, error)
	//GetByUserName() (User, error)
	Insert(ctx context.Context, user *User) (User, error)
	//Update(*User) (User, error)
	//Delete(UserName string)
}

type UserUseCase interface {
	GetAll(ctx context.Context) ([]User, error)
	//GetByUserName() (User, error)
	Insert(ctx context.Context, user *User) (User, error)
	//Update(*User) (User, error)
	//Delete(UserName string)
}
