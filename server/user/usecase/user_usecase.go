package usecase

import (
	"context"
	"math/rand"
	"strconv"

	"github.com/Carmind-Mindia/user-hub/server/domain"
	"github.com/Carmind-Mindia/user-hub/server/roles/usecase"
	"github.com/Carmind-Mindia/user-hub/server/user/delivery/modelview"
	"github.com/Carmind-Mindia/user-hub/server/utils"
	"github.com/labstack/echo/v4"
	"github.com/spf13/viper"
	"golang.org/x/crypto/bcrypt"
)

type UserUseCase struct {
	repo     domain.UserRepository
	rolecase usecase.RolesUseCase
}

func NewUserUseCase(repo domain.UserRepository, roleUsecase usecase.RolesUseCase) UserUseCase {
	uc := UserUseCase{repo: repo, rolecase: roleUsecase}

	superAdminUsername := viper.GetString("superAdmin.username")
	superAdminPassword := viper.GetString("superAdmin.password")

	_, err := uc.GetByUserName(context.Background(), superAdminUsername)
	if err != nil {
		//Si no existe el usuario, lo creamos
		usr := domain.User{
			UserName:           superAdminUsername,
			Password:           superAdminPassword,
			MustChangePassword: viper.GetBool("mustChangeFirstTimePassword"),
			Roles:              []string{"admin"},
		}

		uc.Insert(context.Background(), &usr)
	}

	return uc
}

// Obtiene el usuario por el nombre de usuario
func (uc *UserUseCase) GetByUserName(ctx context.Context, userName string) (domain.User, error) {
	user, err := uc.repo.GetByUserName(ctx, userName)
	if err != nil {
		return user, err
	}

	user.Password = ""

	return user, nil
}

// Actualiza los datos del usuario segun vengan en el request
func (uc *UserUseCase) Update(ctx context.Context, user *domain.User) error {

	// Busca si existe con el username
	usrDb, err := uc.repo.GetByUserName(ctx, user.UserName)

	// Si no lo encuentra, devuelve error
	if err != nil {
		return err
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
	users, err := uc.repo.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	// remove password to all users
	for i := range users {
		users[i].Password = ""
	}

	return users, nil
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

	user.MustChangePassword = viper.GetBool("mustChangeFirstTimePassword")

	finalUser, err := uc.repo.Insert(ctx, user)
	if err != nil {
		user.Password = ""
		return *user, err
	}

	finalUser.Password = ""
	return finalUser, nil
}

// Borramos un user
func (ux *UserUseCase) Delete(ctx context.Context, UserName string) error {
	return ux.repo.Delete(ctx, UserName)
}

// Metodo de login
func (ux *UserUseCase) Login(ctx context.Context, userName string, password string, echoCtx echo.Context) (*modelview.LoginResponse, error) {
	//Buscamos un usuario con el mismo userName
	user, err := ux.repo.GetByUserName(ctx, userName)

	if err != nil {
		//Si no lo encuentra
		return nil, utils.ErrTryLogin
	}

	//Comparamos la pass
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, utils.ErrTryLogin
	}

	changePassword := user.MustChangePassword

	isAdmin := false
	for _, v := range user.Roles {
		if v == "admin" {
			isAdmin = true
		}
	}

	response := &modelview.LoginResponse{
		MustChangePassword: changePassword,
		Username:           user.UserName,
		Admin:              isAdmin,
		Roles:              user.Roles,
	}

	return response, nil
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

	// err = ux.fastEmailClient.SendRecoverPassword(username, strconv.Itoa(u4))
	// if err != nil {
	// 	return err
	// }

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

	user.Password = ""

	return user, nil
}

func (ux *UserUseCase) ResetPasswordWithToken(ctx context.Context, view modelview.ResetPassword) error {
	user, err := ux.ValidateRecoverPasswordToken(ctx, view)
	if err != nil {
		return err
	}

	hashed, _ := bcrypt.GenerateFromPassword([]byte(view.NewPassword), 8)

	if user.Password == string(hashed) {
		return utils.ErrSamePassword
	}

	user.Password = string(hashed)
	user.RecoverPasswordToken = ""

	err = ux.repo.Update(ctx, &user)
	if err != nil {
		return err
	}
	return nil
}

func (ux *UserUseCase) NewPasswordFirstLogin(ctx context.Context, username string, newPass string) error {
	user, err := ux.GetByUserName(ctx, username)
	if err != nil {
		return err
	}

	hashed, _ := bcrypt.GenerateFromPassword([]byte(newPass), 8)

	if user.Password == string(hashed) {
		return utils.ErrSamePassword
	}

	if !user.MustChangePassword {
		return utils.ErrHasChangedPassword
	}

	user.Password = string(hashed)
	user.MustChangePassword = false

	err = ux.repo.Update(ctx, &user)
	if err != nil {
		return err
	}
	return nil
}

func (uc *UserUseCase) ResetPasswordWithoutToken(ctx context.Context, view modelview.ResetPasswordWithoutToken) error {
	user, err := uc.repo.GetByUserName(ctx, view.Username)
	if err != nil {
		return err
	}

	hashed, _ := bcrypt.GenerateFromPassword([]byte(view.NewPassword), 8)

	if user.Password == string(hashed) {
		return utils.ErrSamePassword
	}

	user.Password = string(hashed)
	user.RecoverPasswordToken = ""

	err = uc.repo.Update(ctx, &user)
	if err != nil {
		return err
	}
	return nil
}
