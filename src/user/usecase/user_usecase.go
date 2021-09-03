package usecase

import (
	"context"

	"github.com/Fonzeka/Jame/src/domain"
)

type UserUseCase struct {
	repo domain.UserRepository
}

func NewUserUseCase(repo *domain.UserRepository) domain.UserUseCase {
	return &UserUseCase{repo: *repo}
}

func (uc *UserUseCase) GetAll(ctx context.Context) ([]domain.User, error) {
	return uc.repo.GetAll(ctx)
}

func (uc *UserUseCase) Insert(ctx context.Context, user *domain.User) (domain.User, error) {
	if _, err := user.ValidateData(); err != nil {
		return *user, err
	}
	return uc.repo.Insert(ctx, user)
}
