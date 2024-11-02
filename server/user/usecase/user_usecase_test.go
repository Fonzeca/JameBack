package usecase_test

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Carmind-Mindia/user-hub/server/domain"
	"github.com/Carmind-Mindia/user-hub/server/domain/mocks"
	_rolesUseCase "github.com/Carmind-Mindia/user-hub/server/roles/usecase"
	_userUseCase "github.com/Carmind-Mindia/user-hub/server/user/usecase"
	"github.com/Carmind-Mindia/user-hub/server/utils"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
)

func TestCreateUser(t *testing.T) {
	mockUserRepo := new(mocks.UserRepository)
	mockRoleRepo := new(mocks.RolesRepository)
	mockRoleRepo.On("GetAll", mock.Anything).Return(nil, fmt.Errorf("error")).Once()
	ru := _rolesUseCase.NewRolesUseCase(mockRoleRepo)

	mockUserRepo.On("GetByUserName", mock.Anything, mock.Anything).Return(domain.User{}, nil).Once()
	u := _userUseCase.NewUserUseCase(mockUserRepo, ru, nil)

	var createUserTests = []struct {
		nameTest      string
		userToInsert  domain.User
		prepare       func(repoUser *mock.Mock, repoRole *mock.Mock)
		errorExpected error
	}{
		{
			nameTest: "success",
			userToInsert: domain.User{
				UserName: "username",
				Password: "Password",
				Roles:    []string{"admin"},
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
				UserName: "username",
				Password: "",
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
	mockRoleRepo.On("GetAll", mock.Anything).Return([]domain.Role{
		{Name: "admin"},
		{Name: "test"},
	}, nil)
	mockUserRepo.On("GetByUserName", mock.Anything, mock.Anything).Return(domain.User{}, nil).Once()

	ru := _rolesUseCase.NewRolesUseCase(mockRoleRepo)

	// mockUserRepo.On("GetAll", mock.Anything).Return([]domain.User{
	// 	domain.User{},
	// }, nil)
	u := _userUseCase.NewUserUseCase(mockUserRepo, ru, nil)

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("12345"), 8)
	userInDb := domain.User{
		UserName: "pepito",
		Password: string(hashedPassword),
		Roles:    []string{"admin"},
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

			// Setup
			e := echo.New()
			req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(""))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			result, err := u.Login(context.TODO(), x.userName, x.password, c)

			if x.errorExpected != nil {
				assert.Error(t, err)
				assert.Equal(t, x.errorExpected, err)
				assert.Empty(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.NotEmpty(t, result)

				// cookies := rec.Result().Cookies()
				// assert.NotEmpty(t, cookies)
				// assert.Equal(t, cookies[0].Name, "session")
			}

			mockUserRepo.AssertExpectations(t)
		})
	}
}

func TestUpdateUser(t *testing.T) {

	mockUserRepo := new(mocks.UserRepository)
	mockRoleRepo := new(mocks.RolesRepository)

	mockRoleRepo.On("GetAll", mock.Anything).Return([]domain.Role{
		{Name: "admin"},
		{Name: "test"},
	}, nil)
	ru := _rolesUseCase.NewRolesUseCase(mockRoleRepo)

	mockUserRepo.On("GetByUserName", mock.Anything, mock.Anything).Return(domain.User{}, nil)
	u := _userUseCase.NewUserUseCase(mockUserRepo, ru, nil)

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("12345"), 8)
	userInDb := domain.User{
		UserName: "pepito",
		Password: string(hashedPassword),
		Roles:    []string{"admin"},
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
				UserName: "pepito",
				Password: "asdasd",
				Roles:    []string{"admin", "test"},
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
