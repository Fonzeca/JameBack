package REST

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/Carmind-Mindia/user-hub/server/domain"
	"github.com/Carmind-Mindia/user-hub/server/security"
	"github.com/Carmind-Mindia/user-hub/server/user/delivery/modelview"
	"github.com/Carmind-Mindia/user-hub/server/user/usecase"
	"github.com/Carmind-Mindia/user-hub/server/utils"
	"github.com/labstack/echo/v4"
)

type UserApi struct {
	useCase usecase.UserUseCase
}

type json map[string]interface{}

// Constructor
func NewuserApi(useCase usecase.UserUseCase) *UserApi {

	return &UserApi{useCase: useCase}
}

// Router
func (api *UserApi) Router(e *echo.Echo) {
	userGroup := e.Group("/user")

	userGroup.POST("/", api.InsertOne, security.ParseHeadersMiddleware)
	userGroup.PUT("/", api.UpdateOne, security.ParseHeadersMiddleware)
	userGroup.DELETE("/", api.DeleteOne, security.ParseHeadersMiddleware)
	userGroup.GET("/username/:username", api.GetUserByUserName, security.ParseHeadersMiddleware)
	userGroup.GET("/all", api.GetAllusers, security.ParseHeadersMiddleware)

	pwGroup := e.Group("/password")

	pwGroup.POST("/recover", api.SendEmailToRecoverPassword)
	pwGroup.POST("/validateToken", api.ValidateRecoverPasswordToken)
	pwGroup.POST("/reset", api.ResetPasswordWithToken)

	e.GET("/logged", api.GetUserLogged, security.ParseHeadersMiddleware)
	e.POST("/validate", api.ValidateToken, security.ParseHeadersMiddleware)
	e.POST("/firstLoginResetPassword", api.FirstLoginResetPassword, security.ParseHeadersMiddleware)
	e.POST("/login", api.Login)
}

//Handlers ---------------

// Login de usuarios
func (api *UserApi) Login(c echo.Context) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	user := domain.User{}
	c.Bind(&user)

	fmt.Println(user)

	userName := user.UserName
	password := user.Password

	response, err := api.useCase.Login(ctx, userName, password, c)

	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, response)
}

// Insertar un usuario
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

// Get usuario by username
func (api *UserApi) GetUserByUserName(c echo.Context) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	username := c.Param("username")
	if len(username) <= 0 {
		return utils.ErrBadRequest
	}

	usr, err := api.useCase.GetByUserName(ctx, username)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, usr)
}

// Get all users
func (api *UserApi) GetAllusers(c echo.Context) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	users, err := api.useCase.GetAll(ctx)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, users)
}

// Editar un usuario
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

// Borrar un usuario
func (api *UserApi) DeleteOne(c echo.Context) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	username := c.QueryParams().Get("username")

	err := api.useCase.Delete(ctx, username)
	if err != nil {
		return err
	}

	return c.NoContent(http.StatusOK)
}

// Valida el token
// Obtiene los datos del contexto de echo y los convierte en un objeto User
func (api *UserApi) ValidateToken(c echo.Context) error {

	// Obtener los datos del contexto de echo
	username := c.Get("username").(string)
	isAdmin := c.Get("admin").(bool)
	roles := c.Get("roles").([]string)

	// Crear objeto anónimo para retornar un JSON
	response := struct {
		Username string   `json:"username"`
		IsAdmin  bool     `json:"isAdmin"`
		Roles    []string `json:"roles"`
	}{
		Username: username,
		IsAdmin:  isAdmin,
		Roles:    roles,
	}

	return c.JSON(http.StatusOK, response)
}

// Obtiene un usuario logueado
func (api *UserApi) GetUserLogged(c echo.Context) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	// Obtener los datos del contexto de echo
	username := c.Get("username").(string)

	user, err := api.useCase.GetByUserName(ctx, username)

	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, user)
}

// Envia el email para recuperar la contraseña
func (api *UserApi) SendEmailToRecoverPassword(c echo.Context) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	email := c.QueryParams().Get("email")

	fmt.Println(email)
	fmt.Println(c.Request().URL)

	err := api.useCase.SendEmailRecoverPassword(ctx, email)
	if err != nil {
		return err
	}

	return c.NoContent(http.StatusOK)

}

// TODO: borrar
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

// Resetea la contraseña con el token
func (api *UserApi) ResetPasswordWithToken(c echo.Context) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	modelview := modelview.ResetPassword{}
	c.Bind(&modelview)

	if modelview.NewPassword == "" || modelview.Token == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Bad request")
	}

	err := api.useCase.ResetPasswordWithToken(ctx, modelview)
	if err != nil {
		return err
	}

	return c.NoContent(http.StatusOK)
}

// Cambiamos la contraseña si es el primer inicio de sesion
func (api *UserApi) FirstLoginResetPassword(c echo.Context) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	//Obtenemos la nueva contraseña
	newPass := c.QueryParam("newPassword")
	username := c.Get("username").(string)

	// Validate newPass and username
	if newPass == "" || username == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Bad request")
	}

	err := api.useCase.NewPasswordFirstLogin(ctx, username, newPass)
	if err != nil {
		return err
	}

	return c.NoContent(http.StatusOK)
}
