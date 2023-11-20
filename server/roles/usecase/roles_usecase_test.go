package usecase_test

import (
	"context"
	"testing"

	"github.com/Carmind-Mindia/user-hub/server/domain"
	"github.com/Carmind-Mindia/user-hub/server/domain/mocks"
	"github.com/Carmind-Mindia/user-hub/server/roles/usecase"
	"github.com/Carmind-Mindia/user-hub/server/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestInsertRole(t *testing.T) {
	repoMock := new(mocks.RolesRepository)

	u := usecase.NewRolesUseCase(repoMock)

	insertRoleCases := []struct {
		nameTest      string
		roleToInsert  domain.Role
		prepare       func(repoRole *mock.Mock)
		errorExpected error
	}{
		{
			nameTest: "success",
			roleToInsert: domain.Role{
				Name: "insertTest",
			},
			prepare: func(repoRole *mock.Mock) {
				repoRole.On("Insert", mock.Anything, mock.Anything).Return(nil).Once()
			},
			errorExpected: nil,
		},
		{
			nameTest: "empty-name",
			roleToInsert: domain.Role{
				Name: "",
			},
			prepare: func(repoRole *mock.Mock) {
			},
			errorExpected: utils.ErrBadRequest,
		},
	}

	for _, x := range insertRoleCases {
		t.Run(x.nameTest, func(t *testing.T) {
			x.prepare(&repoMock.Mock)

			err := u.InsertRole(context.TODO(), x.roleToInsert)

			if x.errorExpected != nil {
				assert.Error(t, err)
				assert.Equal(t, x.errorExpected, err)
			} else {
				assert.NoError(t, err)
			}

			repoMock.AssertExpectations(t)
		})
	}
}

func TestValidateRoles(t *testing.T) {
	repoMock := new(mocks.RolesRepository)

	u := usecase.NewRolesUseCase(repoMock)

	rolesInDb := []domain.Role{
		{
			Name: "admin",
		},
		{
			Name: "test",
		},
		{
			Name: "lorem",
		},
		{
			Name: "ipsum",
		},
	}

	validateRolesCases := []struct {
		nameTest        string
		rolesToValidate []string
		prepare         func(repoRole *mock.Mock)
		errorExpected   error
	}{
		{
			nameTest:        "success",
			rolesToValidate: []string{"admin", "test"},
			prepare: func(repoRole *mock.Mock) {
				repoRole.On("GetAll", mock.Anything).Return(rolesInDb, nil).Once()
			},
			errorExpected: nil,
		},
		{
			nameTest:        "bad-db",
			rolesToValidate: []string{"admin", "test"},
			prepare: func(repoRole *mock.Mock) {
				repoRole.On("GetAll", mock.Anything).Return(nil, utils.ErrUnauthorized).Once()
			},
			errorExpected: utils.ErrUnauthorized,
		},
		{
			nameTest:        "bad-validation",
			rolesToValidate: []string{"admin", "test", "not-exist", "other"},
			prepare: func(repoRole *mock.Mock) {
				repoRole.On("GetAll", mock.Anything).Return(rolesInDb, nil).Once()
			},
			errorExpected: utils.ErrNoValidRole("not-exist"),
		},
	}

	for _, x := range validateRolesCases {
		t.Run(x.nameTest, func(t *testing.T) {
			x.prepare(&repoMock.Mock)

			err := u.ValidateRoles(context.TODO(), x.rolesToValidate...)

			if x.errorExpected != nil {
				assert.Error(t, err)
				assert.Equal(t, err, x.errorExpected)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestDeleteRoles(t *testing.T) {
	repoMock := new(mocks.RolesRepository)

	u := usecase.NewRolesUseCase(repoMock)

	deleteRolesCases := []struct {
		nameTest      string
		roleToDelete  string
		prepare       func(repoRole *mock.Mock)
		errorExpected error
	}{
		{
			nameTest:     "success",
			roleToDelete: "test",
			prepare: func(repoRole *mock.Mock) {
				repoRole.On("Delete", mock.Anything, "test").Return(nil).Once()
			},
			errorExpected: nil,
		},
		{
			nameTest:     "bad-validation",
			roleToDelete: "",
			prepare: func(repoRole *mock.Mock) {
			},
			errorExpected: utils.ErrBadRequest,
		},
	}

	for _, x := range deleteRolesCases {
		t.Run(x.nameTest, func(t *testing.T) {
			x.prepare(&repoMock.Mock)

			err := u.DeleteRole(context.TODO(), x.roleToDelete)

			if x.errorExpected != nil {
				assert.Error(t, err)
				assert.Equal(t, err, x.errorExpected)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
