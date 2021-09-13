package usecase_test

import (
	"context"
	"testing"

	"github.com/Fonzeka/Jame/src/domain"
	"github.com/Fonzeka/Jame/src/domain/mocks"
	"github.com/Fonzeka/Jame/src/user/usecase"
	"github.com/Fonzeka/Jame/src/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateUser(t *testing.T) {
	mockUserRepo := new(mocks.UserRepositoryMock)

	mockUser := domain.User{
		UserName:       "username",
		Password:       "Password",
		DocumentType:   1,
		DocumentNumber: "38096937",
	}

	t.Run("success", func(t *testing.T) {
		mockUserRepo.On("Insert", mock.Anything, mock.Anything).Return(domain.User{}, nil).Once()

		u := usecase.NewUserUseCase(mockUserRepo)

		_, err := u.Insert(context.TODO(), &mockUser)

		assert.NoError(t, err)
		mockUserRepo.AssertExpectations(t)

	})

	t.Run("fail-validate-properties", func(t *testing.T) {
		userTest := mockUser
		userTest.Password = ""

		u := usecase.NewUserUseCase(mockUserRepo)

		_, err := u.Insert(context.TODO(), &userTest)
		assert.Error(t, err)

		mockUserRepo.AssertExpectations(t)
	})

}

func TestLogin(t *testing.T) {
	mockUserRepo := new(mocks.UserRepositoryMock)

	mockUser := domain.User{
		UserName:       "username",
		Password:       "Password",
		DocumentType:   1,
		DocumentNumber: "38096937",
		Roles:          []string{"asdasd", "asdasd", "admin", "asdasd", "asdasd"},
	}

	t.Run("insert-and-login", func(t *testing.T) {

		mockUserRepo.On("Insert", mock.Anything, mock.Anything).Return(domain.User{}, nil).Once()
		mockUserRepo.On("GetByUserName", mock.Anything, mock.AnythingOfType("string")).Return(domain.User{}, nil).Once()

		u := usecase.NewUserUseCase(mockUserRepo)

		saltPassword := mockUser.Password

		_, err := u.Insert(context.TODO(), &mockUser)
		assert.NoError(t, err)

		token, err := u.Login(context.TODO(), mockUser.UserName, saltPassword)
		assert.NoError(t, err)
		assert.NotNil(t, token)
		assert.NotEmpty(t, token)
		t.Log(token)

		mockUser.Password = saltPassword

		mockUserRepo.AssertExpectations(t)
	})

	t.Run("invalid-login", func(t *testing.T) {
		mockUserRepo.On("GetByUserName", mock.Anything, mock.AnythingOfType("string")).Return(domain.User{}, utils.NewHTTPError(500, "asd")).Once()
		mockUserRepo.On("GetByUserName", mock.Anything, mock.AnythingOfType("string")).Return(domain.User{}, nil).Once()

		u := usecase.NewUserUseCase(mockUserRepo)

		token, err := u.Login(context.TODO(), "wrongUserName", mockUser.Password)
		assert.Error(t, err)
		assert.Empty(t, token)

		token, err = u.Login(context.TODO(), mockUser.UserName, "asdasd")
		assert.Error(t, err)
		assert.Empty(t, token)

		mockUserRepo.AssertExpectations(t)
	})
}

func TestUpdateUser(t *testing.T) {

	mockUserRepo := new(mocks.UserRepositoryMock)

	mockUser := domain.User{
		UserName:       "username",
		Password:       "Password",
		DocumentType:   1,
		DocumentNumber: "38096937",
		Roles:          []string{"asdasd", "asdasd", "admin", "asdasd", "asdasd"},
	}
	mockUserRepo.On("Insert", mock.Anything, mock.Anything).Return(domain.User{}, nil).Once()
	mockUserRepo.Insert(context.TODO(), &mockUser)

	t.Run("update", func(t *testing.T) {
		mockUserRepo.On("GetByUserName", mock.Anything, mock.AnythingOfType("string")).Return(domain.User{}, nil).Once()
		mockUserRepo.On("Update", mock.Anything, mock.Anything).Return(nil).Once()
		mockUserRepo.On("GetByUserName", mock.Anything, mock.AnythingOfType("string")).Return(domain.User{}, nil).Once()

		u := usecase.NewUserUseCase(mockUserRepo)

		userchange := mockUser

		userchange.DocumentNumber = "45269148"

		err := u.Update(context.TODO(), &userchange)
		assert.NoError(t, err)

		usr, err := u.GetByUserName(context.TODO(), userchange.UserName)
		assert.NoError(t, err)

		assert.Equal(t, usr.UserName, userchange.UserName)
		assert.NotEqual(t, usr.DocumentNumber, mockUser.DocumentNumber)

		mockUserRepo.AssertExpectations(t)
	})

}
