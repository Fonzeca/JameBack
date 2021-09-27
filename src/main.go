package main

import (
	"context"
	"net/http"

	_RESTrole "github.com/Fonzeka/Jame/src/roles/delivery/REST"
	_mongoroles "github.com/Fonzeka/Jame/src/roles/repository/mongodb"
	_usecaseroles "github.com/Fonzeka/Jame/src/roles/usecase"
	"github.com/Fonzeka/Jame/src/security/jwt"
	_RESTuser "github.com/Fonzeka/Jame/src/user/delivery/REST"
	_mongouser "github.com/Fonzeka/Jame/src/user/repository/mongodb"
	_usecaseuser "github.com/Fonzeka/Jame/src/user/usecase"
	"github.com/Fonzeka/Jame/src/utils"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/qiniu/qmgo"
)

func main() {

	utils.InitConfig()

	db, _ := initDataBase()

	reporoles := _mongoroles.NewMongoRolesRepository(db)
	repousers := _mongouser.NewMongoUserRepository(db)

	repoUseCase := _usecaseroles.NewRolesUseCase(reporoles)
	userUseCase := _usecaseuser.NewUserUseCase(repousers, repoUseCase)

	rolesApi := _RESTrole.NewuserApi(repoUseCase)
	userApi := _RESTuser.NewuserApi(userUseCase)

	e := echo.New()

	// Root level middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(jwt.CheckLogged)
	e.HTTPErrorHandler = customHTTPErrorHandler

	userApi.Router(e)
	rolesApi.Router(e)

	e.Logger.Fatal(e.Start(":1323"))
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
	cli, err := qmgo.Open(ctx, &qmgo.Config{Uri: "mongodb://localhost:27017", Database: "UsersHub", Coll: "user"})
	if err != nil {
		return nil, err
	}

	return cli.Database, nil

}
