package REST

import (
	"context"
	"net/http"
	"time"

	"github.com/Fonzeka/Jame/src/domain"
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
	e.POST("/user", api.InsertOne)
	e.GET("/users", api.GetAllusers)
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

func (api *UserApi) GetAllusers(c echo.Context) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	users, err := api.useCase.GetAll(ctx)
	if err != nil {
		c.Error(err)
		return err
	}

	return c.JSON(http.StatusOK, users)
}
