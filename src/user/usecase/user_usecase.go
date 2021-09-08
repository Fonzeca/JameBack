package usecase

import (
	"context"
	"net/http"

	"github.com/Fonzeka/Jame/src/domain"
	"github.com/Fonzeka/Jame/src/utils"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

type jwtCustomClaims struct {
	UserName string `json:"userName"`
	Admin    bool   `json:"admin"`
	jwt.StandardClaims
}

type UserUseCase struct {
	repo domain.UserRepository
}

func NewUserUseCase(repo domain.UserRepository) domain.UserUseCase {
	return &UserUseCase{repo: repo}
}

func (uc *UserUseCase) GetAll(ctx context.Context) ([]domain.User, error) {
	return uc.repo.GetAll(ctx)
}

func (uc *UserUseCase) Insert(ctx context.Context, user *domain.User) (domain.User, error) {
	//Validamos los datos del usuario para insertar
	if err := user.ValidateData(); err != nil {
		return *user, err
	}

	//Hasheamos la pass
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), 8)

	if err != nil {
		//throw Generic Error Message
		return *user, utils.NewHTTPError(http.StatusInternalServerError, "GEM1")
	}

	//Seteamos la pass al user
	user.Password = string(hashedPassword)

	//Insertamos el user
	return uc.repo.Insert(ctx, user)
}

func (ux *UserUseCase) Delete(ctx context.Context, UserName string) error {
	return ux.repo.Delete(ctx, UserName)
}

func (ux *UserUseCase) Login(ctx context.Context, userName string, password string) (string, error) {
	//Buscamos un usuario con el mismo userName
	user, err := ux.repo.GetByUserName(ctx, userName)

	if err != nil {
		//Si no lo encuentra
		return "", utils.NewHTTPError(http.StatusUnauthorized, "UNA1")
	}

	//Comparamos la pass
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", utils.NewHTTPError(http.StatusUnauthorized, "UNA1")
	}

	token, err := generateToken(&user)
	if err != nil {
		return "", err
	}

	return token, nil
}
