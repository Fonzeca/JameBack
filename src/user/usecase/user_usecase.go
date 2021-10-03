package usecase

import (
	"context"
	"time"

	"github.com/Fonzeka/Jame/src/domain"
	"github.com/Fonzeka/Jame/src/roles/usecase"
	myjwt "github.com/Fonzeka/Jame/src/security/jwt"
	"github.com/Fonzeka/Jame/src/utils"
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
		if len(uc.rolecase.GetAllRoles(ctx)) <= {
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
