package domain

import (
	"context"

	"github.com/Carmind-Mindia/user-hub/server/utils"
	"github.com/qiniu/qmgo/field"
)

type User struct {
	field.DefaultField `bson:",inline"`

	UserName             string   `bson:"userName"`
	Password             string   `bson:"password"`
	FirstName            string   `bson:"firstName"`
	LastName             string   `bson:"lastName"`
	Roles                []string `bson:"roles"`
	DocumentType         int      `bson:"documentType"` //1: DNI
	DocumentNumber       string   `bson:"documentNumber"`
	RecoverPasswordToken string   `bson:"recoverPasswordToken"`
	MustChangePassword   bool     `bson:"mustChangePassword"`
	FCMToken             string   `bson:"FCMToken"`
	FCMCreateTimeStamp   string   `bson:"FCMCreateTimeStamp"`
}

type DocumentType int

const (
	DNI DocumentType = iota
)

func (u *User) ValidateData() error {

	if u.UserName == "" {
		return utils.ErrOnInsertNoUsername
	}

	if u.Password == "" {
		return utils.ErrOnInsertNoPassword
	}

	if u.DocumentType < 1 || u.DocumentNumber == "" {
		return utils.ErrOnInsertNoDocument
	}

	return nil
}

type UserRepository interface {
	GetAll(ctx context.Context) ([]User, error)
	GetByUserName(ctx context.Context, userName string) (User, error)
	GetFCMTokensByUserNames(ctx context.Context, userNames []string) ([]struct {
		FCMToken string `bson:"FCMToken"`
	}, error)
	Insert(ctx context.Context, user *User) (User, error)
	Update(ctx context.Context, user *User) error
	Delete(ctx context.Context, UserName string) error
}
