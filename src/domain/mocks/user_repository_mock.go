package mocks

import (
	"context"

	"github.com/Fonzeka/Jame/src/domain"
	"github.com/stretchr/testify/mock"
)

type UserRepositoryMock struct {
	users []domain.User
	mock.Mock
}

func (u *UserRepositoryMock) GetAll(ctx context.Context) ([]domain.User, error) {
	ret := u.Called(ctx)

	return ret.Get(0).([]domain.User), ret.Error(1)
}

func (u *UserRepositoryMock) GetByUserName(ctx context.Context, userName string) (domain.User, error) {
	ret := u.Called(ctx)

	for _, u2 := range u.users {
		if u2.UserName == userName {
			return u2, ret.Error(1)
		}
	}

	return ret.Get(0).(domain.User), ret.Error(1)
}

func (u *UserRepositoryMock) Insert(ctx context.Context, user *domain.User) (domain.User, error) {
	ret := u.Called(ctx)

	u.users = append(u.users, *user)

	return ret.Get(0).(domain.User), ret.Error(1)
}

func (u *UserRepositoryMock) Delete(ctx context.Context, UserName string) error {
	ret := u.Called(ctx)

	return ret.Error(0)
}
