package usecase_test

import (
	"context"
	"testing"

	"github.com/Fonzeca/UserHub/src/domain"
	"github.com/Fonzeca/UserHub/src/domain/mocks"
	_rolesUseCase "github.com/Fonzeca/UserHub/src/roles/usecase"
	_userUseCase "github.com/Fonzeca/UserHub/src/user/usecase"
	"github.com/Fonzeca/UserHub/src/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
)

func TestCreateUser(t *testing.T) {
	mockUserRepo := new(mocks.UserRepository)
	mockRoleRepo := new(mocks.RolesRepository)
	ru := _rolesUseCase.NewRolesUseCase(mockRoleRepo)

	mockUserRepo.On("GetAll", mock.Anything).Return([]domain.User{
		domain.User{},
	}, nil)
	u := _userUseCase.NewUserUseCase(mockUserRepo, ru)

	var createUserTests = []struct {
		nameTest      string
		userToInsert  domain.User
		prepare       func(repoUser *mock.Mock, repoRole *mock.Mock)
		errorExpected error
	}{
		{
			nameTest: "success",
			userToInsert: domain.User{
				UserName:       "username",
				Password:       "Password",
				DocumentType:   1,
				DocumentNumber: "38096937",
				Roles:          []string{"admin"},
			},
			errorExpected: nil,
			prepare: func(repoUser *mock.Mock, repoRole *mock.Mock) {
				roles := []domain.Role{
					{Name: "admin"},
					{Name: "test"},
				}

				repoRole.On("GetAll", mock.Anything).Return(roles, nil).Once()
				repoUser.On("Insert", mock.Anything, mock.Anything).Return(domain.User{}, nil).Once()
			},
		},
		{
			nameTest: "fail-by-property-password",
			userToInsert: domain.User{
				UserName:       "username",
				Password:       "",
				DocumentType:   1,
				DocumentNumber: "38096937",
			},
			errorExpected: utils.ErrOnInsertNoPassword,
			prepare: func(repo *mock.Mock, repoRole *mock.Mock) {
			},
		},
	}

	for _, x := range createUserTests {
		t.Run(x.nameTest, func(t *testing.T) {
			x.prepare(&mockUserRepo.Mock, &mockRoleRepo.Mock)

			_, err := u.Insert(context.TODO(), &x.userToInsert)

			if x.errorExpected != nil {
				assert.Error(t, err)
				assert.Equal(t, err, x.errorExpected)
			} else {
				assert.NoError(t, err)
			}

			mockUserRepo.AssertExpectations(t)
		})
	}

}

func TestLogin(t *testing.T) {
	mockUserRepo := new(mocks.UserRepository)
	mockRoleRepo := new(mocks.RolesRepository)
	ru := _rolesUseCase.NewRolesUseCase(mockRoleRepo)

	mockUserRepo.On("GetAll", mock.Anything).Return([]domain.User{
		domain.User{},
	}, nil)
	u := _userUseCase.NewUserUseCase(mockUserRepo, ru)

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("12345"), 8)
	userInDb := domain.User{
		UserName:       "pepito",
		Password:       string(hashedPassword),
		DocumentType:   1,
		DocumentNumber: "38096937",
		Roles:          []string{"admin"},
	}

	var loginUserTests = []struct {
		nameTest      string
		userName      string
		password      string
		prepare       func(repoUser *mock.Mock, repoRole *mock.Mock)
		errorExpected error
	}{
		{
			nameTest: "success",
			userName: "pepito",
			password: "12345",
			prepare: func(repoUser, repoRole *mock.Mock) {
				mockUserRepo.On("GetByUserName", mock.Anything, userInDb.UserName).Return(userInDb, nil).Once()
			},
			errorExpected: nil,
		},
		{
			nameTest: "invalid-login-pass",
			userName: "pepito",
			password: "badPass",
			prepare: func(repoUser, repoRole *mock.Mock) {
				mockUserRepo.On("GetByUserName", mock.Anything, userInDb.UserName).Return(userInDb, nil).Once()
			},
			errorExpected: utils.ErrTryLogin,
		},
		{
			nameTest: "invalid-login-username",
			userName: "badName",
			password: "12345",
			prepare: func(repoUser, repoRole *mock.Mock) {
				mockUserRepo.On("GetByUserName", mock.Anything, "badName").Return(domain.User{}, utils.ErrUserNotFound).Once()
			},
			errorExpected: utils.ErrTryLogin,
		},
	}

	for _, x := range loginUserTests {
		t.Run(x.nameTest, func(t *testing.T) {

			x.prepare(&mockUserRepo.Mock, &mockRoleRepo.Mock)

			token, err := u.Login(context.TODO(), x.userName, x.password)

			if x.errorExpected != nil {
				assert.Error(t, err)
				assert.Equal(t, x.errorExpected, err)
				assert.Empty(t, token)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, token)
				assert.NotEmpty(t, token)
			}

			mockUserRepo.AssertExpectations(t)
		})
	}
}

func TestUpdateUser(t *testing.T) {

	mockUserRepo := new(mocks.UserRepository)
	mockRoleRepo := new(mocks.RolesRepository)
	ru := _rolesUseCase.NewRolesUseCase(mockRoleRepo)

	mockUserRepo.On("GetAll", mock.Anything).Return([]domain.User{
		domain.User{},
	}, nil)
	u := _userUseCase.NewUserUseCase(mockUserRepo, ru)

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("12345"), 8)
	userInDb := domain.User{
		UserName:       "pepito",
		Password:       string(hashedPassword),
		DocumentType:   1,
		DocumentNumber: "38096937",
		Roles:          []string{"admin"},
	}

	var updateUserTests = []struct {
		nameTest      string
		updateUser    domain.User
		prepare       func(repoUser *mock.Mock, repoRole *mock.Mock)
		errorExpected error
	}{
		{
			nameTest: "success",
			updateUser: domain.User{
				UserName:       "pepito",
				Password:       "asdasd",
				DocumentType:   11,
				DocumentNumber: "20-38096937-9",
				Roles:          []string{"admin", "test"},
			},
			prepare: func(repoUser, repoRole *mock.Mock) {
				roles := []domain.Role{
					{Name: "admin"},
					{Name: "test"},
				}

				repoUser.On("GetByUserName", mock.Anything, userInDb.UserName).Return(userInDb, nil).Once()
				repoRole.On("GetAll", mock.Anything).Return(roles, nil).Once()
				repoUser.On("Update", mock.Anything, mock.Anything).Return(nil).Once()
			},
			errorExpected: nil,
		},
	}

	for _, x := range updateUserTests {
		t.Run(x.nameTest, func(t *testing.T) {
			x.prepare(&mockUserRepo.Mock, &mockRoleRepo.Mock)
			err := u.Update(context.TODO(), &x.updateUser)

			if x.errorExpected != nil {
				assert.Error(t, err)
				assert.Equal(t, err, x.errorExpected)
			} else {
				assert.NoError(t, err)
				//TODO: validar cambios de roles
			}

		})
	}

}

func TestSendEmail(t *testing.T) {
	/* ch := make(chan string)

	go utils.SendEmailChannel(ch)

	t.Run("TestSendEmail", func(t *testing.T) {
		t.Log("Esperando 5 segundos")
		time.Sleep(5 * time.Second)

		email1 := "alexisfonzos@gmail.com"
		t.Log("Enviando a " + email1)
		ch <- email1

		t.Log("Esperando 5 segundos")
		time.Sleep(10 * time.Second)

		email2 := "alexisfonzos_@hotmail.com"
		t.Log("enviando a " + email2)
		ch <- email2

	}) */

}
