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
	Roles                []string `bson:"roles"`
	RecoverPasswordToken string   `bson:"recoverPasswordToken"`
	MustChangePassword   bool     `bson:"mustChangePassword"`
}

func (u *User) ValidateData() error {

	if u.UserName == "" {
		return utils.ErrOnInsertNoUsername
	}

	if u.Password == "" {
		return utils.ErrOnInsertNoPassword
	}

	return nil
}

type UserRepository interface {
	GetAll(ctx context.Context) ([]User, error)
	GetByUserName(ctx context.Context, userName string) (User, error)
	Insert(ctx context.Context, user *User) (User, error)
	Update(ctx context.Context, user *User) error
	Delete(ctx context.Context, UserName string) error
}
