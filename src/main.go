package main

import (
	"context"
	"net/http"

	"github.com/Fonzeka/Jame/src/user/delivery/REST"
	"github.com/Fonzeka/Jame/src/user/repository/mongodb"
	"github.com/Fonzeka/Jame/src/user/usecase"
	"github.com/Fonzeka/Jame/src/utils"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/qiniu/qmgo"
)

func main() {

	utils.InitConfig()

	db, _ := initDataBase()

	repousers := mongodb.NewMongoUserRepository(db)

	userUseCase := usecase.NewUserUseCase(repousers)

	userApi := REST.NewuserApi(userUseCase)

	e := echo.New()

	// Root level middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(usecase.CheckLogged)
	e.HTTPErrorHandler = customHTTPErrorHandler

	userApi.Router(e)

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
