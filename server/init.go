package server

import (
	"context"
	"net/http"

	"github.com/Fonzeca/FastEmail/src/sdk"
	guard_userhub "github.com/Fonzeca/UserHub/guard"
	_RESTrole "github.com/Fonzeca/UserHub/server/roles/delivery/REST"
	_mongoroles "github.com/Fonzeca/UserHub/server/roles/repository/mongodb"
	_usecaseroles "github.com/Fonzeca/UserHub/server/roles/usecase"
	_RESTuser "github.com/Fonzeca/UserHub/server/user/delivery/REST"
	_mongouser "github.com/Fonzeca/UserHub/server/user/repository/mongodb"
	_usecaseuser "github.com/Fonzeca/UserHub/server/user/usecase"
	"github.com/Fonzeca/UserHub/server/utils"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/qiniu/qmgo"
	"github.com/spf13/viper"
)

var Db *qmgo.Database
var Guard *guard_userhub.Guard

func InitServer() {
	fastEmailConfig := sdk.Config{
		Url: viper.GetString("fast-email.url"),
	}

	client := sdk.NewEmailClient(fastEmailConfig)

	reporoles := _mongoroles.NewMongoRolesRepository(Db)
	repousers := _mongouser.NewMongoUserRepository(Db)

	repoUseCase := _usecaseroles.NewRolesUseCase(reporoles)
	userUseCase := _usecaseuser.NewUserUseCase(repousers, repoUseCase, &client)

	rolesApi := _RESTrole.NewuserApi(repoUseCase)
	userApi := _RESTuser.NewuserApi(userUseCase)

	e := echo.New()

	// Root level middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
	}))
	// e.Use(jwt.CheckLogged)
	e.Use(Guard.EchoMiddlewareApiKey)
	e.HTTPErrorHandler = customHTTPErrorHandler

	userApi.Router(e)
	rolesApi.Router(e)

	port := viper.GetString("server.port")

	e.Logger.Fatal(e.Start(":" + port))
}

func customHTTPErrorHandler(err error, c echo.Context) {
	var (
		code = http.StatusInternalServerError
		key  = "ServerError"
		msg  string
	)

	if he, ok := err.(*utils.HttpError); ok {
		code = he.Code
		key = he.Key
		msg = he.Message
	} else if true { //Aca va si estamos en debug
		msg = err.Error()
	} else {
		msg = http.StatusText(code)
	}

	err = c.JSON(code, utils.NewHTTPError(code, key, msg))
	if err != nil {
		c.Logger().Error(err)
	}
}

func InitDataBase() (*qmgo.Database, error) {
	ctx := context.Background()

	url := viper.GetString("db.url")
	port := viper.GetString("db.port")
	username := viper.GetString("db.user")
	password := viper.GetString("db.password")
	database := viper.GetString("db.database")

	uri := "mongodb://" + url + ":" + port + "/" + database

	var config qmgo.Config

	if len(username) > 0 {
		config = qmgo.Config{Uri: uri, Database: database, Coll: "user", Auth: &qmgo.Credential{AuthSource: database, Username: username, Password: password}}
	} else {
		config = qmgo.Config{Uri: uri, Database: database, Coll: "user"}
	}

	cli, err := qmgo.Open(ctx, &config)
	if err != nil {
		return nil, err
	}

	return cli.Database, nil
}
