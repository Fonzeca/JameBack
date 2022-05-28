package domain

import (
	"context"

	"github.com/qiniu/qmgo/field"
)

type Role struct {
	field.DefaultField `bson:",inline"`

	Name string `bson:"name"`
}

type RolesRepository interface {
	GetAll(ctx context.Context) ([]Role, error)
	Insert(ctx context.Context, role *Role) error
	Delete(ctx context.Context, name string) error
}
