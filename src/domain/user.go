package domain

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/qiniu/qmgo/field"
)

type User struct {
	field.DefaultField `bson:",inline"`

	UserName       string   `bson:"userName"`
	Password       string   `bson:"password"`
	FirstName      string   `bson:"firstName"`
	LastName       string   `bson:"lastName"`
	Roles          []string `bson:"roles"`
	DocumentType   int      `bson:"documentType"` //1: DNI
	DocumentNumber string   `bson:"documentNumber"`
}

type DocumentType int

const (
	DNI DocumentType = iota
)

func (u *User) ValidateData() error {

	if u.UserName == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "userName empty")
	}

	if u.Password == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "password empty")
	}

	if u.DocumentType < 1 || u.DocumentNumber == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "document type empty")
	}

	return nil
}

type UserRepository interface {
	GetAll(ctx context.Context) ([]User, error)
	GetByUserName(ctx context.Context, userName string) (User, error)
	Insert(ctx context.Context, user *User) (User, error)
	//Update(*User) (User, error)
	Delete(ctx context.Context, UserName string) error
}

type UserUseCase interface {
	Login(ctx context.Context, userName string, password string) (string, error)
	GetAll(ctx context.Context) ([]User, error)
	//GetByUserName() (User, error)
	Insert(ctx context.Context, user *User) (User, error)
	//Update(*User) (User, error)
	Delete(ctx context.Context, UserName string) error
}
