package usecase

import (
	"context"
	"math/rand"
	"strconv"
	"time"

	"github.com/Fonzeca/UserHub/src/domain"
	"github.com/Fonzeca/UserHub/src/emails"
	"github.com/Fonzeca/UserHub/src/roles/usecase"
	myjwt "github.com/Fonzeca/UserHub/src/security/jwt"
	"github.com/Fonzeca/UserHub/src/user/delivery/modelview"
	"github.com/Fonzeca/UserHub/src/utils"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

type UserUseCase struct {
	repo     domain.UserRepository
	rolecase usecase.RolesUseCase
}

func NewUserUseCase(repo domain.UserRepository, roleUsecase usecase.RolesUseCase) UserUseCase {
	uc := UserUseCase{repo: repo, rolecase: roleUsecase}

	go validateAdminUser(&uc)

	return uc
}

func validateAdminUser(uc *UserUseCase) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	data, _ := uc.GetAll(ctx)
	if len(data) <= 0 {
		data2, _ := uc.rolecase.GetAllRoles(ctx)
		if len(data2) <= 0 {
			uc.rolecase.InsertRole(ctx, domain.Role{
				Name: "admin",
			})
		}

		uc.Insert(ctx, &domain.User{
			UserName:       "afonzo",
			Password:       "123456",
			FirstName:      "Alexis",
			LastName:       "Fonzo",
			DocumentType:   1,
			DocumentNumber: "38096937",
			Roles:          []string{"admin"},
		})
	}
}

// Obtiene el usuario por el nombre de usuario
func (uc *UserUseCase) GetByUserName(ctx context.Context, userName string) (domain.User, error) {
	return uc.repo.GetByUserName(ctx, userName)
}

// Actualiza los datos del usuario segun vengan en el request
func (uc *UserUseCase) Update(ctx context.Context, user *domain.User) error {

	// Busca si existe con el username
	usrDb, err := uc.repo.GetByUserName(ctx, user.UserName)

	// Si no lo encuentra, devuelve error
	if err != nil {
		return err
	}

	//Validamos campos vacios en el request

	if len(user.FirstName) > 0 {
		usrDb.FirstName = user.FirstName
	}

	if len(user.LastName) > 0 {
		usrDb.LastName = user.LastName
	}

	if user.DocumentType > 0 {
		usrDb.DocumentType = user.DocumentType
	}

	if len(user.DocumentNumber) > 0 {
		usrDb.DocumentNumber = user.DocumentNumber
	}

	//Validamos que el usuario tenga bien los roles
	if err := uc.rolecase.ValidateRoles(ctx, user.Roles...); err != nil {
		return err
	}

	usrDb.Roles = user.Roles

	return uc.repo.Update(ctx, &usrDb)
}

// Obtiene todos los usuario
// TODO: paginacion
func (uc *UserUseCase) GetAll(ctx context.Context) ([]domain.User, error) {
	return uc.repo.GetAll(ctx)
}

// Creamos un usuario
func (uc *UserUseCase) Insert(ctx context.Context, user *domain.User) (domain.User, error) {
	//Validamos los datos del usuario para insertar
	if err := user.ValidateData(); err != nil {
		return *user, err
	}

	//Validamos que el usuario tenga bien los roles
	if err := uc.rolecase.ValidateRoles(ctx, user.Roles...); err != nil {
		return *user, err
	}

	//Hasheamos la pass
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), 8)

	if err != nil {
		//throw Generic Error Message
		return *user, utils.ErrInternalError
	}

	//Seteamos la pass al user
	user.Password = string(hashedPassword)

	//Insertamos el user
	return uc.repo.Insert(ctx, user)
}

// Borramos un user
func (ux *UserUseCase) Delete(ctx context.Context, UserName string) error {
	return ux.repo.Delete(ctx, UserName)
}

// Metodo de login
func (ux *UserUseCase) Login(ctx context.Context, userName string, password string) (string, error) {
	//Buscamos un usuario con el mismo userName
	user, err := ux.repo.GetByUserName(ctx, userName)

	if err != nil {
		//Si no lo encuentra
		return "", utils.ErrTryLogin
	}

	//Comparamos la pass
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", utils.ErrTryLogin
	}

	// Generamos el token
	token, err := myjwt.GenerateToken(&user)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (ux *UserUseCase) GetUserByToken(ctx context.Context, claims jwt.MapClaims) (domain.User, error) {
	//Obtenemos el userName desde el mismo contexto de echo
	userName := claims["userName"].(string)

	//Buscamos un usuario con el mismo userName
	user, err := ux.repo.GetByUserName(ctx, userName)
	user.Password = ""
	if err != nil {
		return user, err
	}

	return user, nil
}

func (ux *UserUseCase) SendEmailRecoverPassword(ctx context.Context, username string) error {

	user, err := ux.GetByUserName(ctx, username)
	if err != nil {
		return err
	}

	u4 := 1000 + rand.Intn(8999)

	//Hasheamos la pass
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(strconv.Itoa(u4)), 8)
	if err != nil {
		return err
	}

	user.RecoverPasswordToken = string(hashedPassword)

	emails.ChannelEmails <- emails.Recuperacion{
		Email: username,
		Token: strconv.Itoa(u4),
	}

	err = ux.repo.Update(ctx, &user)
	if err != nil {
		return err
	}

	return nil
}

func (ux *UserUseCase) ValidateRecoverPasswordToken(ctx context.Context, view modelview.ResetPassword) (domain.User, error) {
	user, err := ux.GetByUserName(ctx, view.Email)
	if err != nil {
		return user, err
	}

	//Comparamos los tokens
	if err := bcrypt.CompareHashAndPassword([]byte(user.RecoverPasswordToken), []byte(view.Token)); err != nil {
		return user, utils.ErrOnChangePassword
	}

	return user, nil
}

func (ux *UserUseCase) ResetPasswordWithToken(ctx context.Context, view modelview.ResetPassword) error {
	user, err := ux.ValidateRecoverPasswordToken(ctx, view)
	if err != nil {
		return err
	}

	hashed, _ := bcrypt.GenerateFromPassword([]byte(view.NewPassword), 8)
	user.Password = string(hashed)
	user.RecoverPasswordToken = ""

	err = ux.repo.Update(ctx, &user)
	if err != nil {
		return err
	}
	return nil
}
