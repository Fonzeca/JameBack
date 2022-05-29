package REST

import (
	"context"
	"net/http"
	"time"

	"github.com/Fonzeca/UserHub/server/domain"
	our_jwt "github.com/Fonzeca/UserHub/server/security/jwt"
	"github.com/Fonzeca/UserHub/server/user/delivery/modelview"
	"github.com/Fonzeca/UserHub/server/user/usecase"
	"github.com/Fonzeca/UserHub/server/utils"
	"github.com/labstack/echo/v4"
)

type UserApi struct {
	useCase usecase.UserUseCase
}

type json map[string]interface{}

//Constructor
func NewuserApi(useCase usecase.UserUseCase) *UserApi {
	return &UserApi{useCase: useCase}
}

//Router
func (api *UserApi) Router(e *echo.Echo) {

	e.POST("/admin/user", api.InsertOne)
	e.PUT("/admin/user", api.UpdateOne)
	e.GET("/admin/user", api.GetUserByUserName)
	e.GET("/admin/users", api.GetAllusers)

	e.POST("/public/recoverPassword", api.SendEmailToRecoverPassword)
	e.POST("/public/validateRecoverToken", api.ValidateRecoverPasswordToken)
	e.POST("/public/resetPassword", api.ResetPasswordWithToken)

	e.GET("/logged", api.GetUserLogged)
	e.POST("/validate", api.ValidateToken)
	e.POST("/login", api.Login)
}

//Handlers ---------------

func (api *UserApi) Login(c echo.Context) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	user := domain.User{}
	c.Bind(&user)

	userName := user.UserName
	password := user.Password

	token, err := api.useCase.Login(ctx, userName, password)

	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, json{"token": token})
}

func (api *UserApi) InsertOne(c echo.Context) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	user := domain.User{}
	c.Bind(&user)

	user, err := api.useCase.Insert(ctx, &user)
	if err != nil {
		return err
	}

	return c.NoContent(http.StatusOK)
}

func (api *UserApi) GetUserByUserName(c echo.Context) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	values, _ := c.FormParams()
	username := values.Get("userName")
	if len(username) <= 0 {
		return utils.ErrBadRequest
	}

	usr, err := api.useCase.GetByUserName(ctx, username)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, usr)
}

func (api *UserApi) GetAllusers(c echo.Context) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	users, err := api.useCase.GetAll(ctx)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, users)
}

func (api *UserApi) UpdateOne(c echo.Context) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	user := domain.User{}
	c.Bind(&user)

	err := api.useCase.Update(ctx, &user)
	if err != nil {
		return err
	}

	return c.NoContent(http.StatusOK)
}

func (api *UserApi) ValidateToken(c echo.Context) error {
	_, err := our_jwt.ValidateAuth(c)
	if err != nil {
		return err
	}
	return c.NoContent(http.StatusOK)
}

func (api *UserApi) GetUserLogged(c echo.Context) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	claims, err := our_jwt.ValidateAuth(c)
	if err != nil {
		return err
	}

	user, err := api.useCase.GetUserByToken(ctx, claims)

	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, user)
}

func (api *UserApi) SendEmailToRecoverPassword(c echo.Context) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	email := c.QueryParams().Get("email")
	name := c.QueryParams().Get("name")

	err := api.useCase.SendEmailRecoverPassword(ctx, email, name)
	if err != nil {
		return err
	}

	return c.NoContent(http.StatusOK)

}

func (api *UserApi) ValidateRecoverPasswordToken(c echo.Context) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	modelview := modelview.ResetPassword{}
	c.Bind(&modelview)

	_, err := api.useCase.ValidateRecoverPasswordToken(ctx, modelview)
	if err != nil {
		return err
	}

	return c.NoContent(http.StatusOK)
}

func (api *UserApi) ResetPasswordWithToken(c echo.Context) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	modelview := modelview.ResetPassword{}
	c.Bind(&modelview)

	err := api.useCase.ResetPasswordWithToken(ctx, modelview)
	if err != nil {
		return err
	}

	return c.NoContent(http.StatusOK)
}
