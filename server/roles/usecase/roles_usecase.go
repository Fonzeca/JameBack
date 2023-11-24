package usecase

import (
	"context"

	"github.com/Carmind-Mindia/user-hub/server/domain"
	"github.com/Carmind-Mindia/user-hub/server/utils"
	"github.com/thoas/go-funk"
)

type RolesUseCase struct {
	repo domain.RolesRepository
}

func NewRolesUseCase(repo domain.RolesRepository) RolesUseCase {
	uc := RolesUseCase{repo: repo}

	res, err := uc.GetAllRoles(context.Background())
	if err == nil && len(res) == 0 {

		for _, r := range res {
			if r.Name == "admin" {
				return uc
			}
		}

		rol := domain.Role{
			Name: "admin",
		}
		uc.InsertRole(context.Background(), rol)
	}

	return uc
}

func (uc *RolesUseCase) InsertRole(ctx context.Context, rol domain.Role) error {

	//TODO: hacer un error para este caso
	if len(rol.Name) <= 0 {
		return utils.ErrBadRequest
	}

	return uc.repo.Insert(ctx, &rol)
}

func (uc *RolesUseCase) DeleteRole(ctx context.Context, name string) error {
	if len(name) <= 0 {
		return utils.ErrBadRequest
	}

	return uc.repo.Delete(ctx, name)
}

func (uc *RolesUseCase) GetAllRoles(ctx context.Context) (res []domain.Role, err error) {
	return uc.repo.GetAll(ctx)
}

// TODO: comentar esta funcion
func (uc *RolesUseCase) ValidateRoles(ctx context.Context, roles ...string) error {

	data, err := uc.repo.GetAll(ctx)

	if err != nil {
		return err
	}

	var rolesDb []string
	funk.ForEach(data, func(x domain.Role) {
		rolesDb = append(rolesDb, x.Name)
	})

	intersection := funk.IntersectString(rolesDb, roles)
	if len(intersection) != len(roles) {
		_, r := funk.DifferenceString(intersection, roles)

		return utils.ErrNoValidRole(r[0])
	}

	return nil
}
