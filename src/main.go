package main

import (
	"fmt"

	"context"

	"github.com/qiniu/qmgo"
	"github.com/qiniu/qmgo/options"
)

func main() {

	initApi()
}

type User struct {
	userName string `bson:"userName"`
}

func connect() {
	ctx := context.Background()
	cli, err := qmgo.Open(ctx, &qmgo.Config{Uri: "mongodb://localhost:27017", Database: "Login", Coll: "user"})

	defer func() {
		if err = cli.Close(ctx); err != nil {
			panic(err)
		}
	}()
	cli.CreateOneIndex(ctx, options.IndexModel{Key: []string{"userName"}, Unique: true})

	user := User{
		userName: "afonzo",
	}

	result, _ := cli.InsertOne(ctx, user)

	fmt.Println(result.InsertedID)
}
