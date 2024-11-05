package main

import (
	"context"
	"fmt"
	"net/http"

	_RESTrole "github.com/Carmind-Mindia/user-hub/server/roles/delivery/REST"
	_mongoroles "github.com/Carmind-Mindia/user-hub/server/roles/repository/mongodb"
	_usecaseroles "github.com/Carmind-Mindia/user-hub/server/roles/usecase"
	_RESTuser "github.com/Carmind-Mindia/user-hub/server/user/delivery/REST"
	_mongouser "github.com/Carmind-Mindia/user-hub/server/user/repository/mongodb"
	_usecaseuser "github.com/Carmind-Mindia/user-hub/server/user/usecase"
	"github.com/Carmind-Mindia/user-hub/server/utils"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/qiniu/qmgo"
	"github.com/spf13/viper"
)

var Db *qmgo.Database

func init() {
	utils.InitConfig()

	db, err := initDataBase()
	if err != nil {
		fmt.Println(err)
		return
	}
	Db = db
}

func main() {

	reporoles := _mongoroles.NewMongoRolesRepository(Db)
	repousers := _mongouser.NewMongoUserRepository(Db)

	repoUseCase := _usecaseroles.NewRolesUseCase(reporoles)
	userUseCase := _usecaseuser.NewUserUseCase(repousers, repoUseCase)

	rolesApi := _RESTrole.NewuserApi(repoUseCase)
	userApi := _RESTuser.NewuserApi(userUseCase)

	e := echo.New()

	// Root level middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

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

func initDataBase() (*qmgo.Database, error) {
	ctx := context.Background()

	url := viper.GetString("db.url")
	port := viper.GetString("db.port")
	username := viper.GetString("db.user")
	password := viper.GetString("db.password")
	database := viper.GetString("db.database")
	authSource := viper.GetString("db.authSource")

	uri := "mongodb://" + url + ":" + port + "/" + database

	var config qmgo.Config

	if len(username) > 0 {
		config = qmgo.Config{Uri: uri, Database: database, Coll: "user", Auth: &qmgo.Credential{AuthSource: authSource, Username: username, Password: password}}
	} else {
		config = qmgo.Config{Uri: uri, Database: database, Coll: "user"}
	}

	cli, err := qmgo.Open(ctx, &config)
	if err != nil {
		return nil, err
	}

	return cli.Database, nil
}
