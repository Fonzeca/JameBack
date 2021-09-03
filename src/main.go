package main

import (
	"context"

	"github.com/Fonzeka/Jame/src/user/delivery/REST"
	"github.com/Fonzeka/Jame/src/user/repository/mongodb"
	"github.com/Fonzeka/Jame/src/user/usecase"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/qiniu/qmgo"
)

func main() {

	db, _ := initDataBase()

	repousers := mongodb.NewMongoUserRepository(db)

	userUseCase := usecase.NewUserUseCase(&repousers)

	userApi := REST.NewuserApi(userUseCase)

	e := echo.New()

	// Root level middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	userApi.Router(e)

	e.Logger.Fatal(e.Start(":1323"))
}

func initDataBase() (*qmgo.Database, error) {
	ctx := context.Background()
	cli, err := qmgo.Open(ctx, &qmgo.Config{Uri: "mongodb://localhost:27017", Database: "Login", Coll: "user"})
	if err != nil {
		return nil, err
	}

	return cli.Database, nil

}
