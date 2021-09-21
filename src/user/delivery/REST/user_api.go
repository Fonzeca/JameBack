package REST

import (
	"context"
	"net/http"
	"time"

	"github.com/Fonzeka/Jame/src/domain"
	"github.com/Fonzeka/Jame/src/utils"
	"github.com/labstack/echo/v4"
)

type UserApi struct {
	useCase domain.UserUseCase
}

type json map[string]interface{}

//Constructor
func NewuserApi(useCase domain.UserUseCase) *UserApi {
	return &UserApi{useCase: useCase}
}

//Router
func (api *UserApi) Router(e *echo.Echo) {
	e.POST("/admin/user", api.InsertOne)
	e.PUT("/admin/user", api.UpdateOne)
	e.GET("/admin/user", api.GetUserByUserName)
	e.GET("/admin/users", api.GetAllusers)
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
		return utils.ErrBadRequestGetuser
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
	return c.NoContent(http.StatusOK)
}
